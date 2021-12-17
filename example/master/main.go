package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var serviceRunning = true

func main() {
	loggerInit()
	configInit()
	maintenanceTimer := maintenanceInit()

	addrPort := fmt.Sprintf("%s:%d", config.ListenIP, config.ListenPort)
	pconn, err := net.ListenPacket("udp", addrPort)
	if err != nil {
		LogComponentAlert("server", "unable to bind to %s - [%s]", addrPort, err)
	}
	LogComponent("server", "now listening on [%s]", addrPort)

	// setup control+c hook
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		LogComponentAlert("server", "shutdown initiated...")
		err = pconn.Close()
		if err != nil {
			log.Fatalln(err)
		}
		serviceRunning = false
		maintenanceShutdown(maintenanceTimer)
	}()

	// start listening loop
	buf := make([]byte, config.MaxPacketSize)
	buf2 := make([]byte, config.MaxPacketSize)
	for serviceRunning {
		n, addr, err := pconn.ReadFrom(buf)
		if err != nil {
			switch t := err.(type) {
			case *net.OpError:
				if t.Op == "read" {
					LogComponentAlert("server", "socket closed.")
				}
				continue
			default:
				LogComponentAlert("server", "read error on socket [%s]", err)
			}
		}

		// dedupe packets because wtf dynamix
		if bytes.Equal(buf, buf2) {
			continue
		}
		copy(buf2, buf)

		if addr, ok := addr.(*net.UDPAddr); ok {
			go serve(&pconn, addr, buf[:n])
		}
	}
}
