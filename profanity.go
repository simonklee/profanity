package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/simonz05/profanity/config"
	"github.com/simonz05/profanity/server"
	"github.com/simonz05/profanity/types"
	"github.com/simonz05/util/log"
)

var (
	help           = flag.Bool("h", false, "show help text")
	laddr          = flag.String("http", ":6061", "set bind address for the HTTP server")
	dsn            = flag.String("redis", "redis://:@localhost:6379/15", "Redis data source name")
	filterType     = flag.String("filter", "", "filter type")
	configFilename = flag.String("config", "config.toml", "config file path")
	cpuprofile     = flag.String("debug.cpuprofile", "", "write cpu profile to file")
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

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	conf, err := config.ReadFile(*configFilename)

	if err != nil {
		log.Fatal(err)
	}

	if conf.Listen == "" && *laddr == "" {
		log.Fatal("Listen address required")
	} else if conf.Listen == "" {
		conf.Listen = *laddr
	}

	if conf.Redis.DSN == "" {
		conf.Redis.DSN = *dsn
	}

	if *filterType != "" {
		switch strings.ToLower(*filterType) {
		case "any":
			conf.Filter = types.Any
		default:
			conf.Filter = types.Word
		}
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

	err = server.ListenAndServe(conf)

	if err != nil {
		log.Println(err)
	}
}
