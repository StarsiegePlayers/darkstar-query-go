package darkstar_query_go

import (
	"github.com/StarsiegePlayers/darkstar-query-go/v2/query"
	"github.com/StarsiegePlayers/darkstar-query-go/v2/server"
)

func (q *Query) Masters() (output []*query.MasterQuery, games map[string]*server.Server, errorArray []error) {
	availableMasters := len(q.Addresses)
	await := make(chan *ServerResult)

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

	games = q.dedupeMasterQuery(output)

	return
}

func (q *Query) dedupeMasterQuery(servers []*query.MasterQuery) (output map[string]*server.Server) {
	output = make(map[string]*server.Server)
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
