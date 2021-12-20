package darkstar_query_go

import (
	"github.com/StarsiegePlayers/darkstar-query-go/v2/query"
)

func (q *Query) Servers() ([]*query.PingInfoQuery, []error) {
	availableServers := len(q.Addresses)
	await := make(chan *ServerResult)
	errors := make([]error, 0)

	for _, game := range q.Addresses {
		go q.performServerQuery(game, await)
	}

	var output []*query.PingInfoQuery
	for i := 0; i < availableServers; i++ {
		result := <-await
		if result.Error != nil {
			errors = append(errors, result.Error)
		} else {
			output = append(output, result.Game)
		}
	}

	close(await)

	return output, errors
}

func (q *Query) performServerQuery(address string, ret chan *ServerResult) {
	r := new(ServerResult)

	r.Game = query.NewPingInfoQueryWithOptions(address, q.Options)

	err := r.Game.Query()
	if err != nil {
		r.Error = err
	}

	ret <- r
}
