package query

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/v2/protocol"
)

type PingInfoQuery struct {
	txID uint16

	*protocol.PingInfo

	Ping         time.Duration
	Address      string
	conn         net.Conn
	requestStart time.Time
	requestEnd   time.Time
	options      *protocol.Options
}

func NewPingInfoQuery(address string) *PingInfoQuery {
	options := &protocol.Options{
		Timeout: 2 * time.Second,
		Debug:   false,
	}

	return NewPingInfoQueryWithOptions(address, options)
}

func NewPingInfoQueryWithOptions(address string, options *protocol.Options) *PingInfoQuery {
	output := &PingInfoQuery{
		Address:  address,
		options:  options,
		txID:     0,
		PingInfo: &protocol.PingInfo{},
	}

	return output
}

func newPingInfoPacket() *protocol.Packet {
	packet := protocol.NewPacket()
	packet.Type = protocol.PingInfoQuery
	packet.Number = protocol.RequestAllPackets

	return packet
}

func (s *PingInfoQuery) parseResponse(conn io.Reader) error {
	data := make([]byte, protocol.MaxPacketSize)

	_, err := conn.Read(data)
	if err != nil {
		var netError *net.OpError
		if errors.As(err, &netError) && netError.Timeout() {
			return fmt.Errorf("connection timed out")
		}

		if s.options.Debug {
			return fmt.Errorf("connection read failed: %w", err)
		}

		return fmt.Errorf("connection read failed")
	}

	err = s.UnmarshalBinary(data)
	if err != nil {
		if s.options.Debug {
			return fmt.Errorf("unmarshaling pinginfo data failed: %w", err)
		}

		return fmt.Errorf("unspecified error parsing ping response")
	}

	return nil
}

func (s *PingInfoQuery) Query() error {
	var err error

	s.conn, err = net.Dial("udp", s.Address)
	if err != nil {
		return err
	}

	defer s.conn.Close()

	err = s.conn.SetDeadline(time.Now().Add(s.options.Timeout))
	if err != nil {
		return err
	}

	s.PingInfo.Packet = newPingInfoPacket()
	s.PingInfo.Packet.Key = s.txID

	data, err := s.PingInfo.MarshalBinary()
	if err != nil {
		return fmt.Errorf("pingInfo: [%s]: MarshalBinary failed: %w", s.conn.RemoteAddr(), err)
	}

	s.requestStart = time.Now()

	_, err = s.conn.Write(data)
	if err != nil {
		return fmt.Errorf("pingInfo: [%s]: connection Write failed: %w", s.conn.RemoteAddr(), err)
	}

	_ = s.conn.SetDeadline(time.Now().Add(s.options.Timeout))

	err = s.parseResponse(s.conn)
	if err != nil {
		return fmt.Errorf("pingInfo: [%s]: %w", s.conn.RemoteAddr(), err)
	}

	s.requestEnd = time.Now()
	s.Ping = s.requestEnd.Sub(s.requestStart)

	return nil
}
