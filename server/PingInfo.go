package server

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
)

type Query struct {
	ServerInfo *PingInfo
	Error      error
}

type PingInfo struct {
	GameMode    byte   // ??
	GameName    string // es3a
	GameVersion string // V 001.000r
	GameStatus  StatusByte
	PlayerCount byte
	MaxPlayers  byte
	Name        string
	Address     string
	Ping        time.Duration

	id           int
	conn         net.Conn
	requestStart int64
	requestEnd   int64
}

func (s *PingInfo) String() string {
	return fmt.Sprintf("PingInfoResponse: %s [%s] (%s) Players: %d/%d @ %s",
		s.Name, s.Ping, s.GameStatus, s.PlayerCount, s.MaxPlayers, s.Address)
}

func pingInfoPacket(id int) *protocol.Packet {
	packet := protocol.NewPacket()
	packet.Type = protocol.PingInfoQuery
	packet.Number = protocol.RequestAllPackets
	packet.Key = uint16(id)

	return packet
}

func (s *PingInfo) UnmarshalBinary(p *protocol.Packet) error {
	s.GameMode = p.Total
	s.PlayerCount = byte((p.ID >> 8)& 0xff)
	s.MaxPlayers = byte(p.ID & 0xff)

	s.GameName, p.Data = string(p.Data[0:4]), p.Data[4:]
	s.GameStatus, p.Data = StatusByte(p.Data[0]), p.Data[1:]
	s.GameVersion, p.Data = string(p.Data[0:10]), p.Data[10:]

	p.Data[len(p.Data)-1] = 0x00
	s.Name = string(p.Data[:protocol.Clen(p.Data)])

	s.Address = s.conn.RemoteAddr().String()

	return nil
}

func (s *PingInfo) parseResponse(conn net.Conn, options protocol.Options) error {
	// acquire data
	for {
		data := make([]byte, protocol.MaxPacketSize)
		length, err := conn.Read(data)
		if err != nil {
			var netError *net.OpError
			if errors.As(err, &netError) && netError.Timeout() {
				return fmt.Errorf("connection timed out")
			}
			if options.Debug {
				return fmt.Errorf("connection read failed: %w", err)
			}
			return fmt.Errorf("connection read failed")
		}

		packet := protocol.NewPacket()
		err = packet.UnmarshalBinary(data)
		if err != nil {
			if options.Debug {
				return fmt.Errorf("unmarshaling packet failed: %w", err)
			}
			return fmt.Errorf("unspecified error parsing packet")
		}

		err = s.UnmarshalBinary(packet)
		if err != nil {
			if options.Debug {
				return fmt.Errorf("unmarshaling pinginfo data failed: %w", err)
			}
			return fmt.Errorf("unspecified error parsing ping response")
		}

		if length <= protocol.MaxPacketSize || packet.Total <= 1 {
			break
		}
	}

	return nil
}

func (s *PingInfo) PingInfoQuery(conn net.Conn, id int, options protocol.Options) error {
	s.conn = conn
	s.id = id

	query := pingInfoPacket(id)

	data, err := query.MarshalBinary()
	if err != nil {
		return fmt.Errorf("game: [%s]: MarshalBinary failed: %w", conn.RemoteAddr(), err)
	}

	s.requestStart = time.Now().UnixNano()
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("game: [%s]: connection Write failed: %w", conn.RemoteAddr(), err)
	}

	_ = conn.SetDeadline(time.Now().Add(options.Timeout))
	err = s.parseResponse(conn, options)
	if err != nil {
		return fmt.Errorf("game: [%s]: %w", conn.RemoteAddr(), err)
	}

	s.requestEnd = time.Now().UnixNano()
	s.Ping = time.Duration(s.requestEnd - s.requestStart)

	return nil
}
