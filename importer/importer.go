package importer

import "sourcegraph.com/sourcegraph/thesrc"

var Fetchers = []Fetcher{}

// A Fetcher fetches posts from other sites.
type Fetcher interface {
	// Fetch posts.
	Fetch() ([]*thesrc.Post, error)

	// Site is the name of the site that this Fetcher fetches from.
	Site() string
}

var Store = thesrc.NewClient(nil)

// Import posts fetched by f. If Imported is non-nil, it is called each time a
// post is successfully imported.
func Import(f Fetcher) error {
	posts, err := f.Fetch()
	if err != nil {
		return err
	}

	for _, post := range posts {
		created, err := Store.Posts.Submit(post)
		if err != nil {
			return err
		}
		if Imported != nil {
			Imported(f.Site(), post, created)
		}
	}
	return nil
}

// Imported (if non-nil) is called each time a post is successfully imported.
var Imported func(site string, post *thesrc.Post, created bool)
