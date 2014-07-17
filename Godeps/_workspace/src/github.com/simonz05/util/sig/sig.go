// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package util/sig is simple sig trap closer

package sig

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simonz05/util/log"
)

// TODO remove Cleanup
type Cleanup func() error

func TrapCloser(cl io.Closer, cleanups ...Cleanup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for {
			rawSig := <-c
			sig, ok := rawSig.(syscall.Signal)

			if !ok {
				log.Fatal("not a unix signal")
			}

			switch {
			case sig == syscall.SIGHUP:
				log.Print("SIGHUP: restart process")
				err := RestartProcess()

				if err != nil {
					log.Fatal("failed to restart: " + err.Error())
				}
			case sig == syscall.SIGINT || sig == syscall.SIGKILL || sig == syscall.SIGTERM:
				log.Print("shutting down")
				donec := make(chan bool)
				go func() {
					for _, cb := range cleanups {
						if err := cb(); err != nil {
							log.Error(err)
						}
					}

					if err := cl.Close(); err != nil {
						log.Fatalf("Error shutting down: %v", err)
					}

					donec <- true
				}()
				select {
				case <-donec:
					log.Printf("shutdown")
					os.Exit(0)
				case <-time.After(5 * time.Second):
					log.Fatal("Timeout shutting down. Exiting uncleanly.")
				}
			default:
				log.Fatal("Received another signal, should not happen.")
			}
		}
	}()
}
