package darkstar_query_go

import (
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/v2/protocol"
	"github.com/StarsiegePlayers/darkstar-query-go/v2/query"
)

type Query struct {
	*protocol.Options
	Addresses []string
}

func NewQuery(timeout time.Duration, debug bool) *Query {
	return &Query{
		Options: &protocol.Options{
			Timeout:              timeout,
			Debug:                debug,
			MaxServerPacketSize:  protocol.MaxDataSize,
			MaxNetworkPacketSize: protocol.MaxPacketSize,
		},
		Addresses: []string{},
	}
}

type Results struct {
	Masters []*query.MasterQuery
	Games   []*query.PingInfoQuery
	Errors  []error
}

type ServerResult struct {
	Master *query.MasterQuery
	Game   *query.PingInfoQuery
	Error  error
}
