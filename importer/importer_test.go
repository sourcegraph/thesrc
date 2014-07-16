package importer

import (
	"testing"

	"github.com/sourcegraph/thesrc"
	"github.com/sourcegraph/thesrc/datastore"
)

type mockFetcher struct {
	posts []*thesrc.Post
	err   error
}

func (f *mockFetcher) Fetch() ([]*thesrc.Post, error) { return f.posts, f.err }
func (f *mockFetcher) Site() string                   { return "mock" }

func TestImport(t *testing.T) {
	want := &thesrc.Post{Title: "t"}

	var createCalled bool
	store = &datastore.Datastore{
		Posts: &thesrc.MockPostsService{
			Create_: func(post *thesrc.Post) error {
				if post.Title != want.Title {
					t.Errorf("got title %q, want %q", post.Title, want.Title)
				}
				createCalled = true
				return nil
			},
		},
	}

	var imported int
	Imported = func(site string, post *thesrc.Post) {
		imported++
	}

	f := &mockFetcher{posts: []*thesrc.Post{want}}
	if err := Import(f); err != nil {
		t.Fatal(err)
	}

	if !createCalled {
		t.Error("!createCalled")
	}

	if want := 1; imported != want {
		t.Errorf("got imported == %d, want %d", imported, want)
	}
}
