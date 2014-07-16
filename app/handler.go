package app

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/sqs/mux"
)

var (
	// ReloadTemplates is whether to reload templates on each request.
	ReloadTemplates bool

	// StaticDir is the directory containing static assets.
	StaticDir string
)

func Handler() *mux.Router {
	m := mux.NewRouter()
	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(StaticDir))))
	// TODO(sqs): add handlers for /favicon.ico and /robots.txt
	m.Path("/").Methods("GET").Handler(handler(serveHome))
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
	err2 := renderTemplate(w, r, "error/error.html", status, &struct {
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
