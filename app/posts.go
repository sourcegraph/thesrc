package app

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sourcegraph/thesrc"
	"github.com/sourcegraph/thesrc/router"
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

	opt.CodeOnly = true

	if opt.PerPage == 0 {
		opt.PerPage = 60
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

func serveSubmitPostForm(w http.ResponseWriter, r *http.Request) error {
	// Populate form from querystring.
	q := r.URL.Query()
	post := &thesrc.Post{
		Title:   getCaseOrLowerCaseQuery(q, "Title"),
		LinkURL: getCaseOrLowerCaseQuery(q, "LinkURL") + getCaseOrLowerCaseQuery(q, "URL"), // support both
		Body:    getCaseOrLowerCaseQuery(q, "Body"),
	}

	return renderTemplate(w, r, "posts/submit_form.html", http.StatusOK, struct {
		Post *thesrc.Post
	}{
		Post: post,
	})
}

func serveSubmitPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	var post thesrc.Post
	if err := schemaDecoder.Decode(&post, r.Form); err != nil {
		return err
	}

	if _, err := apiclient.Posts.Submit(&post); err != nil {
		return err
	}

	postURL := urlTo(router.Post, "ID", strconv.Itoa(post.ID))
	http.Redirect(w, r, postURL.String(), http.StatusSeeOther)
	return nil
}

func getCaseOrLowerCaseQuery(q url.Values, name string) string {
	if v, present := q[name]; present {
		return v[0]
	}
	return q.Get(strings.ToLower(name))
}
