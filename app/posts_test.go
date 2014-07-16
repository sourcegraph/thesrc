package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
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

func TestSubmitPostForm(t *testing.T) {
	setup()
	defer teardown()

	want := &thesrc.Post{
		Title:   "t",
		LinkURL: "http://example.com",
		Body:    "b",
	}

	url_, _ := router.App().Get(router.SubmitPostForm).URL()
	url_.RawQuery = url.Values{"Title": []string{want.Title}, "url": []string{want.LinkURL}, "body": []string{want.Body}}.Encode()
	html, resp := getHTML(t, url_)

	if want := http.StatusOK; resp.Code != want {
		t.Errorf("got HTTP status %d, want %d", resp.Code, want)
	}

	if got, _ := html.Find("input[name=Title]").Attr("value"); got != want.Title {
		t.Errorf("got title %q, want %q", got, want.Title)
	}
	if got, _ := html.Find("input[name=LinkURL]").Attr("value"); got != want.LinkURL {
		t.Errorf("got link href %q, want %q", got, want.LinkURL)
	}
	if body := html.Find("textarea[name=Body]").Text(); body != want.Body {
		t.Errorf("got post body %q, want %q", body, want.Body)
	}
}

func TestSubmitPosts(t *testing.T) {
	setup()
	defer teardown()

	post := &thesrc.Post{ID: 0, Title: "t", LinkURL: "http://example.com", Body: "b"}

	var called bool
	apiclient = &thesrc.Client{
		Posts: &thesrc.MockPostsService{
			Submit_: func(post *thesrc.Post) (bool, error) {
				called = true
				post.ID = 1
				return true, nil
			},
		},
	}

	v := url.Values{
		"Title":   []string{post.Title},
		"LinkURL": []string{post.LinkURL},
		"Body":    []string{post.Body},
	}

	url, _ := router.App().Get(router.SubmitPost).URL()
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	resp.Body = new(bytes.Buffer)
	testMux.ServeHTTP(resp, req)

	if want := http.StatusSeeOther; resp.Code != want {
		t.Errorf("got HTTP status %d, want %d", resp.Code, want)
	}

	if !called {
		t.Error("!called")
	}

	if loc, want := resp.Header().Get("location"), urlTo(router.Post, "ID", "1").String(); loc != want {
		t.Errorf("got Location %q, want %q", loc, want)
	}
}
