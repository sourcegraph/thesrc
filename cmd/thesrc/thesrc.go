package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go/build"

	"github.com/sourcegraph/thesrc/api"
	"github.com/sourcegraph/thesrc/app"
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
	{"serve", "start web server", serveCmd},
}

func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	httpAddr := flag.String("http", ":5000", "HTTP service address")
	templateDir := fs.String("tmpl-dir", filepath.Join(defaultBase("github.com/sourcegraph/thesrc/app"), "tmpl"), "template directory")
	staticDir := fs.String("static-dir", filepath.Join(defaultBase("github.com/sourcegraph/thesrc/app"), "static"), "static assets directory")
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

	m := http.NewServeMux()
	m.Handle("/api", api.Handler())
	m.Handle("/", app.Handler())

	log.Print("Listening on ", *httpAddr)
	err := http.ListenAndServe(*httpAddr, m)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func defaultBase(path string) string {
	p, err := build.Default.Import(path, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}
