package main

import (
	"fmt"
	"github.com/StarsiegePlayers/darkstar-query-go/master"
	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
	"log"
	"net"
	"time"
)

var (
	service        = master.NewMaster()
	bannedService  = master.NewMaster()
	options        = protocol.Options{}
	ipServiceCount = make(map[string]int)
)

func serve(conn *net.PacketConn, addr *net.UDPAddr, buf []byte) {
	// we use an ip-port combo as a unique identifier
	ipPort := fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)

	// parse packet
	p := protocol.NewPacket()
	err := p.UnmarshalBinary(buf)
	if err != nil {
		switch err {
		case protocol.ErrorUnknownPacketVersion:
			log.Printf("[%s]: ! Unknown protocol number\n", ipPort)
		case protocol.ErrorEmptyPacket:
			log.Printf("[%s]: ! Empty packet recieved\n", ipPort)
		default:
			log.Printf("[%s]: ! Error %s while parsing packet\n", ipPort, err)
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
			log.Printf("[%s] ! Recieved a %s packet from banned host", addr.IP.String(), p.Type.String())
			return
		}
		registerHeartbeat(conn, addr, ipPort, p)
		break

	// client is requesting a server list
	case protocol.PingInfoQuery:
		if isBanned {
			log.Printf("[%s] ! Recieved a %s packet from banned host", addr.IP.String(), p.Type.String())
			sendBanned(conn, addr, ipPort, p)
			return
		}
		sendList(conn, addr, ipPort, p)
		break

	default:
		if isBanned {
			log.Printf("[%s]: ! Recieved unsociliated packet type %s from banned host", ipPort, p.Type.String())
			return
		}
		log.Printf("[%s]: ! Recieved unsociliated packet type %s", ipPort, p.Type.String())
	}
}

func registerHeartbeat(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	log.Printf("[%s]: Heartbeat - ", ipPort)
	s, exist := service.Servers[ipPort]
	if exist {
		log.Printf("delta: %s\n", time.Now().Sub(s.LastSeen).String())
		s.LastSeen = time.Now()
		return
	}
	log.Printf("New\n")
	service.Servers[ipPort] = &protocol.Server{
		Address:    addr,
		Connection: conn,
		LastSeen:   time.Now(),
	}
}

func sendList(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	log.Printf("[%s] - Servers list sent", ipPort)
	options.PacketKey = p.Key
	service.SendResponse(conn, addr, options)
}

func sendBanned(conn *net.PacketConn, addr *net.UDPAddr, ipPort string, p *protocol.Packet) {
	log.Printf("[%s] - Banned message sent", ipPort)
	options.PacketKey = p.Key
	bannedService.SendResponse(conn, addr, options)
}
