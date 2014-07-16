package app

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/sourcegraph/thesrc"
	"github.com/sourcegraph/thesrc/router"
)

func TestPost(t *testing.T) {
	setup()
	defer teardown()

	post := &thesrc.Post{ID: 1, Title: "t", LinkURL: "http://example.com", Body: "b"}

	var called bool
	apiclient = &thesrc.Client{
		Posts: &thesrc.MockPostsService{
			Get_: func(id int) (*thesrc.Post, error) {
				if id != post.ID {
					t.Errorf("got post ID %v, want %v", id, post.ID)
				}
				called = true
				return post, nil
			},
		},
	}

	url, _ := router.App().Get(router.Post).URL("ID", strconv.Itoa(post.ID))
	html, resp := getHTML(t, url)

	if want := http.StatusOK; resp.Code != want {
		t.Errorf("got HTTP status %d, want %d", resp.Code, want)
	}

	if !called {
		t.Error("!called")
	}

	a := html.Find("a.post-link")
	if a.Text() != post.Title {
		t.Errorf("got link text %q, want %q", a.Text(), post.Title)
	}
	if got, _ := a.Attr("href"); got != post.LinkURL {
		t.Errorf("got link href %q, want %q", got, post.LinkURL)
	}
	body := html.Find("p.post-body")
	if body.Text() != post.Body {
		t.Errorf("got post body %q, want %q", body.Text(), post.Body)
	}
}

func TestPosts(t *testing.T) {
	setup()
	defer teardown()

	posts := []*thesrc.Post{{ID: 1, Title: "t", LinkURL: "http://example.com", Body: "b"}}

	var called bool
	apiclient = &thesrc.Client{
		Posts: &thesrc.MockPostsService{
			List_: func(opt *thesrc.PostListOptions) ([]*thesrc.Post, error) {
				called = true
				return posts, nil
			},
		},
	}

	url, _ := router.App().Get(router.Posts).URL()
	html, resp := getHTML(t, url)

	if want := http.StatusOK; resp.Code != want {
		t.Errorf("got HTTP status %d, want %d", resp.Code, want)
	}

	if !called {
		t.Error("!called")
	}

	for _, post := range posts {
		a := html.Find("a.post-link")
		if a.Text() != post.Title {
			t.Errorf("got link text %q, want %q", a.Text(), post.Title)
		}
		if got, _ := a.Attr("href"); got != post.LinkURL {
			t.Errorf("got link href %q, want %q", got, post.LinkURL)
		}
		body := html.Find("p.post-body")
		if body.Text() != post.Body {
			t.Errorf("got post body %q, want %q", body.Text(), post.Body)
		}
	}
}
