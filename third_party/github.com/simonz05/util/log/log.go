// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// pkg log implements a logger. It supports console, file and raven logging.
package log

import (
	"flag"
	"fmt"
	"os"
)

var (
	// Severity stores the log level
	Severity Level
	std      logger
	filename *string
	ravenDSN *string
)

func init() {
	flag.Var(&Severity, "log", "log level")
	filename = flag.String("log-file", "", "If non-empty, write log to this file")
	ravenDSN = flag.String("log-raven-dsn", "", "If non-empty, write to raven dsn")
	std = new(multiLogger)
}

type logger interface {
	Output(calldepth int, s string, sev Level) error
}

type multiLogger struct {
	loggers []logger
}

func (l multiLogger) Output(calldepth int, s string, sev Level) (err error) {
	if len(l.loggers) == 0 {
		l.init()
	}

	for _, w := range l.loggers {
		err = w.Output(calldepth, s, sev)

		if err != nil {
			return
		}
	}

	return
}

func (l *multiLogger) init() {
	l.loggers = append(l.loggers, &consoleLogger{sev: Severity})

	if *filename != "" {
		l.loggers = append(l.loggers, &fileLogger{fname: *filename, sev: Severity})
	}

	if *ravenDSN != "" {
		l.loggers = append(l.loggers, &ravenLogger{dsn: *ravenDSN, sev: LevelError})
	}
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	if Severity >= LevelInfo {
		std.Output(5, fmt.Sprint(v...), LevelInfo)
	}
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	if Severity >= LevelInfo {
		std.Output(5, fmt.Sprintf(format, v...), LevelInfo)
	}
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	if Severity >= LevelInfo {
		std.Output(5, fmt.Sprintln(v...), LevelInfo)
	}
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Error(v ...interface{}) {
	if Severity >= LevelError {
		std.Output(5, fmt.Sprint(v...), LevelError)
	}
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Errorf(format string, v ...interface{}) {
	if Severity >= LevelError {
		std.Output(5, fmt.Sprintf(format, v...), LevelError)
	}
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Errorln(v ...interface{}) {
	if Severity >= LevelError {
		std.Output(5, fmt.Sprintln(v...), LevelError)
	}
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	if Severity >= LevelFatal {
		std.Output(5, fmt.Sprint(v...), LevelFatal)
		os.Exit(1)
	}
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	if Severity >= LevelFatal {
		std.Output(5, fmt.Sprintf(format, v...), LevelFatal)
		os.Exit(1)
	}
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	if Severity >= LevelFatal {
		std.Output(5, fmt.Sprintln(v...), LevelFatal)
		os.Exit(1)
	}
}
