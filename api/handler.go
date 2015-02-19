package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"sourcegraph.com/sourcegraph/thesrc/datastore"
	"sourcegraph.com/sourcegraph/thesrc/router"
)

var (
	store         = datastore.NewDatastore(nil)
	schemaDecoder = schema.NewDecoder()
)

func Handler() *mux.Router {
	m := router.API()
	m.Get(router.Post).Handler(handler(servePost))
	m.Get(router.SubmitPost).Handler(handler(serveSubmitPost))
	m.Get(router.Posts).Handler(handler(servePosts))
	return m
}

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		log.Println(err)
	}
}
