package darkstar_query_go

import (
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
)

func NewOptions() protocol.Options {
	return protocol.Options{
		Search:  nil,
		Timeout: 5 * time.Second,
		Debug:   false,
	}
}
