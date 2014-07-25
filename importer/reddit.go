package importer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sourcegraph/thesrc"
)

func init() {
	Fetchers = append(Fetchers,
		&subreddit{"programming"}, &subreddit{"golang"}, &subreddit{"postgresql"},
		// &subreddit{"ruby"}, &subreddit{"node"},
		// &subreddit{"python"}, , &subreddit{"django"},
		// &subreddit{"rust"},
	)
}

type subreddit struct {
	name string
}

func (f *subreddit) Fetch() ([]*thesrc.Post, error) {
	postsMap := map[string]*thesrc.Post{}
	patterns := []string{
		"http://www.reddit.com/r/%s/hot.json",
		"http://www.reddit.com/r/%s/new.json",
		"http://www.reddit.com/r/%s/top.json",
	}
	for _, pat := range patterns {
		posts2, err := f.fetchOne(fmt.Sprintf(pat, f.name))
		if err != nil {
			return nil, err
		}
		for _, p2 := range posts2 {
			postsMap[p2.LinkURL] = p2
		}
	}

	posts := make([]*thesrc.Post, 0, len(postsMap))
	for _, post := range postsMap {
		posts = append(posts, post)
	}
	return posts, nil
}

func (f *subreddit) fetchOne(urlStr string) ([]*thesrc.Post, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 HTTP response status: %d", resp.StatusCode)
	}

	var results *struct {
		Data struct {
			Children []*struct {
				Data struct {
					Title string
					URL   string
					Score int
				}
			}
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	posts := make([]*thesrc.Post, len(results.Data.Children))
	for i, s := range results.Data.Children {
		posts[i] = &thesrc.Post{
			Title:   s.Data.Title,
			LinkURL: s.Data.URL,
			Score:   s.Data.Score,
		}
	}

	return posts, nil
}

func (f *subreddit) Site() string { return "/r/" + f.name }
