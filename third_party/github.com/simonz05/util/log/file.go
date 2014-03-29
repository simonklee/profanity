// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"log"
	"os"
	"path/filepath"
)

type fileLogger struct {
	l     *log.Logger
	sev   Level
	fname string
}

func (l *fileLogger) Output(calldepth int, s string, sev Level) error {
	if l.sev < sev {
		return nil
	}

	if l.l == nil {
		l.init()
	}

	return l.l.Output(calldepth, s)
}

func (l *fileLogger) init() {
	f, err := os.OpenFile(filepath.Clean(l.fname), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}

	l.l = log.New(f, "", log.Ldate|log.Lmicroseconds)
}
