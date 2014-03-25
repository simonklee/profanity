// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"log"
	"os"

	"github.com/simonz05/util/raven"
)

type ravenLogger struct {
	l   *log.Logger
	sev Level
	dsn string
}

func (r *ravenLogger) Output(calldepth int, s string, sev Level) error {
	if r.sev < sev {
		return nil
	}

	if r.l == nil {
		r.init()
	}

	return r.l.Output(calldepth, s)
}

func (r *ravenLogger) init() {
	c, err := raven.NewClient(r.dsn, "")

	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}

	r.l = log.New(&ravenWriter{c: c}, "", log.Lshortfile)
}

type ravenWriter struct {
	c *raven.Client
}

func (w *ravenWriter) Write(p []byte) (int, error) {
	return len(p), w.c.Error(string(p))
}
