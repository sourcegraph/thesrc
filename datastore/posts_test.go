package datastore

import (
	"reflect"
	"testing"

	"github.com/sourcegraph/thesrc"
)

func TestPostsStore_Get_db(t *testing.T) {
	want := &thesrc.Post{ID: 1}

	tx, _ := DB.Begin()
	defer tx.Rollback()
	tx.Exec(`DELETE FROM post;`) // test on a clean DB
	if err := tx.Insert(want); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	post, err := d.Posts.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.SubmittedAt)
	if !reflect.DeepEqual(post, want) {
		t.Errorf("got post %+v, want %+v", post, want)
	}
}

func TestPostsStore_List_db(t *testing.T) {
	want := []*thesrc.Post{{ID: 1}}

	// tx, _ := DB.Begin()
	// defer tx.Rollback()
	tx := DBH
	tx.Exec(`DELETE FROM post;`) // test on a clean DB
	if err := tx.Insert(want[0]); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	posts, err := d.Posts.List(&thesrc.PostListOptions{ListOptions: thesrc.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range want {
		normalizeTime(&p.SubmittedAt)
	}
	if !reflect.DeepEqual(posts, want) {
		t.Errorf("got posts %+v, want %+v", posts, want)
	}
}
