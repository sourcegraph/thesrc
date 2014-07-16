package app

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sourcegraph/thesrc"
	"github.com/sourcegraph/thesrc/router"
)

var (
	// ReloadTemplates is whether to reload templates on each request.
	ReloadTemplates bool

	// StaticDir is the directory containing static assets.
	StaticDir = filepath.Join(defaultBase("github.com/sourcegraph/thesrc/app"), "static")
)

var (
	apiclient     = thesrc.NewClient(nil)
	schemaDecoder = schema.NewDecoder()
	appRouter     = router.App()
)

func Handler() *mux.Router {
	m := appRouter
	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(StaticDir))))
	// TODO(sqs): add handlers for /favicon.ico and /robots.txt
	m.Get(router.Post).Handler(handler(servePost))
	m.Get(router.Posts).Handler(handler(servePosts))
	m.Get(router.CreatePostForm).Handler(handler(serveCreatePostForm))
	m.Get(router.CreatePost).Handler(handler(serveCreatePost))
	return m
}

type handler func(resp http.ResponseWriter, req *http.Request) error

func (h handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if ReloadTemplates {
		LoadTemplates()
	}
	runHandler(resp, req, h)
}

func runHandler(w http.ResponseWriter, r *http.Request, fn func(http.ResponseWriter, *http.Request) error) {
	var err error

	defer func() {
		if rv := recover(); rv != nil {
			err = errors.New("handler panic")
			logError(r, err, rv)
			handleError(w, r, http.StatusInternalServerError, err)
		}
	}()

	err = fn(w, r)
	if err != nil {
		logError(r, err, nil)
		handleError(w, r, http.StatusInternalServerError, err)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("cache-control", "no-cache")
	err2 := renderTemplate(w, r, "error.html", status, &struct {
		StatusCode int
		Status     string
		Err        error
		templateCommon
	}{
		StatusCode: status,
		Status:     http.StatusText(status),
		Err:        err,
	})
	if err2 != nil {
		logError(r, fmt.Errorf("during execution of error template: %s", err2), nil)
	}
}

func logError(req *http.Request, err error, rv interface{}) {
	if err != nil {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Error serving %s (route %s): %s\n", req.URL, mux.CurrentRoute(req).GetName(), err)
		if rv != nil {
			fmt.Fprintln(&buf, rv)
			buf.Write(debug.Stack())
		}
		log.Print(buf.String())
	}
}
