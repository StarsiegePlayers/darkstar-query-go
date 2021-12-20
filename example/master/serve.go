package main

import (
	"fmt"
	darkstar "github.com/StarsiegePlayers/darkstar-query-go/v2"
	"github.com/StarsiegePlayers/darkstar-query-go/v2/protocol"
	"github.com/StarsiegePlayers/darkstar-query-go/v2/server"
	"net"
	"sync"
	"time"
)

type Service struct {
	sync.Mutex
	Service          *protocol.Master
	BannedService    *protocol.Master
	Options          *protocol.Options
	IPServiceCount   map[string]uint16
	SolicitedServers map[string]bool
}

var thisMaster = Service{
	Service:          protocol.NewMaster(),
	BannedService:    protocol.NewMaster(),
	Options:          &protocol.Options{},
	IPServiceCount:   make(map[string]uint16),
	SolicitedServers: make(map[string]bool),
}

func serve(conn *net.PacketConn, addr *net.UDPAddr, buf []byte) {
	// we use an ip-port combo as a unique identifier
	ipPort := fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)

	// parse packet
	p := protocol.NewPacket()
	err := p.UnmarshalBinary(buf)
	if err != nil {
		switch err {
		case protocol.ErrorUnknownPacketVersion:
			LogServerAlert(ipPort, "Unknown protocol number")
		case protocol.ErrorEmptyPacket:
			LogServerAlert(ipPort, "Empty packet received")
		default:
			LogServerAlert(ipPort, "Error %s while parsing packet", err)
		}
		return
	}

	isBanned := false
	for _, v := range config.parsedBannedNets {
		if v.Contains(addr.IP) {
			isBanned = true
		}
	}

	switch p.Type {
	// server has sent in a heartbeat
	case protocol.MasterServerHeartbeat:
		if isBanned {
			LogServerAlert(addr.IP.String(), "Received a %s packet from banned host", p.Type.String())
			return
		}
		registerHeartbeat(conn, addr, ipPort, p)
		break

	// client is requesting a server list
	case protocol.PingInfoQuery:
		if isBanned {
			LogServerAlert(addr.IP.String(), "Received a %s packet from banned host", p.Type.String())
			sendBanned(conn, addr, ipPort, p)
			return
		}
		sendList(conn, addr, ipPort, p)
		break

	default:
		if isBanned {
			LogServerAlert(ipPort, "Received unsolicited packet type %s from banned host", p.Type.String())
			return
		}
		LogServerAlert(ipPort, "Received unsolicited packet type %s", p.Type.String())
	}
}

func registerHeartbeat(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	thisMaster.Lock()
	thisMaster.SolicitedServers[ipPort] = true
	thisMaster.Unlock()

	q := darkstar.NewQuery(2*time.Second, true)
	q.Addresses = append(q.Addresses, ipPort)
	response, err := q.Servers()
	if len(err) > 0 || len(response) <= 0 {
		LogServerAlert(ipPort, "error during server verification [%s, %d]", err, len(response))
		return
	}

	registerPingInfo(conn, addr, ipPort, p)
}

func registerPingInfo(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	thisMaster.Lock()
	if _, exist := thisMaster.Service.Servers[ipPort]; !exist {
		count := thisMaster.IPServiceCount[addr.IP.String()]
		if uint16(count)+1 > config.ServersPerIP {
			LogServerAlert(ipPort, "Rejecting additional server for IP - count: %d/%d", count, config.ServersPerIP)
			thisMaster.Unlock()
			return
		}

		// log and add new
		LogServer(ipPort, "Heartbeat - New Server")
		thisMaster.Service.Servers[ipPort] = &server.Server{
			Address:    addr,
			Connection: conn,
			LastSeen:   time.Now(),
		}
		count++
		LogServer(ipPort, "New Server for IP - total server count for IP: %d/%d", count, config.ServersPerIP)
		thisMaster.IPServiceCount[addr.IP.String()] = count
	}

	s := thisMaster.Service.Servers[ipPort]
	LogServer(ipPort, "Heartbeat - delta: %s", time.Now().Sub(s.LastSeen).String())
	s.LastSeen = time.Now()
	thisMaster.Service.Servers[ipPort] = s
	thisMaster.Unlock()
}

func sendList(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	packets := thisMaster.Service.GeneratePackets(thisMaster.Options, p.Key)
	for _, v := range packets {
		_, err := (*conn).WriteTo(v, addr)
		if err != nil {
			LogServerAlert(ipPort, "error sending list packet [%s]", err)
			return
		}
	}
	LogServer(ipPort, "servers list sent")
}

func sendBanned(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	packets := thisMaster.BannedService.GeneratePackets(thisMaster.Options, p.Key)
	for _, v := range packets {
		_, err := (*conn).WriteTo(v, addr)
		if err != nil {
			LogServerAlert(ipPort, "error sending banned message packet [%s]", err)
			return
		}
	}
	LogServer(ipPort, "banned message sent")
}
