package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PagerDuty/godspeed"
	log "github.com/Sirupsen/logrus"
	r "gopkg.in/dancannon/gorethink.v2"
)

type RethinkStats struct {
	addr    []string
	tags    []string
	verbose bool
	stats   chan []Stat
	g       *godspeed.Godspeed
}

func NewRethinkStats(addr, tags, env string, verbose bool) *RethinkStats {
	rs := &RethinkStats{
		addr:    strings.Split(addr, ","),
		tags:    strings.Split(tags, ","),
		stats:   make(chan []Stat),
		verbose: verbose,
	}

	// DogStatsD
	datadogHost := os.Getenv("DATADOG")
	if datadogHost == "" {
		datadogHost = "datadog"
	}
	if env == "" || env == "dev" {
		datadogHost = "127.0.0.1"
	}
	g, _ := godspeed.New(datadogHost, 8125, false)
	defer g.Conn.Close()
	g.AddTags(rs.tags)

	rs.g = g
	go rs.procStats()

	return rs
}

func (rs *RethinkStats) Query() {
	session, err := r.Connect(r.ConnectOpts{
		Addresses: rs.addr,
		Database:  "rethinkdb",
	})
	defer session.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to open database connection")
		return
	}

	res, err := r.Table("stats").Run(session)
	defer res.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to query stats table")
		return
	}

	var stats []Stat
	if res.All(&stats) != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to retrieve results")
		return
	}

	rs.stats <- stats
}

func (rs *RethinkStats) procStats() {
	var err error
	for stats := range rs.stats {
		if rs.verbose {
			log.Println("Trying to send stats to DataDog..")
		}
		for _, stat := range stats {
			switch stat.ID[0] {
			case "cluster":
				{
					err = rs.g.Gauge("rethinkdb.read_docs_per_sec.cluster", stat.QueryEngine.ReadDocsPerSec, nil)
					err = rs.g.Gauge("rethinkdb.written_docs_per_sec.cluster", stat.QueryEngine.WrittenDocsPerSec, nil)
				}
			case "server":
				{
					err = rs.g.Gauge(
						fmt.Sprintf("rethinkdb.read_docs_per_sec.server.%s", stat.Server),
						stat.QueryEngine.ReadDocsPerSec,
						nil)
					err = rs.g.Gauge(
						fmt.Sprintf("rethinkdb.written_docs_per_sec.server.%s", stat.Server),
						stat.QueryEngine.WrittenDocsPerSec,
						nil)
				}
			case "table":
				{
					err = rs.g.Gauge(
						fmt.Sprintf("rethinkdb.read_docs_per_sec.table.%s", stat.Table),
						stat.QueryEngine.ReadDocsPerSec,
						nil)
					err = rs.g.Gauge(
						fmt.Sprintf("rethinkdb.written_docs_per_sec.table.%s", stat.Table),
						stat.QueryEngine.WrittenDocsPerSec,
						nil)
				}
			case "table_server":
				{
					// replicas
				}
			}
		}
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to send stats to DataDog")
		}
	}
}
