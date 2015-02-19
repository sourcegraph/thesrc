package importer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sourcegraph/thesrc"
)

func init() {
	Fetchers = append(Fetchers, &qTag{"Ruby", "4", ""}, &qTag{"golang", "5", ""}, &qTag{"postgresql", "2", ""})
}

type qTag struct {	//Tags should all be in one string, semi-colon delimited
	tags string
	pageSize string	//Number of posts to be grabbed with these tags
	link string	//Link to be generated later, from above parameters
}

func (f *qTag) urlSetUp() {

	s := []string{f.pageSize, "&order=desc&sort=activity&tagged=", f.tags, "&site=stackoverflow"}
	f.link = strings.Join(s, "");

}

func (f *qTag) Fetch() ([]*thesrc.Post, error) {

	postsMap := map[string]*thesrc.Post{}
	patterns := []string{
		"https://api.stackexchange.com/2.2/search?page=1&pagesize=",
	}	

	f.urlSetUp();

	for _, pat := range patterns {
		posts2, err := f.fetchOne(strings.Join([]string{pat,f.link},""))
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

func (f *qTag) fetchOne(urlStr string) ([]*thesrc.Post, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 HTTP response status: %d", resp.StatusCode)
	}

	var results *struct {
		Items []*struct {
			Title string
			Score int
			Link  string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	posts := make([]*thesrc.Post, len(results.Items))
	for i, s := range results.Items {
		posts[i] = &thesrc.Post{
			Title:   s.Title,
			LinkURL: s.Link,
			Score:   s.Score,
		}

		fmt.Printf("TITLE: %s\nURL: %s\nScore: %d",s.Title, s.Link, s.Score)

	}

	return posts, nil
}

func (f *qTag) Site() string { return f.link }
