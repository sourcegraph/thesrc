package app

import (
	"fmt"
	htmpl "html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

var (
	// TemplateDir is the directory containing the html/template template files.
	TemplateDir string
)

func LoadTemplates() {
	err := parseHTMLTemplates([][]string{
		{"posts/show.html", "common.html", "layout.html"},
		{"posts/list.html", "common.html", "layout.html"},
		{"error.html", "common.html", "layout.html"},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// templateCommon is data that is passed to (and available to) all templates.
type templateCommon struct {
	CurrentURL         *url.URL
	PageGenerationTime time.Duration
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, status int, data interface{}) error {
	w.WriteHeader(status)
	if ct := w.Header().Get("content-type"); ct == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	t := templates[name]
	if t == nil {
		return fmt.Errorf("Template %s not found", name)
	}
	return t.Execute(w, data)
}

var templates = map[string]*htmpl.Template{}

func parseHTMLTemplates(sets [][]string) error {
	for _, set := range sets {
		t := htmpl.New("")
		t.Funcs(htmpl.FuncMap{})

		_, err := t.ParseFiles(joinTemplateDir(TemplateDir, set)...)
		if err != nil {
			return fmt.Errorf("template %v: %s", set, err)
		}

		t = t.Lookup("ROOT")
		if t == nil {
			return fmt.Errorf("ROOT template not found in %v", set)
		}
		templates[set[0]] = t
	}
	return nil
}

func joinTemplateDir(base string, files []string) []string {
	result := make([]string, len(files))
	for i := range files {
		result[i] = filepath.Join(base, files[i])
	}
	return result
}
