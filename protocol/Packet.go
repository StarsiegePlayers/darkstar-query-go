package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const Version = 0x10
const MaxPacketSize = 1500 // standard packet MTU
const HeaderSize = 8
const MaxDataSize = MaxPacketSize - HeaderSize

var ErrorAllDataMarshaled = errors.New("all data has been marshaled")
var ErrorUnknownPacketVersion = errors.New("unknown packet version")

type Packet struct {
	Version byte
	Type    PacketType
	Number  byte   // packet number out of total; starts at 1
	Total   byte   // total packets of info
	Key     uint16 // for verification and id purposes
	ID      uint16 // master server id (read from config file)
	Data    []byte // MaxSize = (MaxPacketSize - HeaderSize)

	currentPacket uint16
}

func NewPacket() *Packet {
	return &Packet{
		Version:       Version,
		currentPacket: 0,
	}
}

func (p *Packet) marshalBinaryInt() ([]byte, error) {
	if len(p.Data) < int(p.currentPacket*MaxDataSize) {
		return nil, ErrorAllDataMarshaled
	}
	data := bytes.Trim(p.Data, "\x00")

	packetSize := MaxPacketSize
	if len(data) <= 0 {
		packetSize = HeaderSize
	} else if len(data) <= MaxPacketSize {
		packetSize = len(data)
	} else if len(data[p.currentPacket*MaxDataSize:]) > MaxDataSize {
		packetSize = MaxPacketSize
	}

	out := make([]byte, packetSize)
	out[0] = p.Version
	out[1] = byte(p.Type)
	out[2] = p.Number
	out[3] = byte(len(p.Data[p.currentPacket*MaxDataSize:]) % MaxDataSize)
	binary.LittleEndian.PutUint16(out[4:], p.Key) // key index 5-6
	binary.LittleEndian.PutUint16(out[6:], p.ID)  // ID index 7-8
	copy(out[8:], p.Data[p.currentPacket*MaxDataSize:])

	return out, nil
}

func (p *Packet) MarshalBinary() ([]byte, error) {
	out, err := p.marshalBinaryInt()
	if err != nil {
		return nil, err
	}

	p.currentPacket++

	return out, nil
}

func (p *Packet) UnmarshalBinary(data []byte) error {
	data = bytes.Trim(data, "\x00")
	if len(data) <= 0 {
		return fmt.Errorf("no data recieved")
	}
	p.Data = make([]byte, len(data)-HeaderSize)
	if data[0] != p.Version {
		return ErrorUnknownPacketVersion
	}
	p.Version = data[0]
	p.Type = PacketType(data[1])
	p.Number = data[2]
	p.Total = data[3]
	p.Key = binary.LittleEndian.Uint16(data[3:5])
	p.ID = binary.LittleEndian.Uint16(data[5:7])
	copy(p.Data, data[8:])

	return nil
}

func (p *Packet) String() string {
	packet, err := p.marshalBinaryInt()
	if err != nil {
		return "![unable to create packet]"
	}
	return fmt.Sprintf("%x", packet)
}
