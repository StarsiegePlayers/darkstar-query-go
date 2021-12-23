package protocol

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Protocol_PacketTestSuite struct {
	suite.Suite
	Packet *Packet
}

func (t *Protocol_PacketTestSuite) SetupTest() {
	t.Packet = NewPacket()
}

func (t *Protocol_PacketTestSuite) TestProtocol_NewPacket() {
	question := &Packet{
		Version: Version,
	}

	t.Assert().Equal(question, t.Packet)
}

func (t *Protocol_PacketTestSuite) TestProtocol_UnmarshalBinary() {
	question := []byte{
		0x10, 0x06, 0x04, 0x01, 0x45, 0x00, 0x00, 0x02, 0xAA, 0xAA,
	}

	err := t.Packet.UnmarshalBinary(question)
	t.Assert().Nil(err)

	t.Assert().Equal(question[0], t.Packet.Version)
	t.Assert().Equal(question[1], byte(t.Packet.Type))
	t.Assert().Equal(question[2], t.Packet.Number)
	t.Assert().Equal(question[3], t.Packet.Total)
	t.Assert().Equal(binary.BigEndian.Uint16(question[4:4+2]), t.Packet.Key)
	t.Assert().Equal(binary.BigEndian.Uint16(question[6:6+2]), t.Packet.ID)
	t.Assert().Equal(question[8:], t.Packet.Data)
}

func (t *Protocol_PacketTestSuite) TestProtocol_UnmarshalBinary_TrimNull() {
	question := []byte{
		0x10, 0x06, 0x04, 0x01, 0x45, 0x00, 0x00, 0x02, 0x00, 0x00,
	}

	err := t.Packet.UnmarshalBinary(question)
	t.Assert().Nil(err)

	t.Assert().Equal(question[0], t.Packet.Version)
	t.Assert().Equal(question[1], byte(t.Packet.Type))
	t.Assert().Equal(question[2], t.Packet.Number)
	t.Assert().Equal(question[3], t.Packet.Total)
	t.Assert().Equal(binary.BigEndian.Uint16(question[4:4+2]), t.Packet.Key)
	t.Assert().Equal(binary.BigEndian.Uint16(question[6:6+2]), t.Packet.ID)
	t.Assert().Equal([]byte{}, t.Packet.Data)
}

func (t *Protocol_PacketTestSuite) TestProtocol_UnmarshalBinary_SmallPacket() {
	question := []byte{
		0x10, 0x06, 0x04, 0x01, 0x45, 0x00,
	}

	err := t.Packet.UnmarshalBinary(question)
	t.Assert().Nil(err)

	t.Assert().Equal(question[0], t.Packet.Version)
	t.Assert().Equal(question[1], byte(t.Packet.Type))
	t.Assert().Equal(question[2], t.Packet.Number)
	t.Assert().Equal(question[3], t.Packet.Total)
	t.Assert().Equal(binary.BigEndian.Uint16(question[4:4+2]), t.Packet.Key)
	t.Assert().Equal(uint16(0), t.Packet.ID)
	t.Assert().Equal([]byte{}, t.Packet.Data)
}

func (t *Protocol_PacketTestSuite) TestProtocol_MarshalBinary() {
	t.Packet.Type = MasterServerHeartbeat
	t.Packet.Number = 0xff
	t.Packet.Total = 0x00
	t.Packet.Key = 0x02
	t.Packet.ID = 0x00
	t.Packet.Data = []byte{}

	question := []byte{
		0x10, 0x05, 0xff, 0x00, 0x00, 0x02, 0x00, 0x00,
	}

	output, err := t.Packet.MarshalBinary()
	t.Assert().Nil(err)

	t.Assert().Equal(question, output)
}

func TestProtocol_PacketTestSuite(t *testing.T) {
	suite.Run(t, new(Protocol_PacketTestSuite))
}
