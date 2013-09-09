package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/simonz05/profanity/server"
	"github.com/simonz05/profanity/util"
)

var (
	verbose  = flag.Bool("v", false, "verbose mode")
	help     = flag.Bool("h", false, "show help text")
	laddr    = flag.String("http", ":8080", "set bind address for the HTTP server")
	wordlist = flag.String("wordlist", "", `filepath to use a '\\n' separated word list 
		which will be used as the default profanity filter`)
	logLevel   = flag.Int("log", 0, "set log level")
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

	if *verbose {
		util.LogLevel = 2
	} else {
		util.LogLevel = *logLevel
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err := server.ListenAndServe(*laddr, *wordlist)

	if err != nil {
		log.Println(err)
	}
}
