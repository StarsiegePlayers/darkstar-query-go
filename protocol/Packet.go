package protocol

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
)

const Version = 0x10
const VersionExt = 0x69
const MaxPacketSize = 1500 // standard packet MTU
const HeaderSize = 8
const MaxDataSize = MaxPacketSize - HeaderSize

var ErrorUnknownPacketVersion = errors.New("unknown packet version")
var ErrorEmptyPacket = errors.New("empty packet received")

type Packet struct {
	Version byte
	Number  byte   // packet number out of total; starts at 1
	Total   byte   // total packets of info
	Key     uint16 // used for verification and (transaction) id purposes
	ID      uint16 // master server id (read from config file)
	Data    []byte // MaxSize = (MaxPacketSize - HeaderSize)
	Type    PacketType

	// implements
	encoding.BinaryMarshaler   `json:"-" csv:"-"`
	encoding.BinaryUnmarshaler `json:"-" csv:"-"`
	fmt.Stringer               `json:"-" csv:"-"`
}

func NewPacket() *Packet {
	return &Packet{
		Version: Version,
	}
}

func NewPacketWithData(data []byte) (out *Packet, err error) {
	out = NewPacket()
	err = out.UnmarshalBinary(data)

	return
}

func (p *Packet) MarshalBinary() ([]byte, error) {
	out := make([]byte, HeaderSize+len(p.Data))
	out[0] = p.Version
	out[1] = byte(p.Type)
	out[2] = p.Number
	out[3] = p.Total
	binary.BigEndian.PutUint16(out[4:4+2], p.Key) // BigE - key byte numbers 5-6
	binary.BigEndian.PutUint16(out[6:6+2], p.ID)  // BigE - ID byte numbers 7-8
	copy(out[HeaderSize:], p.Data)

	return out, nil
}

func (p *Packet) UnmarshalBinary(data []byte) error {
	data = bytes.Trim(data, "\x00")

	if len(data) < 0 {
		return ErrorEmptyPacket
	}

	// pad data to a minimum of 8 bytes
	if len(data) < 8 {
		tmp := make([]byte, 8)
		copy(tmp, data)
		data = tmp
	}

	p.Data = make([]byte, len(data)-HeaderSize)

	if data[0] != Version && data[0] != VersionExt {
		return ErrorUnknownPacketVersion
	}

	p.Version = data[0]
	p.Type = PacketType(data[1])
	p.Number = data[2]
	p.Total = data[3]
	p.Key = binary.BigEndian.Uint16(data[4 : 4+2]) // BigE - key index 5-6
	p.ID = binary.BigEndian.Uint16(data[6 : 6+2])  // BigE - ID index 7-8
	copy(p.Data, data[8:])                         // payload

	return nil
}

func (p *Packet) String() string {
	packet, err := p.MarshalBinary()
	if err != nil {
		return "![unable to create packet]"
	}

	return fmt.Sprintf("%x", packet)
}
