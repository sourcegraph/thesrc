package router

import "github.com/gorilla/mux"

// App-only routes
const (
	CreatePostForm = "post:create-form"
)

func App() *mux.Router {
	m := mux.NewRouter()
	m.Path("/").Methods("GET").Name(Posts)
	m.Path("/p/{ID:.+}").Methods("GET").Name(Post)
	m.Path("/submit").Methods("GET").Name(CreatePostForm)
	m.Path("/posts").Methods("POST").Name(CreatePost)
	return m
}
