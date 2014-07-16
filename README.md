# thesrc [![Build Status](https://travis-ci.org/sourcegraph/thesrc.png?branch=master)](https://travis-ci.org/sourcegraph/thesrc)

## Installation

Use the `thesrc` command to interact with the app.

You can either run it directly:

```
go get github.com/sourcegraph/thesrc
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
thesrc createdb
thesrc import
thesrc classify

# then, in a persistent terminal window:
thesrc serve

# now open your browser to localhost:5000
```
