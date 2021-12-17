package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	query "github.com/StarsiegePlayers/darkstar-query-go"
	"github.com/StarsiegePlayers/darkstar-query-go/master"
	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
	"github.com/StarsiegePlayers/darkstar-query-go/server"
)

const ServerStatsPathFormat = "2006/01/02"

type ServerListData struct {
	RequestTime time.Time
	Masters     []*master.Master
	Games       []*server.PingInfo
	Errors      []string
}

func main() {
	data := performServerListUpdate()
	masterStats(data)
	recordServerListUpdate(data)
}

func performServerListUpdate() ServerListData {
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
	}

	masterServerInfo, gameAddresses, errs := query.Masters(masterQueryOptions)
	if len(errs) >= 0 {
		for _, v := range errs {
			errors = append(errors, v.Error())
		}
	}

	serverQueryOptions := &protocol.Options{
		Search:  gameAddresses,
		Timeout: 5 * time.Second,
	}
	games, errs := query.Servers(serverQueryOptions)
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
