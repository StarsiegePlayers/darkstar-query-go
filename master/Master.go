package master

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
)

type Query struct {
	MasterData *Master
	Error      error
}

type Master struct {
	Address         string
	CommonName      string
	MOTD            string
	ServerCount     int
	ServerAddresses []string `json:"-"`
	Ping            time.Duration

	id           int
	conn         net.Conn
	requestStart int64
	requestEnd   int64
}

func (m *Master) UnmarshalBinary(p *protocol.Packet) error {
	data := p.Data
	commonNameLength, data := data[0], data[1:]
	commonName := string(data[0:commonNameLength])
	commonName = strings.ReplaceAll(commonName, `\n`, "")
	m.CommonName, data = commonName, data[commonNameLength:]

	motdLength, data := data[0], data[1:]
	motd, data := string(data[0:motdLength]), data[motdLength:]
	motd = strings.ReplaceAll(motd, `\n`, " ")
	if len(motd) > 0 {
		m.MOTD = motd[10:] // first 10 characters are junk
	}

	data = data[1:] // skip null terminator

	serverCount, data := data[0], data[1:]
	for i := byte(0); i < serverCount; i++ {
		data = data[1:] // skip separator byte "0x6"

		address := fmt.Sprintf("%d.%d.%d.%d", data[0], data[1], data[2], data[3])
		port := fmt.Sprintf("%d", binary.LittleEndian.Uint16(data[4:6]))
		data = data[6:]

		if address == "127.0.0.1" { // skip all servers reporting as localhost
			continue
		}
		m.ServerAddresses = append(m.ServerAddresses, fmt.Sprintf("%s:%s", address, port))
	}

	// log.Printf("Servercount: %d, datalen %d, countlen %d\n", serverCount, len(data), len(m.ServerAddresses))

	return nil
}

func queryPacket(id int) *protocol.Packet {
	packet := protocol.NewPacket()
	packet.Type = protocol.PingInfoQuery
	packet.Number = protocol.RequestAllPackets
	packet.Key = uint16(id)

	return packet
}

func (m *Master) parseResponse(conn net.Conn) error {
	// acquire data
	_ = conn.SetDeadline(time.Now().Add(1 * time.Second))
	for {
		data := make([]byte, protocol.MaxPacketSize)
		length, err := conn.Read(data)
		var netError *net.OpError
		if err != nil && errors.As(err, &netError) && netError.Timeout() {
			return fmt.Errorf("connection timed out")
		} else if err != nil {
			return fmt.Errorf("connection read failed: %w", err)
		}

		packet := protocol.NewPacket()
		err = packet.UnmarshalBinary(data)
		if err != nil {
			return fmt.Errorf("unmarshaling packet failed: %w", err)
		}

		err = m.UnmarshalBinary(packet)
		if err != nil {
			return fmt.Errorf("unmarshaling master data failed: %w", err)
		}

		if length <= protocol.MaxPacketSize || packet.Total <= 1 {
			break
		}
	}

	return nil
}

func (m *Master) Query(conn net.Conn, id int) error {
	m.conn = conn
	m.id = id

	query := queryPacket(id)

	// log.Printf("Server: %s - %s\n", conn.RemoteAddr(), query)

	data, err := query.MarshalBinary()
	if err != nil {
		return fmt.Errorf("master: [%s]: MarshalBinary failed: %w", m.Address, err)
	}

	m.requestStart = time.Now().UnixNano()
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("master: [%s]: connection Write failed: %w", m.Address, err)
	}

	err = m.parseResponse(conn)
	if err != nil {
		return fmt.Errorf("master: [%s]: %w", m.Address, err)
	}

	m.requestEnd = time.Now().UnixNano()

	m.ServerCount = len(m.ServerAddresses)
	m.Ping = time.Duration(m.requestEnd - m.requestStart)

	return nil
}
