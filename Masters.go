package darkstar_query_go

import (
	"github.com/StarsiegePlayers/darkstar-query-go/v2/query"
	"github.com/StarsiegePlayers/darkstar-query-go/v2/server"
)

func (q *Query) Masters() ([]*query.MasterQuery, map[string]*server.Server, []error) {
	availableMasters := len(q.Addresses)
	await := make(chan *ServerResult)
	output := make([]*query.MasterQuery, 0)
	errorArray := make([]error, 0)

	for _, address := range q.Addresses {
		go q.performMasterQuery(address, await)
	}

	for i := 0; i < availableMasters; i++ {
		result := <-await
		if result.Master.Ping > 0 {
			output = append(output, result.Master)
		}
		if result.Error != nil {
			errorArray = append(errorArray, result.Error)
		}
	}

	close(await)

	games := q.dedupeMasterQuery(output)

	return output, games, errorArray
}

func (q *Query) dedupeMasterQuery(servers []*query.MasterQuery) map[string]*server.Server {
	output := make(map[string]*server.Server)
	for _, svr := range servers {
		for k, v := range svr.Servers {
			if _, ok := output[k]; ok {
				continue
			}
			output[k] = v
		}
	}

	return output
}

func (q *Query) performMasterQuery(address string, ret chan *ServerResult) {
	r := new(ServerResult)

	r.Master = query.NewMasterQueryWithOptions(address, q.Options)
	r.Error = r.Master.Query()

	ret <- r
}
