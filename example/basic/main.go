package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	query "github.com/StarsiegePlayers/darkstar-query-go"
	"github.com/StarsiegePlayers/darkstar-query-go/master"
	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
	"github.com/StarsiegePlayers/darkstar-query-go/server"
)

func main() {
	errors := make([]string, 0)
	masterQueryOptions := &protocol.Options{
		Search: protocol.NewServersMapFromList([]string{
			"master1.starsiegeplayers.com:29000",
			"master2.starsiegeplayers.com:29000",
			"master3.starsiegeplayers.com:29000",
			"starsiege1.no-ip.org:29000",
			"starsiege.noip.us:29000",
			"southerjustice.dyndns-server.com:29000",
			"dustersteve.ddns.net:29000",
			"starsiege.from-tx.com:29000",
		}),
		Timeout: 5 * time.Second,
		Debug:   true,
	}

	masterServerInfo, servers, errs := query.Masters(masterQueryOptions)
	if len(errs) >= 0 {
		for _, v := range errs {
			errors = append(errors, v.Error())
		}
	}

	masterStats(masterServerInfo)

	log.Printf("Acquired %d unique servers\n", len(servers))

	gameQueryOptions := &protocol.Options{
		Search:  servers,
		Timeout: 5 * time.Second,
		Debug:   true,
	}
	games, errs := query.Servers(gameQueryOptions)
	for _, err := range errs {
		errors = append(errors, err.Error())
	}

	output, _ := json.MarshalIndent(struct {
		RequestTime time.Time
		Masters     []*master.Master
		Games       []*server.PingInfo
		Errors      []string
	}{
		RequestTime: time.Now(),
		Masters:     masterServerInfo,
		Games:       games,
		Errors:      errors,
	}, "", "\t")

	fmt.Printf("%s", output)
}

func masterStats(servers []*master.Master) {
	for _, m := range servers {
		log.Printf("Master: %s [%s] returned %d servers in %s\n", m.Address, m.CommonName, len(m.Servers), m.Ping)
	}
}
