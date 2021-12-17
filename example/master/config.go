package main

import (
	"encoding/json"
	"net"
	"os"
)

type Configuration struct {
	ListenIP   string
	ListenPort uint16

	MaxPacketSize uint16
	MaxBufferSize uint16

	Hostname  string
	MOTD      string
	ServerTTL int

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
	ServerTTL:     300,
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
	thisMaster.Service.MOTD = config.MOTD
	thisMaster.Service.MasterID = config.ID
	thisMaster.Service.CommonName = config.Hostname

	thisMaster.BannedService.MOTD = config.BannedMessage
	thisMaster.BannedService.MasterID = config.ID
	thisMaster.BannedService.CommonName = config.Hostname

	thisMaster.Options.MaxServerPacketSize = config.MaxPacketSize

	for _, v := range config.BannedNetworks {
		_, network, err := net.ParseCIDR(v)
		if err != nil {
			LogComponent("config", "error unable to parse BannedNetwork %s, %s", v, err)
			os.Exit(1)
		}
		config.parsedBannedNets = append(config.parsedBannedNets, network)
	}

	configJson, _ := json.Marshal(config)
	LogComponent("config", string(configJson))
}
