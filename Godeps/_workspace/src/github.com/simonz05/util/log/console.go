// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"log"
	"os"
)

type consoleLogger struct {
	l   *log.Logger
	sev Level
}

func (l *consoleLogger) Output(calldepth int, s string, sev Level) error {
	if l.sev < sev {
		return nil
	}

	if l.l == nil {
		l.l = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
	}

	return l.l.Output(calldepth, s)
}
