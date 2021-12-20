package server

import (
	"net"
	"time"
)

type Server struct {
	Address    *net.UDPAddr
	Connection *net.PacketConn `json:"-" csv:"-"`
	LastSeen   time.Time
}

func NewServerFromString(input string) (*Server, error) {
	address, err := net.ResolveUDPAddr("udp", input)
	if err != nil {
		return nil, err
	}

	return &Server{
		Address:    address,
		Connection: nil,
		LastSeen:   time.Now(),
	}, nil
}

func NewServersMapFromList(input []string) map[string]*Server {
	output := make(map[string]*Server)
	for _, v := range input {
		thisServer, _ := NewServerFromString(v)
		output[v] = thisServer
	}
	return output
}

func (s Server) IsExpired(ttl time.Duration) bool {
	return time.Since(s.LastSeen) >= ttl
}

func (s Server) String() string {
	return s.Address.String()
}
