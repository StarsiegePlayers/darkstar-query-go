package master

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"sort"
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
	Address    string
	CommonName string
	MOTDJunk   string `json:"-"`
	MOTD       string
	Servers    map[string]*protocol.Server `json:"-"`
	Ping       time.Duration
	MasterID   uint16

	id           int
	conn         net.Conn
	requestStart int64
	requestEnd   int64
}

func NewMaster() *Master {
	output := new(Master)
	output.Servers = make(map[string]*protocol.Server)
	output.MOTDJunk = "0000000000" // anything except all <0x00> will show the MOTD
	return output
}

func (m *Master) UnmarshalBinary(p *protocol.Packet) error {
	m.MasterID = p.ID
	data := p.Data
	if len(data) <= 2 {
		return nil
	}

	// if it's the first packet (and only the first packet)
	// parse out the common name and MOTD
	if p.Number == 0x01 {
		m.CommonName, data = protocol.ReadPascalStringStream(data)
		m.CommonName = strings.ReplaceAll(m.CommonName, `\n`, "")

		m.MOTD, data = protocol.ReadPascalStringStream(data)
		m.MOTD = strings.ReplaceAll(m.MOTD, `\n`, " ")
		if len(m.MOTD) > 10 {
			// the first 10 characters are classified as "junk"
			m.MOTDJunk, m.MOTD = m.MOTD[0:10], m.MOTD[10:]
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
		addressPort := fmt.Sprintf("%s:%s", address, port)
		data = data[6:]

		if address == "127.0.0.1" { // skip all servers reporting as localhost
			continue
		}

		server, err := protocol.NewServerFromString(addressPort, 300)
		if err != nil {
			continue
		}

		m.Servers[addressPort] = server
	}

	// log.Printf("Servercount: %d, datalen %d, countlen %d\n", serverCount, len(data), len(m.ServerAddresses))

	return nil
}

func (m *Master) MarshalBinaryHeader() []byte {
	// field 01 - pascal common name, string
	commonName := make([]byte, len(m.CommonName)+1)
	commonName[0] = byte(len(m.CommonName))
	copy(commonName[1:], m.CommonName)

	motd := make([]byte, 2)
	if len(m.MOTD) > 0 {
		// field 02 - pascal MOTD string, incl 10 character spacer
		motd = make([]byte, len(m.MOTD)+12)
		motd[0] = byte(len(motd) - 2)  // exclude size byte and trailer null
		copy(motd[1:1+10], m.MOTDJunk) // magic 10 characters
		copy(motd[11:len(m.MOTD)+11], m.MOTD)
		copy(motd[len(motd)-1:], "\x00")
	}

	// combine
	output := make([]byte, len(commonName)+len(motd))
	copy(output, commonName)
	copy(output[len(commonName):], motd)

	return output
}

func (m *Master) MarshalBinarySet(input map[string]*protocol.Server) []byte {
	set := make([]string, 0)
	for k := range input {
		set = append(set, k)
	}
	sort.Strings(set)

	// work with a byte buffer
	hold := make([]byte, len(set)*7)

	for index, hostPort := range set {
		// split ip:port
		temp := strings.Split(hostPort, ":")
		port := temp[1]
		host := strings.Split(temp[0], ".")

		// individual entry byte buffer
		out := make([]byte, 7)

		// <length: 0x06><4 bytes ipv4 addr><2 bytes port>
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

		// current entry * length of entry : current entry * length of entry + length of entry
		copy(hold[index*7:index*7+7], out)
	}

	output := make([]byte, len(hold)+1)
	output[0] = byte(len(set))
	copy(output[1:], hold)

	return output
}

func (m *Master) SendResponse(conn *net.PacketConn, addr *net.UDPAddr, options *protocol.Options) {
	serverAddresses := make([]string, 0)
	for k := range m.Servers {
		serverAddresses = append(serverAddresses, k)
	}
	sort.Strings(serverAddresses)

	// generate header
	header := m.MarshalBinaryHeader()
	firstPacketOverhead := uint16(len(header) + protocol.HeaderSize + 2) // uint16 little endian trailer for payload length

	// calculate packet sizes relative to payload data
	// 7 bytes per entry: <0x06 pascal-style entry length in bytes><4 bytes ipv4 address><2 bytes udpPort>
	firstPacketMax := (options.MaxServerPacketSize - firstPacketOverhead) / 7
	overflowPacketMax := (options.MaxServerPacketSize - (protocol.HeaderSize + 2)) / 7

	// calculate overflow from first packet
	overflowSize := len(serverAddresses) - int(firstPacketMax)
	overflowPackets := 0
	if overflowSize > 0 {
		overflowPackets = int(math.Ceil(float64(overflowSize)/float64(overflowPacketMax))) + 1
	}

	localAddresses := make([]string, len(serverAddresses))
	copy(localAddresses, serverAddresses)

	// send first packet
	pkt := protocol.NewPacket()
	pkt.Type = protocol.MasterServerList
	pkt.ID = m.MasterID
	pkt.Key = options.PacketKey

	// simple logic for non spanned packets
	if overflowPackets <= 0 {
		// setting pkt 1 of 1 is distinctly different from ping/game info
		pkt.Number = 0x01
		pkt.Total = 0x01
		dataset := m.MarshalBinarySet(m.Servers)
		tempData := make([]byte, len(header)+len(dataset))
		copy(tempData[0:len(header)], header)
		copy(tempData[len(header):len(header)+len(dataset)], dataset)
		pkt.Data = tempData
		output, err := pkt.MarshalBinary()
		if err != nil {
			// todo: log
			return
		}

		_, err = (*conn).WriteTo(output, addr)
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
	tmpAddresses := make(map[string]*protocol.Server)
	for k, v := range localAddresses {
		if uint16(k) >= firstPacketMax {
			break
		}
		tmpAddresses[v] = m.Servers[v]
	}

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
	_, err = (*conn).WriteTo(output, addr)
	if err != nil {
		// todo: log
		return
	}

	// do the above for each overflow packet
	for i := 1; i < overflowPackets; i++ {
		pkt.Number++ // increment packet number

		// make sure we don't exceed the array for the last packet
		if uint16(len(localAddresses)) <= overflowPacketMax {
			overflowPacketMax = uint16(len(localAddresses))
		}

		// copy the next subset of overflow addresses
		tmpAddresses = make(map[string]*protocol.Server)
		for k, v := range localAddresses {
			if uint16(k) >= overflowPacketMax {
				break
			}
			tmpAddresses[v] = m.Servers[v]
		}
		localAddresses = localAddresses[overflowPacketMax:]

		// marshal and send
		pkt.Data = m.MarshalBinarySet(tmpAddresses)
		output, err := pkt.MarshalBinary()
		if err != nil {
			// todo: log
			return
		}
		_, err = (*conn).WriteTo(output, addr)
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

func (m *Master) parseResponse(conn net.Conn, options *protocol.Options) error {
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

func (m *Master) Query(conn net.Conn, id int, options *protocol.Options) error {
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
