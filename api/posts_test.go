package api

import (
	"testing"

	"github.com/sourcegraph/thesrc"
)

func TestPost(t *testing.T) {
	setup()

	wantPost := &thesrc.Post{ID: 1}

	calledGet := false
	store.Posts.(*thesrc.MockPostsService).Get_ = func(id int) (*thesrc.Post, error) {
		if id != wantPost.ID {
			t.Errorf("wanted request for post %d but got %d", wantPost.ID, id)
		}
		calledGet = true
		return wantPost, nil
	}

	gotPost, err := apiClient.Posts.Get(wantPost.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(wantPost, gotPost) {
		t.Errorf("got post %+v but wanted post %+v", wantPost, gotPost)
	}
}

func TestPost_Submit(t *testing.T) {
	setup()

	wantPost := &thesrc.Post{ID: 1}

	calledPost := false
	store.Posts.(*thesrc.MockPostsService).Submit_ = func(post *thesrc.Post) (bool, error) {
		if !normalizeDeepEqual(wantPost, post) {
			t.Errorf("wanted request for post %+v but got %+v", wantPost, post)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Posts.Submit(wantPost)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPost {
		t.Error("!calledPost")
	}
	if !success {
		t.Errorf("!success")
	}
}

func TestPosts_List(t *testing.T) {
	setup()

	wantPosts := []*thesrc.Post{{ID: 1}}
	wantOpt := &thesrc.PostListOptions{ListOptions: thesrc.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Posts.(*thesrc.MockPostsService).List_ = func(opt *thesrc.PostListOptions) ([]*thesrc.Post, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted list options %+v but got %+v", wantOpt, opt)
		}
		calledList = true
		return wantPosts, nil
	}

	posts, err := apiClient.Posts.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}
	for i, _ := range posts {
		if !normalizeDeepEqual(wantPosts[i], posts[i]) {
			t.Errorf("got post %+v but wanted post %+v", posts[i], wantPosts[i])
		}
	}

}
