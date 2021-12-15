package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	configInit()

	addrPort := fmt.Sprintf("%s:%d", config.ListenIP, config.ListenPort)
	pconn, err := net.ListenPacket("udp", addrPort)
	if err != nil {
		log.Fatalf("Unable to bind to %s\n%s", addrPort, err)
	}
	log.Println("Now listening on ", addrPort)
	defer pconn.Close()

	for {
		buf := make([]byte, config.MaxPacketSize)
		n, addr, err := pconn.ReadFrom(buf)
		if err != nil {
			log.Println("read error from", addr.String())
			continue
		}

		if addr, ok := addr.(*net.UDPAddr); ok {
			go serve(&pconn, addr, buf[:n])
		}

	}
}
