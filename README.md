# profanity — word filter service

[![Build Status](https://travis-ci.org/simonz05/profanity.png?branch=master)](https://travis-ci.org/simonz05/profanity)

## Filter

[wordfilter](http://godoc.org/github.com/simonz05/profanity/wordfilter)
is a simple library which implements a word filter.  The
library takes a list of words which are used to sanitize
text. The sanitizer will replace all words which match a
word in the list with **** (stars). 

## Profanity

`profanity` is a HTTP server which implements a simple API.
It exposes the `wordfilter.Wordfilter`. 

Usage:

    profanity [flag]

The flags are:

    -h
            help text
    -http=":8080"
            set bind address for the HTTP server
    -log=0
            set log level
    -redis="redis://:@localhost:6379/15"
            redis DSN
    -config=filename
            config filename
    -debug.cpuprofile=""
            run cpu profiler

### API

Create/overwrite blacklist.

    POST --data "blacklist=x&blacklist=xx&blacklist=xxx" /v1/profanity/blacklist/?lang=en_US

    HTTP/1.1 201 Created
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Return current blacklist.

    GET /v1/profanity/blacklist/?lang=en_US&count=10&offset=0

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Date: Mon, 12 Aug 2013 09:34:44 GMT
    Transfer-Encoding: chunked

    {"blacklist": ["x", "xx", "xxx"], "total": 3}

Update blacklist.

    PUT --data "blacklist=y" /v1/profanity/blacklist/?lang=en_US

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Remove from blacklist.

    PUT --data "blacklist=y" /v1/profanity/blacklist/remove/?lang=en_US

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Sanitize text.

    GET /v1/profanity/sanitize/?text=foo%20bar%20xxx&lang=en_US

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:38 GMT
    Content-Type: application/json; charset=utf-8
    Content-Length: 33

    {"text":"foo bar ***"}
