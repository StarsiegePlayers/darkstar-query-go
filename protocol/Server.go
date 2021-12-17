package protocol

import (
	"net"
	"time"
)

type Server struct {
	Address    *net.UDPAddr
	Connection *net.PacketConn
	LastSeen   time.Time
	TTL        int
	Info       *Packet
}

func NewServerFromString(input string, ttl int) (*Server, error) {
	address, err := net.ResolveUDPAddr("udp", input)
	if err != nil {
		return nil, err
	}

	return &Server{
		Address:    address,
		Connection: nil,
		LastSeen:   time.Now(),
		TTL:        ttl,
	}, nil
}

func NewServersMapFromList(input []string) map[string]*Server {
	output := make(map[string]*Server)
	for _, v := range input {
		thisServer, _ := NewServerFromString(v, 300)
		output[v] = thisServer
	}
	return output
}

func (s Server) IsExpired() bool {
	return time.Since(s.LastSeen) >= time.Duration(s.TTL)*time.Second
}

func (s Server) String() string {
	return s.Address.String()
}
