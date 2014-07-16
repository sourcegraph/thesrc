package api

import "github.com/gorilla/mux"

func Handler() *mux.Router {
	return mux.NewRouter()
}
