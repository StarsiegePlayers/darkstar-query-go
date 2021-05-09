package darkstar_query_go

import (
	"errors"
	"fmt"
	"net"

	"github.com/StarsiegePlayers/darkstar-query-go/master"
	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
)

func Masters(options protocol.Options) ([]*master.Master, []string, []error) {
	masters := options.Search
	availableMasters := len(masters)
	await := make(chan *master.Query)
	output := make([]*master.Master, 0)
	errorArray := make([]error, 0)

	for id, server := range masters {
		conn, err := net.Dial("udp", server)
		if err != nil {
			var dnsError *net.DNSError
			if errors.As(err, &dnsError) {
				errorArray = append(errorArray, fmt.Errorf("master: [%s]: no such host", dnsError.Name))
			}

			availableMasters--
			continue
		}
		go masterQuery(conn, server, id, await, options)
	}

	for i := 0; i < availableMasters; i++ {
		result := <-await
		if result.MasterData.Ping > 0 {
			output = append(output, result.MasterData)
		}
		if result.Error != nil {
			errorArray = append(errorArray, result.Error)
		}
	}

	close(await)

	games := dedupe(output)

	return output, games, errorArray
}

func dedupe(servers []*master.Master) []string {
	games := make(map[string]bool)
	for _, server := range servers {
		for _, game := range server.ServerAddresses {
			games[game] = true
		}
	}

	output := make([]string, len(games))
	i := 0
	for server := range games {
		output[i] = server
		i++
	}
	return output
}

func masterQuery(conn net.Conn, hostname string, id int, ret chan *master.Query, options protocol.Options) {
	query := new(master.Query)
	query.MasterData = new(master.Master)
	query.MasterData.Address = hostname
	query.Error = query.MasterData.Query(conn, id, options)
	ret <- query
}
