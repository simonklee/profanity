# profanity — word filter service

[![Build Status](https://travis-ci.org/simonz05/profanity.png?branch=master)](https://travis-ci.org/simonz05/profanity)

## Filter

[filter](http://godoc.org/github.com/simonz05/profanity/filter) is a
simple library which implements a word filter.  The library
takes a list of words which are used to sanitize text. The
sanitizer will replace all words which match a word in the
list with **** (stars). 

## Profanity

`profanity` is a HTTP server which implements a simple API.
It exposes the `filter.Update()`, `filter.Replace()`,
`filter.Remove()` and `filter.Sanitize()` methods. 

Usage:

    profanity [flag]

The flags are:

    -v
            verbose mode
    -h
            help text
    -http=":8080"
            set bind address for the HTTP server
    -wordlist=""
            filepath to use a '\n' separated word list which
            will be used as the default profanity filter
    -log=0
            set log level
    -version=false
            display version number and exit
    -debug.cpuprofile=""
            run cpu profiler

### API

Create/overwrite blacklist.

    POST --data "blacklist=x&blacklist=xx&blacklist=xxx" /api/1.0/blacklist/?lang=en_US

    HTTP/1.1 201 Created
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Return current blacklist.

    GET /api/1.0/blacklist/?lang=en_US

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Date: Mon, 12 Aug 2013 09:34:44 GMT
    Transfer-Encoding: chunked

    {"blacklist": ["x", "xx", "xxx"], "total": 3}

Update blacklist.

    PUT --data "blacklist=y" /api/1.0/blacklist/?lang=en_US

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Remove from blacklist.

    PUT --data "blacklist=y" /api/1.0/blacklist/remove/?lang=en_US

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Sanitize text.

    GET /api/1.0/sanitize/?text=foo%20bar%20xxx&lang=en_US

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:38 GMT
    Content-Type: application/json; charset=utf-8
    Content-Length: 33

    {"text":"foo bar ***"}
