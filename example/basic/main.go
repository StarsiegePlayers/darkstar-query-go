package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	darkstar "github.com/StarsiegePlayers/darkstar-query-go/v2"
	"github.com/StarsiegePlayers/darkstar-query-go/v2/query"
)

type ServerListData struct {
	RequestTime time.Time
	Masters     []*query.MasterQuery
	Games       []*query.PingInfoQuery
	Errors      []string
}

const (
	timeout = 5 * time.Second
	debug   = true
)

func main() {
	data := performServerListUpdate()
	masterStats(data)
	printOut(data)
}

func printOut(data ServerListData) {
	output, _ := json.MarshalIndent(data, "", "\t")

	fmt.Println(string(output))
}

func performServerListUpdate() ServerListData {
	errors := make([]string, 0)
	q := darkstar.NewQuery(timeout, debug)
	q.Addresses = []string{
		"master1.starsiegeplayers.com:29000",
		"master2.starsiegeplayers.com:29000",
		"master3.starsiegeplayers.com:29000",
		"starsiege1.no-ip.org:29000",
		"starsiege.noip.us:29000",
		"southerjustice.dyndns-server.com:29000",
		"dustersteve.ddns.net:29000",
		"starsiege.from-tx.com:29000",
	}

	masterServerInfo, gameAddresses, errs := q.Masters()
	if len(errs) >= 0 {
		for _, v := range errs {
			errors = append(errors, v.Error())
		}
	}

	q = darkstar.NewQuery(timeout, debug)
	q.Addresses = []string{}
	for k := range gameAddresses {
		q.Addresses = append(q.Addresses, k)
	}

	games, errs := q.Servers()
	for _, err := range errs {
		errors = append(errors, err.Error())
	}

	return ServerListData{
		RequestTime: time.Now(),
		Masters:     masterServerInfo,
		Games:       games,
		Errors:      errors,
	}
}

func masterStats(data ServerListData) {
	for _, m := range data.Masters {
		log.Printf("Master: %s [%s] returned %d servers in %s\n", m.Address, m.CommonName, len(m.Servers), m.Ping)
	}
}
