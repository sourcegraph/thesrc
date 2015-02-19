package importer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sourcegraph.com/sourcegraph/thesrc"
)

func init() {
	Fetchers = append(Fetchers, &lobsters{"hottest"}, &lobsters{"newest"})
}

type lobsters struct {
	which string
}

func (f *lobsters) Fetch() ([]*thesrc.Post, error) {
	resp, err := http.Get(fmt.Sprintf("https://lobste.rs/%s.json", f.which))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 HTTP response status: %d", resp.StatusCode)
	}

	var results []*struct {
		Title string
		URL   string
		Score int
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	posts := make([]*thesrc.Post, len(results))
	for i, s := range results {
		posts[i] = &thesrc.Post{
			Title:   s.Title,
			LinkURL: s.URL,
			Score:   s.Score,
		}
	}

	return posts, nil
}

func (f *lobsters) Site() string { return "lobsters/" + f.which }
