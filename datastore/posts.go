package datastore

import "github.com/sourcegraph/thesrc"

func init() {
	DB.AddTableWithName(thesrc.Post{}, "post").SetKeys(true, "ID")
	createSQL = append(createSQL,
		`CREATE INDEX post_submittedat ON post(submittedat DESC);`,
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

func (s *postsStore) Create(post *thesrc.Post) error {
	return s.dbh.Insert(post)
}
