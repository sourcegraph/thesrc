package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/sourcegraph/thesrc"
	"github.com/sourcegraph/thesrc/datastore"
)

func init() {
	serveMux.Handle("/", http.StripPrefix("/api", Handler()))
}

var (
	serveMux   = http.NewServeMux()
	httpClient = http.Client{Transport: (*muxTransport)(serveMux)}
	apiClient  = thesrc.NewClient(&httpClient)
)

func setup() {
	store = datastore.NewMockDatastore()
}

type muxTransport http.ServeMux

func (t *muxTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rw := httptest.NewRecorder()
	rw.Body = new(bytes.Buffer)
	(*http.ServeMux)(t).ServeHTTP(rw, req)
	return &http.Response{
		StatusCode:    rw.Code,
		Status:        http.StatusText(rw.Code),
		Header:        rw.HeaderMap,
		Body:          ioutil.NopCloser(rw.Body),
		ContentLength: int64(rw.Body.Len()),
		Request:       req,
	}, nil
}
