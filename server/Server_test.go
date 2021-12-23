package server

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	Server  *Server
	Timeout time.Duration
}

func (t *ServerTestSuite) SetupTest() {
	t.Server = nil
	t.Timeout = 500 * time.Second
}

func (t *ServerTestSuite) TestNewServerFromString() {
	server, err := NewServerFromString("127.0.0.1:29001")
	t.Assert().Nil(err)

	address, err := net.ResolveUDPAddr("udp", "127.0.0.1:29001")
	t.Assert().Nil(err)
	t.Server = &Server{
		Address:    address,
		Connection: nil,
		LastSeen:   server.LastSeen,
	}

	t.Assert().Equal(t.Server, server)
}

func (t *ServerTestSuite) TestNewServerFromString_NoPort() {
	_, err := NewServerFromString("127.0.0.1")
	t.Assert().Contains(err.Error(), "missing port in address")
}

func (t *ServerTestSuite) TestNewServerFromString_Invalid() {
	_, err := NewServerFromString("256.256.256.256:29001")
	t.Assert().Contains(err.Error(), "no such host")
}

func (t *ServerTestSuite) TestIsExpired() {
	var err error
	t.Server, err = NewServerFromString("127.0.0.1:29001")
	t.Assert().Nil(err)

	expired := t.Server.IsExpired(t.Timeout)
	t.Assert().False(expired)

	expired = t.Server.IsExpired(time.Now().Sub(time.Now()))
	t.Assert().True(expired)
}

func (t *ServerTestSuite) TestString() {
	server, err := NewServerFromString("127.0.0.1:29001")
	t.Assert().Nil(err)

	t.Assert().Equal("127.0.0.1:29001", server.String())
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
