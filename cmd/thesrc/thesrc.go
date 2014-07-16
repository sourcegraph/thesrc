package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/sourcegraph/thesrc"
	"github.com/sourcegraph/thesrc/api"
	"github.com/sourcegraph/thesrc/app"
	"github.com/sourcegraph/thesrc/datastore"
	"github.com/sourcegraph/thesrc/importer"
	"github.com/sourcegraph/thesrc/router"
)

var (
	baseURLStr = flag.String("url", "http://thesrc.org", "base URL of thesrc")
	baseURL    *url.URL
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `thesrc is a web news and link server.

Usage:

        thesrc [options] command [arg...]

The commands are:
`)
		for _, c := range subcmds {
			fmt.Fprintf(os.Stderr, "    %-24s %s\n", c.name, c.description)
		}
		fmt.Fprintln(os.Stderr, `
Use "thesrc command -h" for more information about a command.

The options are:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}
	log.SetFlags(0)

	var err error
	baseURL, err = url.Parse(*baseURLStr)
	if err != nil {
		log.Fatal(err)
	}

	subcmd := flag.Arg(0)
	for _, c := range subcmds {
		if c.name == subcmd {
			c.run(flag.Args()[1:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcmd %q\n", subcmd)
	fmt.Fprintln(os.Stderr, `Run "thesrc -h" for usage.`)
	os.Exit(1)
}

type subcmd struct {
	name        string
	description string
	run         func(args []string)
}

var subcmds = []subcmd{
	{"post", "submit a post", postCmd},
	{"import", "import posts from other sites", importCmd},
	{"serve", "start web server", serveCmd},
	{"createdb", "create the database schema", createDBCmd},
}

var apiclient = thesrc.NewClient(nil)

func postCmd(args []string) {
	fs := flag.NewFlagSet("post", flag.ExitOnError)
	title := fs.String("title", "", "title of post")
	linkURL := fs.String("link", "", "link URL")
	body := fs.String("body", "", "body of post")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: thesrc post [options]

Submits a post.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	if *title == "" {
		log.Fatal(`Title must not be empty. See "thesrc post -h" for usage.`)
	}
	if *linkURL == "" {
		log.Fatal(`Link URL must not be empty. See "thesrc post -h" for usage.`)
	}

	post := &thesrc.Post{
		Title:   *title,
		LinkURL: *linkURL,
		Body:    *body,
	}
	err := apiclient.Posts.Create(post)
	if err != nil {
		log.Fatal(err)
	}

	url, err := router.App().Get(router.Post).URL("ID", strconv.Itoa(post.ID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(baseURL.ResolveReference(url))
}

func importCmd(args []string) {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	hnTop := fs.Bool("hn-top", true, "import news.ycombinator.com/news")
	hnNewest := fs.Bool("hn-newest", true, "import news.ycombinator.com/newest")
	hnBest := fs.Bool("hn-best", true, "import from news.ycombinator.com/best")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: thesrc import [options]

Imports posts from other sites.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	n := 0
	var printMu sync.Mutex
	importer.Imported = func(site string, post *thesrc.Post) {
		printMu.Lock()
		defer printMu.Unlock()
		fmt.Printf("%-12s  %-50s\n              %-60s\n", site, post.Title, post.LinkURL)
		n++
	}

	var fetchers []importer.Fetcher
	if *hnTop {
		fetchers = append(fetchers, importer.HackerNewsTop)
	}
	if *hnNewest {
		fetchers = append(fetchers, importer.HackerNewsNewest)
	}
	if *hnBest {
		fetchers = append(fetchers, importer.HackerNewsBest)
	}

	datastore.Connect()
	var failed bool
	var wg sync.WaitGroup
	for _, f_ := range fetchers {
		f := f_
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := importer.Import(f); err != nil {
				log.Printf("Error fetching from %s: %s.", f.Site(), err)
				failed = true
			}
		}()
	}
	wg.Wait()

	log.Printf("# %d posts imported", n)
	if failed {
		os.Exit(1)
	}
}

func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	httpAddr := flag.String("http", ":5000", "HTTP service address")
	templateDir := fs.String("tmpl-dir", app.TemplateDir, "template directory")
	staticDir := fs.String("static-dir", app.StaticDir, "static assets directory")
	reload := flag.Bool("reload", true, "reload templates on each request (dev mode)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: thesrc serve [options] 

Starts the web server that serves the app and API.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	app.StaticDir = *staticDir
	app.TemplateDir = *templateDir
	app.ReloadTemplates = *reload
	app.LoadTemplates()

	datastore.Connect()

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	m.Handle("/", app.Handler())

	log.Print("Listening on ", *httpAddr)
	err := http.ListenAndServe(*httpAddr, m)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func createDBCmd(args []string) {
	fs := flag.NewFlagSet("createdb", flag.ExitOnError)
	drop := fs.Bool("drop", false, "drop DB before creating")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: thesrc createdb [options] 

Creates the necessary DB tables and indexes.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	datastore.Connect()
	if *drop {
		datastore.Drop()
	}
	datastore.Create()
}
