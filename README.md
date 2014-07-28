# thesrc [![Build Status](https://travis-ci.org/sourcegraph/thesrc.png?branch=master)](https://travis-ci.org/sourcegraph/thesrc) [![docs examples](https://sourcegraph.com/api/repos/github.com/sourcegraph/thesrc/.badges/docs-examples.png)](https://sourcegraph.com/github.com/sourcegraph/thesrc) [![status](https://sourcegraph.com/api/repos/github.com/sourcegraph/thesrc/.badges/status.png)](https://sourcegraph.com/github.com/sourcegraph/thesrc) [![views](https://sourcegraph.com/api/repos/github.com/sourcegraph/thesrc/.counters/views.png)](https://sourcegraph.com/github.com/sourcegraph/thesrc)

<img width=300 align=right src="https://s3-us-west-2.amazonaws.com/sourcegraph-assets/thesrc-screenshot.png">

thesrc is a news site for programmers that's intended to be an example of how to
structure a large Go web app. While this app is not large itself, it
demonstrates the same patterns as in the web app that powers
[Sourcegraph.com](https://sourcegraph.com).

The web application architecture and patterns demonstrated here were presented
in a talk at Google I/O 2014 entitled *[Building Sourcegraph, a large-scale code
search engine in Go](https://sourcegraph.com/blog/google-io-2014-building-sourcegraph-a-large-scale-code-search-engine-in-go)*. See that talk for more details.

thesrc has a few special features of interest to programmers:

* just the good stuff: an automated classifier rejects links that don't contain code or involve programming;
* not a popularity contest: you can only see a link's score by mousing over it for a couple of seconds, and (TODO) freshly posted links are randomly rotated into the homepage;

**Browse the code on [Sourcegraph](https://sourcegraph.com/github.com/sourcegraph/thesrc).**

## Installation

Use the `thesrc` command to interact with the app.

You can either run it directly:

```
go get github.com/sourcegraph/thesrc/...
go install github.com/sourcegraph/thesrc/...
thesrc
```

Or inside Docker:

```
docker build -t thesrc && docker run thesrc
```

If you want to run it in Docker, substitute `docker run thesrc` for every
instance of `thesrc`. (Also note that you'll have to pass Docker the necessary
`PG*` environment variables to connect to the PostgreSQL database.)

## Running

First, set the `PG*` environment variables so that `psql` works.

Then run these commands to create the DB, import posts from other sites, and classify their links:

```
# start the server:
thesrc serve

# then, in a separate terminal window, run:
thesrc -url=http://localhost:5000 createdb
thesrc -url=http://localhost:5000 import
thesrc -url=http://localhost:5000 classify

# now open your browser to localhost:5000
```
