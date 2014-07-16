package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sourcegraph/thesrc"
)

func servePost(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		return err
	}

	post, err := store.Posts.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, post)
}

func serveCreatePost(w http.ResponseWriter, r *http.Request) error {
	var post thesrc.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return err
	}

	if err := store.Posts.Create(&post); err != nil {
		return err
	}

	return writeJSON(w, post)
}

func servePosts(w http.ResponseWriter, r *http.Request) error {
	var opt thesrc.PostListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	posts, err := store.Posts.List(&opt)
	if err != nil {
		return err
	}
	if posts == nil {
		posts = []*thesrc.Post{}
	}

	return writeJSON(w, posts)
}
