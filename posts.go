package thesrc

import (
	"errors"
	"time"

	"github.com/sourcegraph/thesrc/router"
)

// A Post is a link and short body submitted to and displayed on thesrc.
type Post struct {
	// ID is a short, unique alphanumeric identifier for the post.
	ID string

	// Title of the post.
	Title string

	// LinkURL is the URL to a link that this post is about.
	LinkURL string

	// Body of the post.
	Body string

	// SubmittedAt is when the post was submitted.
	SubmittedAt time.Time

	// AuthorUserID is the user ID of this post's author.
	AuthorUserID int
}

// PostsService interacts with the post-related endpoints in thesrc's API.
type PostsService interface {
	// Get a post.
	Get(id string) (*Post, error)

	// List posts.
	List(opt *PostListOptions) ([]*Post, error)
}

var (
	ErrPostNotFound = errors.New("post not found")
)

type postsService struct{ client *Client }

func (s *postsService) Get(id string) (*Post, error) {
	url, err := s.client.url(router.Post, map[string]string{"ID": id}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var post *Post
	_, err = s.client.Do(req, &post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

type PostListOptions struct {
	ListOptions
}

func (s *postsService) List(opt *PostListOptions) ([]*Post, error) {
	url, err := s.client.url(router.Posts, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	_, err = s.client.Do(req, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

type MockPostsService struct {
	Get_  func(id string) (*Post, error)
	List_ func(opt *PostListOptions) ([]*Post, error)
}

func (s *MockPostsService) Get(id string) (*Post, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockPostsService) List(opt *PostListOptions) ([]*Post, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}
