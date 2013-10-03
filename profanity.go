package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/simonz05/profanity/server"
	"github.com/simonz05/util/log"
)

var (
	help       = flag.Bool("h", false, "show help text")
	laddr      = flag.String("http", ":6061", "set bind address for the HTTP server")
	dsn        = flag.String("redis", "redis://:@localhost:6379/15", "Redis data source name")
	version    = flag.Bool("version", false, "show version number and exit")
	cpuprofile = flag.String("debug.cpuprofile", "", "write cpu profile to file")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.Println("Start")

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

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err := server.ListenAndServe(*laddr, *dsn)

	if err != nil {
		log.Errorln(err)
	}
}
