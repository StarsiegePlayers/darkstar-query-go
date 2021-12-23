package protocol

import (
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PingInfoTestSuite struct {
	suite.Suite
	PingInfo *PingInfo
	ConnPipe net.Conn
}

func (t *PingInfoTestSuite) SetupTest() {
	t.PingInfo = new(PingInfo)
	t.PingInfo.Packet = NewPacket()
}

func (t *PingInfoTestSuite) TestServer_PingInfoUnmarshal() {
	question := []byte{
		0x10, 0x04, 0xFF, 0xFD, 0x00, 0x00, 0x40, 0x00, 0x65, 0x73, 0x33, 0x61, 0x06, 0x56, 0x20, 0x30,
		0x30, 0x31, 0x2E, 0x30, 0x30, 0x30, 0x72, 0x44, 0x4F, 0x56, 0x3A, 0x20, 0x43, 0x69, 0x74, 0x79,
		0x20, 0x4F, 0x6E, 0x20, 0x54, 0x68, 0x65,
	}
	response := &PingInfo{
		Packet: &Packet{
			Version: 0x10,
			Type:    0x04,
			Number:  0xff,
			Total:   0xfd,
			Key:     0x00,
			ID:      0x4000,
			Data:    question[23:],
		},
		GameMode:    0xfd,
		GameName:    "es3a",
		GameVersion: "V 001.000r",
		GameStatus:  0x6,
		PlayerCount: 0x0,
		MaxPlayers:  0x40,
		Name:        "DOV: City On The",
	}

	err := t.PingInfo.UnmarshalBinary(question)
	t.Assert().Nil(err)

	t.Assert().Equal(response, t.PingInfo)
}

func TestPingInfoTestSuite(t *testing.T) {
	suite.Run(t, new(PingInfoTestSuite))
}
