package protocol

import (
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/StarsiegePlayers/darkstar-query-go/v2/server"
)

type Master struct {
	*Packet    `json:"-" csv:"-"`
	Address    string
	CommonName string
	MOTDJunk   string `json:"-" csv:"-"`
	MOTD       string
	Servers    map[string]*server.Server
	MasterID   uint16
}

func NewMasterWithAddress(address string) (output *Master) {
	output = NewMaster()
	output.Address = address
	return
}

func NewMaster() (output *Master) {
	output = new(Master)
	output.Servers = make(map[string]*server.Server)
	output.MOTDJunk = "0000000000" // anything except all <0x00> will show the MOTD
	output.Packet = NewPacket()
	return
}

func (m *Master) UnmarshalBinarySet(data [][]byte) (err error) {
	for _, v := range data {
		err = m.UnmarshalBinary(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Master) UnmarshalBinary(data []byte) (err error) {
	err = m.Packet.UnmarshalBinary(data)
	if err != nil {
		return
	}

	p := m.Packet
	m.MasterID = p.ID
	if len(p.Data) <= 2 {
		return
	}

	// if it's the first packet (and only the first packet)
	// parse out the common name and MOTD
	if p.Number == 0x01 {
		m.CommonName, p.Data = ReadPascalStringStream(p.Data)
		m.CommonName = strings.ReplaceAll(m.CommonName, `\n`, "")

		m.MOTD, p.Data = ReadPascalStringStream(p.Data)
		m.MOTD = strings.ReplaceAll(m.MOTD, `\n`, " ")
		if len(m.MOTD) > 10 {
			// the first 10 characters are classified as "junk"
			m.MOTDJunk, m.MOTD = m.MOTD[0:10], m.MOTD[10:]
		}
	}

	if len(p.Data) <= 0 {
		return
	}
	p.Data = p.Data[1:] // null header separator

	serverCount := byte(0)
	serverCount, p.Data = p.Data[0], p.Data[1:]
	if serverCount <= 0 || len(p.Data) <= 0 {
		return
	}

	for i := byte(0); i < serverCount; i++ {
		p.Data = p.Data[1:] // skip separator byte "0x6"

		address := fmt.Sprintf("%d.%d.%d.%d", p.Data[0], p.Data[1], p.Data[2], p.Data[3])
		port := fmt.Sprintf("%d", binary.LittleEndian.Uint16(p.Data[4:4+2]))
		addressPort := fmt.Sprintf("%s:%s", address, port)
		p.Data = p.Data[6:]

		if address == "127.0.0.1" { // skip all servers reporting as localhost
			continue
		}

		svr, err := server.NewServerFromString(addressPort)
		if err != nil {
			continue
		}

		m.Servers[addressPort] = svr
	}

	// log.Printf("Servercount: %d, datalen %d, countlen %d\n", serverCount, len(data), len(m.ServerAddresses))

	return
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

func (m *Master) MarshalBinarySet(input map[string]*server.Server) []byte {
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

func (m *Master) GeneratePackets(options *Options, key uint16) [][]byte {
	serverAddresses := make([]string, 0)
	for k := range m.Servers {
		serverAddresses = append(serverAddresses, k)
	}
	sort.Strings(serverAddresses)

	output := make([][]byte, 0)

	// generate header
	header := m.MarshalBinaryHeader()
	firstPacketOverhead := uint16(len(header) + HeaderSize + 2) // uint16 little endian trailer for payload length

	// calculate packet sizes relative to payload data
	// 7 bytes per entry: <0x06 pascal-style entry length in bytes><4 bytes ipv4 address><2 bytes udpPort>
	firstPacketMax := (options.MaxServerPacketSize - firstPacketOverhead) / 7
	overflowPacketMax := (options.MaxServerPacketSize - (HeaderSize + 2)) / 7

	// calculate overflow from first packet
	overflowSize := len(serverAddresses) - int(firstPacketMax)
	overflowPackets := 0
	if overflowSize > 0 {
		overflowPackets = int(math.Ceil(float64(overflowSize)/float64(overflowPacketMax))) + 1
	}

	localAddresses := make([]string, len(serverAddresses))
	copy(localAddresses, serverAddresses)

	// send first packet
	pkt := NewPacket()
	pkt.Type = MasterServerList
	pkt.ID = m.MasterID
	pkt.Key = key

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
		binOut, err := pkt.MarshalBinary()
		if err != nil {
			// todo: log
			return [][]byte{}
		}

		output = append(output, binOut)

		// exit early
		return output
	}

	// otherwise, time to do some convoluted craziness
	pkt.Number = 0x01                 // start at 0x1
	pkt.Total = byte(overflowPackets) // overflow packets should be > 2

	// deep copy the first subset of addresses
	tmpAddresses := make(map[string]*server.Server)
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

	// copy to output
	pkt.Data = tempData
	binOut, err := pkt.MarshalBinary()
	if err != nil {
		// todo: log
		return [][]byte{}
	}
	output = append(output, binOut)

	// do the above for each overflow packet
	for i := 1; i < overflowPackets; i++ {
		pkt.Number++ // increment packet number

		// make sure we don't exceed the array for the last packet
		if uint16(len(localAddresses)) <= overflowPacketMax {
			overflowPacketMax = uint16(len(localAddresses))
		}

		// copy the next subset of overflow addresses
		tmpAddresses = make(map[string]*server.Server)
		for k, v := range localAddresses {
			if uint16(k) >= overflowPacketMax {
				break
			}
			tmpAddresses[v] = m.Servers[v]
		}
		localAddresses = localAddresses[overflowPacketMax:]

		// marshal and send
		pkt.Data = m.MarshalBinarySet(tmpAddresses)
		binOut, err = pkt.MarshalBinary()
		if err != nil {
			// todo: log
			return [][]byte{}
		}
		output = append(output, binOut)
	}

	return output
}
