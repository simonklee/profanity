// Copyright (c) 2013 Simon Zimmermann
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// `profanity` is a HTTP server which implements a simple API.
// It exposes the `wordfilter.Wordfilter`.
//
//    Usage:
//    	profanity [flag]
//
//    The flags are:
//    	-v
//    		verbose mode
//    	-h
//    		help text
//    	-http=":8080"
//    		set bind address for the HTTP server
//    	-log=0
//		set log level
//    	-version=false
//    		display version number and exit
//    	-debug.cpuprofile=""
//    		run cpu profiler
//
// API
//
// Create/overwrite blacklist.
//
//     POST --data "blacklist=x&blacklist=xx&blacklist=xxx" /api/1.0/blacklist/?lang=en_US
//
//     HTTP/1.1 201 Created
//     Date: Mon, 12 Aug 2013 09:37:17 GMT
//     Content-Length: 0
//     Content-Type: text/plain; charset=utf-8
//
// Return current blacklist.
//
//     GET /api/1.0/blacklist/?lang=en_US&count=10&offset=0
//
//     HTTP/1.1 200 OK
//     Content-Type: application/json; charset=utf-8
//     Date: Mon, 12 Aug 2013 09:34:44 GMT
//     Transfer-Encoding: chunked
//
//     {"blacklist": ["x", "xx", "xxx"], "total": 3}
//
// Update blacklist.
//
//     PUT --data "blacklist=y" /api/1.0/blacklist/?lang=en_US
//
//     HTTP/1.1 200 OK
//     Date: Mon, 12 Aug 2013 09:37:17 GMT
//     Content-Length: 0
//     Content-Type: text/plain; charset=utf-8
//
// Remove from blacklist.
//
//     PUT --data "blacklist=y" /api/1.0/blacklist/remove/?lang=en_US
//
//     HTTP/1.1 200 OK
//     Date: Mon, 12 Aug 2013 09:37:17 GMT
//     Content-Length: 0
//     Content-Type: text/plain; charset=utf-8
//
// Sanitize text.
//
//     GET /api/1.0/sanitize/?text=foo%20bar%20xxx&lang=en_US
//
//     HTTP/1.1 200 OK
//     Date: Mon, 12 Aug 2013 09:37:38 GMT
//     Content-Type: application/json; charset=utf-8
//     Content-Length: 33
//
//     {"text":"foo bar ***"}
package main
