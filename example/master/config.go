package main

import (
	"log"
	"net"
)

type Configuration struct {
	ListenIP   string
	ListenPort uint16

	MaxPacketSize uint16
	MaxBufferSize uint16

	Hostname string

	MOTD string

	ID           uint16
	ServersPerIP uint16

	BannedNetworks []string
	BannedMessage  string

	parsedBannedNets []*net.IPNet
}

var config = Configuration{
	ListenIP:      "",
	ListenPort:    29000,
	MaxPacketSize: 512,
	MaxBufferSize: 32768,
	Hostname:      "SlimThiccMaster",
	MOTD:          "Welcome to Neo's MiniMaster",
	BannedMessage: "Welcome to bansville, population: you\\nVisit the discord to appeal!",
	ID:            99,
	ServersPerIP:  15,
	BannedNetworks: []string{
		"224.0.0.0/4",
	},
}

func configInit() {
	service.MOTD = config.MOTD
	service.MasterID = config.ID
	service.CommonName = config.Hostname

	bannedService.MOTD = config.BannedMessage
	bannedService.MasterID = config.ID
	bannedService.CommonName = config.Hostname

	options.MaxServerPacketSize = config.MaxPacketSize

	for _, v := range config.BannedNetworks {
		_, network, err := net.ParseCIDR(v)
		if err != nil {
			log.Fatalf("Error unable to parse BannedNetwork %s, %s", v, err)
		}
		config.parsedBannedNets = append(config.parsedBannedNets, network)
	}
}
