package main

import (
	"fmt"
	"github.com/StarsiegePlayers/darkstar-query-go/master"
	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
	"net"
	"sync"
	"time"
)

type Service struct {
	sync.Mutex
	Service          *master.Master
	BannedService    *master.Master
	Options          *protocol.Options
	IPServiceCount   map[string]uint16
	SolicitedServers map[string]bool
}

var thisMaster = Service{
	Service:          master.NewMaster(),
	BannedService:    master.NewMaster(),
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

	case protocol.PingInfoResponse:
		if isBanned {
			LogServerAlert(addr.IP.String(), "Received a %s packet from banned host", p.Type.String())
			return
		}
		registerPingInfo(conn, addr, ipPort, p)
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
	thisMaster.SolicitedServers[ipPort] = true
	sendVerificationPacket(conn, addr, ipPort, p)
}

func registerPingInfo(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	thisMaster.Lock()
	defer thisMaster.Unlock()

	if _, exist := thisMaster.SolicitedServers[ipPort]; !exist {
		LogServerAlert(ipPort, "Received unsolicited packet type %s", ipPort)
		return
	}

	if _, exist := thisMaster.Service.Servers[ipPort]; !exist {
		count := thisMaster.IPServiceCount[addr.IP.String()]
		if uint16(count)+1 > config.ServersPerIP {
			LogServerAlert(ipPort, "Rejecting additional server for IP - count: %d/%d", count, config.ServersPerIP)
			return
		}

		// log and add new
		LogServer(ipPort, "Heartbeat - New Server")
		thisMaster.Service.Servers[ipPort] = &protocol.Server{
			Address:    addr,
			Connection: conn,
			LastSeen:   time.Now(),
			TTL:        config.ServerTTL,
		}
		count++
		LogServer(ipPort, "New Server for IP - total server count for IP: %d/%d", count, config.ServersPerIP)
		thisMaster.IPServiceCount[addr.IP.String()] = count
	}

	s := thisMaster.Service.Servers[ipPort]
	LogServer(ipPort, "Heartbeat - delta: %s", time.Now().Sub(s.LastSeen).String())
	s.LastSeen = time.Now()
	thisMaster.Service.Servers[ipPort] = s
	thisMaster.Service.Servers[ipPort].Info = p
}

func sendVerificationPacket(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	response := protocol.NewPacket()
	response.ID = config.ID
	response.Key = p.Key
	response.Type = protocol.PingInfoQuery

	data, err := response.MarshalBinary()
	if err != nil {
		LogServerAlert(ipPort, "error marshalling verification response - [%s]", err)
		return
	}

	// send two packets because wtf dynamix
	for i := 0; i < 2; i++ {
		_, err = (*conn).WriteTo(data, addr)
		if err != nil {
			LogServerAlert(ipPort, "error sending verification response - [%s]", ipPort, err)
			return
		}
	}
}

func sendList(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	LogServer(ipPort, "Servers list sent")
	thisMaster.Options.PacketKey = p.Key
	thisMaster.Service.SendResponse(conn, addr, thisMaster.Options)
}

func sendBanned(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	LogServer(ipPort, "Banned message sent", ipPort)
	thisMaster.Options.PacketKey = p.Key
	thisMaster.BannedService.SendResponse(conn, addr, thisMaster.Options)
}
