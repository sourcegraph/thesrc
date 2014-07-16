package datastore

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jmoiron/modl"
	"github.com/sourcegraph/thesrc"
)

func init() {
	DB.AddTableWithName(thesrc.Post{}, "post").SetKeys(true, "ID")
	createSQL = append(createSQL,
		`CREATE INDEX post_submittedat ON post(submittedat DESC);`,
		`CREATE UNIQUE INDEX post_linkurl ON post(linkurl);`,
	)

}

type postsStore struct{ *Datastore }

func (s *postsStore) Get(id int) (*thesrc.Post, error) {
	var posts []*thesrc.Post
	if err := s.dbh.Select(&posts, `SELECT * FROM post WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, thesrc.ErrPostNotFound
	}
	return posts[0], nil
}

func (s *postsStore) List(opt *thesrc.PostListOptions) ([]*thesrc.Post, error) {
	if opt == nil {
		opt = &thesrc.PostListOptions{}
	}
	var posts []*thesrc.Post
	err := s.dbh.Select(&posts, `SELECT * FROM post LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postsStore) Submit(post *thesrc.Post) (bool, error) {
	retries := 3
	var wantRetry bool

retry:
	retries--
	wantRetry = false
	if retries == 0 {
		return false, fmt.Errorf("failed to submit post with URL %q after retrying", post.LinkURL)
	}

	var created bool
	err := transact(s.dbh, func(tx modl.SqlExecutor) error {
		var existing []*thesrc.Post
		if err := tx.Select(&existing, `SELECT * FROM post WHERE linkurl=$1 LIMIT 1;`, post.LinkURL); err != nil {
			return err
		}
		if len(existing) > 0 {
			*post = *existing[0]
			return nil
		}

		if err := tx.Insert(post); err != nil {
			if strings.Contains(err.Error(), `violates unique constraint "post_linkurl"`) {
				time.Sleep(time.Duration(rand.Intn(75)) * time.Millisecond)
				wantRetry = true
				return err
			}
			return err
		}

		created = true
		return nil
	})
	if wantRetry {
		goto retry
	}
	return created, err
}
