package importer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sourcegraph/thesrc"
)

type lobsters struct {
	endpoint string
}

var (
	LobstersHottest Fetcher = &lobsters{"https://lobste.rs/hottest.json"}
	LobstersNewest  Fetcher = &lobsters{"https://lobste.rs/newest.json"}
)

func (f *lobsters) Fetch() ([]*thesrc.Post, error) {
	resp, err := http.Get(f.endpoint)
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

func (f *lobsters) Site() string { return "lobsters" }
