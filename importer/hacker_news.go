package importer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sourcegraph.com/sourcegraph/thesrc"
)

func init() {
	Fetchers = append(Fetchers, &hackerNews{"top"}, &hackerNews{"newest"}, &hackerNews{"best"})
}

type hackerNews struct {
	which string
}

func (f *hackerNews) Fetch() ([]*thesrc.Post, error) {
	resp, err := http.Get("http://hnify.herokuapp.com/get/" + f.which)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 HTTP response status: %d", resp.StatusCode)
	}

	var results *struct {
		Stories []*struct {
			Title  string
			Link   string
			Points int
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	posts := make([]*thesrc.Post, len(results.Stories))
	for i, s := range results.Stories {
		posts[i] = &thesrc.Post{
			Title:   s.Title,
			LinkURL: s.Link,
			Score:   s.Points,
		}
	}

	return posts, nil
}

func (f *hackerNews) Site() string { return "hn/" + f.which }
