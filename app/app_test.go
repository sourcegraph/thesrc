package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var testMux *http.ServeMux

func init() {
	LoadTemplates()
}

func setup() {
	testMux = http.NewServeMux()
	testMux.Handle("/", Handler())
	APIClient = nil
}

func teardown() {
	APIClient, testMux = nil, nil
}

func getHTML(t *testing.T, uri *url.URL) (*goquery.Document, *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rw := httptest.NewRecorder()
	rw.Body = new(bytes.Buffer)
	testMux.ServeHTTP(rw, req)

	doc, err := goquery.NewDocumentFromReader(rw.Body)
	if err != nil {
		t.Fatal(err)
	}
	return doc, rw
}
