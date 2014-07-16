package router

import "github.com/gorilla/mux"

func App() *mux.Router {
	m := mux.NewRouter()
	m.Path("/").Methods("GET").Name(Posts)
	m.Path("/p/{ID:.+}").Methods("GET").Name(Post)
	return m
}
