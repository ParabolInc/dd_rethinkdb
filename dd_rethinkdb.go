package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	env     = os.Getenv("DD_RETHINKDB_ENV")
	addr    = flag.String("addr", "localhost:28015", "Database cluster address, comma separated")
	tags    = flag.String("tags", "rethinkdb", "Tags to associate with metrics, comma separated")
	tick    = flag.Duration("tick", 15*time.Second, "Statistics check interval")
	verbose = flag.Bool("verbose", false, "Enable verbose logging")
)

func init() {
	if env == "" || env == "dev" {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	flag.Usage = func() {
		fmt.Println(NameVersion())
		fmt.Println()
		fmt.Println("usage: dd_rethinkdb [options]")
		fmt.Println()
		fmt.Println("options:")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	flag.Parse()

	if len(*addr) == 0 {
		flag.Usage()
	}

	log.Println(NameVersion())

	sigs := make(chan os.Signal)
	done := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	stats := NewRethinkStats(*addr, *tags, env, *verbose)
	stats.Query()

	ticker := time.NewTicker(*tick)
	for {
		select {
		case <-ticker.C:
			stats.Query()
		case <-done:
			ticker.Stop()
			stats.Close()
			return
		}
	}
}
