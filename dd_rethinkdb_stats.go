package main

import (
	"strings"

	//  "github.com/PagerDuty/godspeed"
	log "github.com/Sirupsen/logrus"
	r "gopkg.in/dancannon/gorethink.v2"
)

type RethinkStats struct {
	addr []string
	tags []string
	proc chan []Stat
}

func NewRethinkStats(addr, tags string) *RethinkStats {
	rs := &RethinkStats{
		addr: strings.Split(addr, ","),
		tags: strings.Split(tags, ","),
		proc: make(chan []Stat),
	}

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
		}).Error("failed to open database connection")
		return
	}

	res, err := r.Table("stats").Run(session)
	defer res.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to query stats table")
		return
	}

	var stats []Stat
	if res.All(&stats) != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to retrieve results")
		return
	}

	rs.proc <- stats
}

func (rs *RethinkStats) procStats() {
	for stats := range rs.proc {

		countServers := 0
		countTables := 0
		countReplicas := 0

		for _, stat := range stats {
			//log.Printf("%+v", stat)

			switch stat.ID[0] {
			case "server":
				{
					countServers++
				}
			case "table":
				{
					countTables++
				}
			case "table_server":
				{
					countReplicas++
				}
			}
		}

		log.Printf("%d, %d, %d", countServers, countTables, countReplicas)
	}
}
