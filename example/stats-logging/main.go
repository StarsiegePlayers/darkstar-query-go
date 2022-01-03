package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/v2"
	"github.com/StarsiegePlayers/darkstar-query-go/v2/query"
)

const ServerStatsPathFormat = "2006/01/02"

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
	recordServerListUpdate(data)
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
	for _, v := range errs {
		errors = append(errors, v.Error())
	}

	q = darkstar.NewQuery(timeout, debug)
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

func recordServerListUpdate(data ServerListData) {
	for _, s := range data.Games {
		currentTime := time.Now()

		directory := fmt.Sprintf("./stats/%s", currentTime.Format(ServerStatsPathFormat))

		err := os.MkdirAll(directory, 755)
		if err != nil {
			log.Println("Error", err.Error())
		}

		fileName := strings.Replace(s.Address, ":", "_", 1)
		filePath := fmt.Sprintf("%s/%s.csv", directory, fileName)

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Error", err.Error())
			continue
		}

		output := fmt.Sprintf("\"%s\",\"%s\",%s,%s,%d,%d,%s", currentTime.Format(time.RFC1123Z), s.Name, s.Ping, s.GameStatus, s.PlayerCount, s.MaxPlayers, s.Address)

		_, err = f.WriteString(output)
		if err != nil {
			log.Println("error", err.Error())

			_ = f.Close()

			continue
		}

		log.Printf("Wrote file %s\n", filePath)

		_ = f.Close()
	}
}
