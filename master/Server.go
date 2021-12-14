package master

import (
	"net"
	"time"
)

type Server struct {
	Address  *net.UDPAddr
	LastSeen time.Time
	TTL      int
}

func (s Server) String() string {
	return s.Address.String()
}
