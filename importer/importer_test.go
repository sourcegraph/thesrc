package importer

import (
	"testing"

	"sourcegraph.com/sourcegraph/thesrc"
)

type mockFetcher struct {
	posts []*thesrc.Post
	err   error
}

func (f *mockFetcher) Fetch() ([]*thesrc.Post, error) { return f.posts, f.err }
func (f *mockFetcher) Site() string                   { return "mock" }

func TestImport(t *testing.T) {
	want := &thesrc.Post{Title: "t"}

	var submitCalled bool
	Store = &thesrc.Client{
		Posts: &thesrc.MockPostsService{
			Submit_: func(post *thesrc.Post) (bool, error) {
				if post.Title != want.Title {
					t.Errorf("got title %q, want %q", post.Title, want.Title)
				}
				submitCalled = true
				return true, nil
			},
		},
	}

	var imported int
	Imported = func(site string, post *thesrc.Post, created bool) {
		imported++
	}

	f := &mockFetcher{posts: []*thesrc.Post{want}}
	if err := Import(f); err != nil {
		t.Fatal(err)
	}

	if !submitCalled {
		t.Error("!submitCalled")
	}

	if want := 1; imported != want {
		t.Errorf("got imported == %d, want %d", imported, want)
	}
}
