package router

import "github.com/gorilla/mux"

// App-only routes
const (
	SubmitPostForm = "post:submit-form"
)

func App() *mux.Router {
	m := mux.NewRouter()
	m.Path("/").Methods("GET").Name(Posts)
	m.Path("/p/{ID:.+}").Methods("GET").Name(Post)
	m.Path("/submit").Methods("GET").Name(SubmitPostForm)
	m.Path("/posts").Methods("POST").Name(SubmitPost)
	return m
}
