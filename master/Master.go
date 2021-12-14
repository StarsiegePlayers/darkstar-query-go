package master

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
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
	ServerAddresses []string `json:"-"`
	Ping            time.Duration
	MasterID        uint16

	id           int
	conn         net.Conn
	requestStart int64
	requestEnd   int64
}

func (m *Master) UnmarshalBinary(p *protocol.Packet) error {
	m.MasterID = p.ID
	data := p.Data
	if len(data) <= 2 {
		return nil
	}

	// byte 3 after the packet header will be 0x6 on a spanned master response
	// if it's not, parse the canonical name / motd header
	if p.Data[2] != 0x06 {
		m.CommonName, data = protocol.ReadPascalStringStream(data)
		m.CommonName = strings.ReplaceAll(m.CommonName, `\n`, "")

		m.MOTD, data = protocol.ReadPascalStringStream(data)
		m.MOTD = strings.ReplaceAll(m.MOTD, `\n`, " ")
		if len(m.MOTD) > 10 {
			m.MOTD = m.MOTD[10:] // first 10 characters are junk
		}
	}

	if len(data) <= 0 {
		return nil
	}
	data = data[1:] // null header separator

	serverCount, data := data[0], data[1:]
	if serverCount <= 0 || len(data) <= 0 {
		return nil
	}

	for i := byte(0); i < serverCount; i++ {
		data = data[1:] // skip separator byte "0x6"

		address := fmt.Sprintf("%d.%d.%d.%d", data[0], data[1], data[2], data[3])
		port := fmt.Sprintf("%d", binary.LittleEndian.Uint16(data[4:4+2]))
		data = data[6:]

		if address == "127.0.0.1" { // skip all servers reporting as localhost
			continue
		}
		m.ServerAddresses = append(m.ServerAddresses, fmt.Sprintf("%s:%s", address, port))
	}

	// log.Printf("Servercount: %d, datalen %d, countlen %d\n", serverCount, len(data), len(m.ServerAddresses))

	return nil
}

func (m *Master) MarshalBinaryHeader() []byte {
	// field 01 - pascal common name, string
	commonName := make([]byte, len(m.CommonName)+1)
	commonName[0] = byte(len(m.CommonName))
	copy(commonName[1:], m.CommonName)

	// field 02 - pascal MOTD string, incl 10 character spacer
	motd := make([]byte, len(m.MOTD)+12)
	motd[0] = byte(len(motd) - 2)    // exclude size byte and trailer null
	copy(motd[1:1+10], "dummythicc") // magic 10 characters
	copy(motd[11:len(m.MOTD)+11], m.MOTD)
	copy(motd[len(motd)-1:], "\x00")

	// combine
	output := make([]byte, len(commonName)+len(motd))
	copy(output, commonName)
	copy(output[len(commonName):], motd)

	return output
}

func (m *Master) MarshalBinarySet(set []string) []byte {
	// work with a byte buffer
	hold := make([]byte, len(set)*7)
	for index, hostPort := range set {
		// split ip:port
		temp := strings.Split(hostPort, ":")
		port := temp[1]
		host := strings.Split(temp[0], ".")

		// individual entry byte buffer
		out := make([]byte, 7)

		// 6 bytes is the length of the ip+host as a BCD
		out[0] = byte(6)
		for k2, v2 := range host {
			h, err := strconv.Atoi(v2)
			if err != nil {
				continue
			}
			out[k2+1] = byte(h)
		}

		// ports are sent as a uint16 little endian stream
		p, _ := strconv.Atoi(port)
		binary.LittleEndian.PutUint16(out[5:], uint16(p))

		// current entry * length of entry -> current entry * length of entry + length of entry
		copy(hold[index*7:index*7+7], out)
	}

	output := make([]byte, len(hold)+1)
	output[0] = byte(len(set))
	copy(output[1:], hold)

	return output
}

func (m *Master) SendResponse(conn net.Conn, options protocol.Options) {
	// generate header
	header := m.MarshalBinaryHeader()
	overhead := uint16(len(header) + protocol.HeaderSize + 2) // uint16 little endian trailer for payload length

	// calculate packet sizes relative to payload data
	firstPacketMax := (options.MaxServerPacketSize - overhead) / 7
	overflowPacketMax := options.MaxServerPacketSize - uint16(len(header)+2)/7

	// calculate overflow from first packet
	overflowSize := len(m.ServerAddresses) - int(firstPacketMax)
	overflowPackets := 0
	if overflowSize > 0 {
		overflowPackets = int(math.Ceil(float64(overflowSize)/float64(overflowPacketMax))) + 1
	}

	localAddresses := make([]string, len(m.ServerAddresses))
	copy(localAddresses, m.ServerAddresses)

	// send first packet
	pkt := protocol.NewPacket()
	pkt.Type = protocol.MasterServerList
	pkt.ID = m.MasterID

	// simple logic for non spanned packets
	if overflowPackets <= 0 {
		// setting pkt 1 of 1 is distinctly different from ping/game info
		pkt.Number = 0x01
		pkt.Total = 0x01
		dataset := m.MarshalBinarySet(localAddresses)
		tempData := make([]byte, len(header)+len(dataset))
		copy(tempData[0:len(header)], header)
		copy(tempData[len(header):len(header)+len(dataset)], dataset)
		pkt.Data = tempData
		output, err := pkt.MarshalBinary()
		if err != nil {
			// todo: log
			return
		}

		_, err = conn.Write(output)
		if err != nil {
			// todo: log
			return
		}

		// exit early
		return
	}

	// otherwise, time to do some convoluted craziness
	pkt.Number = 0x01                 // start at 0x1
	pkt.Total = byte(overflowPackets) // overflow packets should be > 2

	// deep copy the first subset of addresses
	tmpAddresses := make([]string, firstPacketMax)
	copy(tmpAddresses, localAddresses[:firstPacketMax])

	// pop the elements we just copied
	localAddresses = localAddresses[firstPacketMax:]

	// marshal data
	dataset := m.MarshalBinarySet(tmpAddresses)
	tempData := make([]byte, len(header)+len(dataset))
	copy(tempData[0:len(header)], header)
	copy(tempData[len(header):len(header)+len(dataset)], dataset)

	// send
	pkt.Data = tempData
	output, err := pkt.MarshalBinary()
	if err != nil {
		// todo: log
		return
	}
	_, err = conn.Write(output)
	if err != nil {
		// todo: log
		return
	}

	// do the above for each overflow packet
	for i := 1; i < overflowPackets; i++ {
		pkt.Number++ // increment packet number

		// make sure we don't exceed the array for the last packet
		if uint16(len(localAddresses)) <= overflowPacketMax {
			overflowPacketMax = uint16(len(localAddresses) - 1)
		}

		// copy the next subset of overflow addresses
		tmpAddresses := make([]string, overflowPacketMax)
		copy(tmpAddresses, localAddresses[:overflowPacketMax])
		localAddresses = localAddresses[overflowPacketMax:]

		// marshal and send
		pkt.Data = m.MarshalBinarySet(tmpAddresses)
		output, err := pkt.MarshalBinary()
		if err != nil {
			// todo: log
			return
		}
		_, err = conn.Write(output)
		if err != nil {
			// todo: log
			return
		}
	}

}

func queryPacket(id int) *protocol.Packet {
	packet := protocol.NewPacket()
	packet.Type = protocol.PingInfoQuery
	packet.Number = protocol.RequestAllPackets
	packet.Key = uint16(id)

	return packet
}

func (m *Master) parseResponse(conn net.Conn, options protocol.Options) error {
	spannedSet := true
	// acquire data
	for spannedSet == true {
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
		if (packet.Number == 0xff || packet.Total == 0x00) || (packet.Number == packet.Total) {
			spannedSet = false
		}

		err = m.UnmarshalBinary(packet)
		if err != nil {
			if options.Debug {
				return fmt.Errorf("unmarshaling master data failed: %w", err)
			}
			return fmt.Errorf("unspecified error parsing master response")
		}

		if length <= protocol.MaxPacketSize || packet.Total <= 1 {
			break
		}
	}

	return nil
}

func (m *Master) Query(conn net.Conn, id int, options protocol.Options) error {
	m.conn = conn
	m.id = id

	query := queryPacket(id)

	// log.Printf("Server: %s - %s\n", conn.RemoteAddr(), query)

	data, err := query.MarshalBinary()
	if err != nil {
		if options.Debug {
			return fmt.Errorf("master: [%s]: MarshalBinary failed: %w", m.Address, err)
		}
		return fmt.Errorf("master: [%s]: Error parsing response", m.Address)
	}

	m.requestStart = time.Now().UnixNano()
	_, err = conn.Write(data)
	if err != nil {
		if options.Debug {
			return fmt.Errorf("master: [%s]: connection Write failed: %w", m.Address, err)
		}
		return fmt.Errorf("master: [%s]: connection refused", m.Address)
	}

	_ = conn.SetDeadline(time.Now().Add(options.Timeout))
	err = m.parseResponse(conn, options)
	if err != nil {
		return fmt.Errorf("master: [%s]: %w", m.Address, err)
	}

	m.requestEnd = time.Now().UnixNano()
	m.Ping = time.Duration(m.requestEnd - m.requestStart)

	return nil
}
