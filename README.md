# profanity — profanity filter service 

[![Build Status](https://travis-ci.org/simonz05/profanity.png?branch=master)](https://travis-ci.org/simonz05/profanity)

## Filter

`filter` is a simple library which implements a word filter.
The library takes a list of words which are used to sanitize
text. The sanitizer will replace all words which match a
word in the list with **** (stars). 

## Profanity

`profanity` is a HTTP server which implements a simple API.
It exposes the `filter.Update()`, `filter.Replace` and
`filter.Sanitize()` methods. 

### API

Create/overwrite blacklist.

    POST --data "blacklist=x&blacklist=xx&blacklist=xxx" /api/1.0/blacklist/

    HTTP/1.1 201 Created
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Return current blacklist.

    GET /api/1.0/blacklist/

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Date: Mon, 12 Aug 2013 09:34:44 GMT
    Transfer-Encoding: chunked

    ["x", "xx", "xxx"]    

Update blacklist.

    PUT --data "blacklist=y" /api/1.0/blacklist/

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:17 GMT
    Content-Length: 0
    Content-Type: text/plain; charset=utf-8

Sanitize text.

    GET /api/1.0/sanitize/?text=foo%20bar%20xxx

    HTTP/1.1 200 OK
    Date: Mon, 12 Aug 2013 09:37:38 GMT
    Content-Type: application/json; charset=utf-8
    Content-Length: 33

    {"text":"foo bar ***","lang":""}
