package app

import (
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

	post, err := apiclient.Posts.Get(id)
	if err != nil {
		return err
	}

	return renderTemplate(w, r, "posts/show.html", http.StatusOK, struct {
		Post *thesrc.Post
	}{
		Post: post,
	})
}

func servePosts(w http.ResponseWriter, r *http.Request) error {
	var opt thesrc.PostListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	posts, err := apiclient.Posts.List(&opt)
	if err != nil {
		return err
	}

	return renderTemplate(w, r, "posts/list.html", http.StatusOK, struct {
		Posts []*thesrc.Post
	}{
		Posts: posts,
	})
}
