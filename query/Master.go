package query

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/v2/protocol"
)

type MasterQuery struct {
	*protocol.Master

	Ping time.Duration

	conn         net.Conn
	requestStart time.Time
	requestEnd   time.Time
	options      *protocol.Options
	txId         uint16
}

func newPacket() *protocol.Packet {
	packet := protocol.NewPacket()
	packet.Type = protocol.PingInfoQuery
	packet.Number = protocol.RequestAllPackets

	return packet
}

func NewMasterQuery(address string) *MasterQuery {
	options := &protocol.Options{
		Timeout: 2 * time.Second,
		Debug:   false,
	}
	return NewMasterQueryWithOptions(address, options)
}

func NewMasterQueryWithOptions(address string, options *protocol.Options) (output *MasterQuery) {
	output = &MasterQuery{
		Master:  protocol.NewMasterWithAddress(address),
		options: options,
		txId:    0,
	}
	return
}

func (m *MasterQuery) parseResponse() error {
	// acquire data
	m.Total = 1
	for i := byte(0); i < m.Total; i++ {
		data := make([]byte, protocol.MaxPacketSize)
		length, err := m.conn.Read(data)
		if err != nil {
			var netError *net.OpError
			if errors.As(err, &netError) && netError.Timeout() {
				return fmt.Errorf("connection timed out")
			}
			if m.options.Debug {
				return fmt.Errorf("connection read failed: %w", err)
			}
			return fmt.Errorf("connection read failed")
		}

		m.Packet = protocol.NewPacket()
		err = m.Packet.UnmarshalBinary(data)
		if err != nil {
			if m.options.Debug {
				return fmt.Errorf("unmarshaling packet failed: %w", err)
			}
			return fmt.Errorf("unspecified error parsing packet")
		}

		m.Master = protocol.NewMasterWithAddress(m.Address)
		err = m.Master.UnmarshalBinary(data)
		if err != nil {
			if m.options.Debug {
				return fmt.Errorf("unmarshaling master data failed: %w", err)
			}
			return fmt.Errorf("unspecified error parsing master response")
		}

		if length <= protocol.MaxPacketSize || m.Packet.Total <= 1 {
			break
		}
	}

	return nil
}

func (m *MasterQuery) Query() error {
	var err error
	m.conn, err = net.Dial("udp", m.Address)
	if err != nil {
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) {
			if m.options.Debug {
				return fmt.Errorf("master: [%s]: dns error during dial [%s]", dnsError.Name, err)
			}
			return fmt.Errorf("master: [%s]: no such host", dnsError.Name)
		}
		if m.options.Debug {
			return fmt.Errorf("master: [%s]: error during dial [%s]", dnsError.Name, err)
		}
		return fmt.Errorf("master: [%s]: unspecified error during network connection", m.Address)
	}
	defer m.conn.Close()

	err = m.conn.SetDeadline(time.Now().Add(m.options.Timeout))
	if err != nil {
		if m.options.Debug {
			return fmt.Errorf("master: [%s]: error during m.conn.SetDeadline [%s]", m.Address, err)
		}
		return fmt.Errorf("master: [%s]: unspecified error while setting connection timeout", m.Address)
	}

	query := newPacket()
	query.Key = m.txId
	m.txId++

	// log.Printf("Server: %s - %s\n", conn.RemoteAddr(), query)

	data, err := query.MarshalBinary()
	if err != nil {
		if m.options.Debug {
			return fmt.Errorf("master: [%s]: MarshalBinary failed: %w", m.Address, err)
		}
		return fmt.Errorf("master: [%s]: Error parsing response", m.Address)
	}

	m.requestStart = time.Now()
	_, err = m.conn.Write(data)
	if err != nil {
		if m.options.Debug {
			return fmt.Errorf("master: [%s]: m.Conn.Write failed: %w", m.Address, err)
		}
		return fmt.Errorf("master: [%s]: connection refused", m.Address)
	}

	_ = m.conn.SetDeadline(time.Now().Add(m.options.Timeout))
	err = m.parseResponse()
	if err != nil {
		return fmt.Errorf("master: [%s]: %w", m.Address, err)
	}

	m.requestEnd = time.Now()
	m.Ping = m.requestEnd.Sub(m.requestStart)

	return nil
}
