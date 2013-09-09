package util

import (
	"log"
	"os"
)

var (
	LogLevel int = 0
	Logger       = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
)

func Logf(fmt string, args ...interface{}) {
	if LogLevel == 0 {
		return
	}

	log.Printf(fmt, args...)
}

func Logln(args ...interface{}) {
	if LogLevel == 0 {
		return
	}

	log.Println(args...)
}
