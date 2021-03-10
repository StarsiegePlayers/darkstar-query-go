package darkstar_query_go

import (
	"net"
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/server"
)

func Servers(options Options) ([]*server.PingInfo, []error) {
	servers := options.Search
	availableServers := len(servers)
	await := make(chan *server.Query)
	errors := make([]error, 0)

	for id, game := range servers {
		conn, err := net.Dial("udp", game)
		if err != nil {
			errors = append(errors, err)
			availableServers--
			continue
		}
		go pingInfoQuery(conn, id, await, options.Timeout)
	}

	var output []*server.PingInfo
	for i := 0; i < availableServers; i++ {
		result := <-await
		if result.Error != nil {
			errors = append(errors, result.Error)
		} else {
			output = append(output, result.ServerInfo)
		}
	}

	close(await)

	return output, errors
}

func pingInfoQuery(conn net.Conn, id int, ret chan *server.Query, timeout time.Duration) {
	query := new(server.Query)
	query.ServerInfo = new(server.PingInfo)
	err := query.ServerInfo.PingInfoQuery(conn, id, timeout)
	if err != nil {
		query.Error = err
	}
	ret <- query
}
