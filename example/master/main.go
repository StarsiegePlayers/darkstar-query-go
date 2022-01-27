package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var serviceRunning = true

var (
	VERSION string
	DATE    string
	TIME    string
	DEBUG   string
)

func init() {
	if DEBUG == "" {
		DEBUG = "Release"
	}
}

func main() {
	loggerInit(false)
	configInit()
	thisMaster.Options.LocalNetworks = config.localNetworks

	LogComponent("startup", "~~~ Neo's MiniMaster Starting Up ~~~")
	LogComponent("startup", "Version %s %s - Built on [%s@%s]", VERSION, DEBUG, DATE, TIME)

	maintenanceTimer := maintenanceInit()

	addrPort := fmt.Sprintf("%s:%d", config.ListenIP, config.ListenPort)

	pconn, err := net.ListenPacket("udp", addrPort)
	if err != nil {
		LogComponentAlert("server", "unable to bind to %s - [%s]", addrPort, err)
	}

	LogComponent("server", "now listening on [%s]", addrPort)

	// setup kill / rehash hooks
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for serviceRunning {
			sig := <-c
			LogComponentAlert("server", "received [%s]", sig.String())

			switch sig {
			case os.Interrupt:
				fallthrough
			case syscall.SIGTERM:
				LogComponentAlert("server", "shutdown initiated...")

				err = pconn.Close()
				if err != nil {
					log.Fatalln(err)
				}

				serviceRunning = false

				maintenanceShutdown(maintenanceTimer)
			}
		}
	}()

	// start listening loop
	buf := make([]byte, config.MaxPacketSize)
	buf2 := make([]byte, config.MaxPacketSize)
	prevIPPort := ""

	for serviceRunning {
		n, addr, err := pconn.ReadFrom(buf)
		if err != nil {
			t := &net.OpError{}
			if errors.As(err, &t) && t.Op == "read" {
				LogComponentAlert("server", "socket closed.")
				continue
			}
			LogComponentAlert("server", "read error on socket [%s]", err)
		}

		// dedupe packets because wtf dynamix
		if prevIPPort == addr.String() && bytes.Equal(buf2[:n], buf[:n]) {
			// blank out the stored header and discord the packet silently
			prevIPPort = ""
			continue
		}

		copy(buf2, buf)

		prevIPPort = addr.String()

		if addr, ok := addr.(*net.UDPAddr); ok {
			go serve(&pconn, addr, buf[:n])
		}
	}
	LogComponent("shutdown", "process complete")
}
