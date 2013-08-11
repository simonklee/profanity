package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/simonz05/profanity/server"
)

var (
	help       = flag.Bool("h", false, "this help")
	laddr      = flag.String("a", ":8080", "bind address")
	filename   = flag.String("filename", "", "filename which contains db")
	logLevel   = flag.Int("l", 0, "set logging level")
	version    = flag.Bool("v", false, "show version and exit")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Fprintln(os.Stderr, server.Version)
		return
	}

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if *laddr == "" {
		fmt.Fprintln(os.Stderr, "listen address required")
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	server.LogLevel = *logLevel

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err := server.ListenAndServe(*laddr, *filename)

	if err != nil {
		log.Println(err)
	}
}
