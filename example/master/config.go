package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"net"
	"os"
	"sync"
)

type Configuration struct {
	sync.Mutex

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

	ColorLogs bool

	parsedBannedNets []*net.IPNet
}

const (
	DefaultConfigFileName = "mstrsvr.yaml"
	EnvPrefix             = "mstrsvr"
)

var config = new(Configuration)

func configInit() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(EnvPrefix)

	v.SetEnvPrefix(EnvPrefix)
	v.AllowEmptyEnv(true)

	v.SetDefault("ListenIP", "")
	v.SetDefault("ListenPort", 29000)
	v.SetDefault("MaxPacketSize", 512)
	v.SetDefault("MaxBufferSize", 32768)
	v.SetDefault("ServerTTL", 300)
	v.SetDefault("Hostname", "SlimThiccMaster")
	v.SetDefault("MOTD", "Welcome to Neo's MiniMaster")
	v.SetDefault("BannedMessage", "Welcome to bansville, population: you\\nVisit the discord to appeal!")
	v.SetDefault("ID", 99)
	v.SetDefault("ServersPerIP", 15)
	v.SetDefault("BannedNetworks", []string{"224.0.0.0/4"})
	v.SetDefault("ColorLogs", true)

	v.OnConfigChange(func(in fsnotify.Event) {
		LogComponent("config", "configuration change detected, updating...")
		rehashConfig(v)
	})
	v.WatchConfig()

	rehashConfig(v)

	loggerInit(config.ColorLogs)
}

func rehashConfig(v *viper.Viper) {
	err := v.ReadInConfig()
	if _, configFileNotFound := err.(viper.ConfigFileNotFoundError); err != nil && configFileNotFound {
		LogComponentAlert("config", "config file not found, creating...")
		err := v.WriteConfigAs(DefaultConfigFileName)
		if err != nil {
			LogComponentAlert("config", "unable to create config! [%s]", err)
			os.Exit(1)
		}
	} else if err != nil {
		LogComponentAlert("config", "error while reading config file [%s]", err)
	}

	config = new(Configuration)
	config.Lock()
	defer config.Unlock()
	err = v.Unmarshal(&config)
	if err != nil {
		LogComponentAlert("config", "error unmarshalling config [%s]", err)
	}

	config.parsedBannedNets = make([]*net.IPNet, 0)
	for _, v := range config.BannedNetworks {
		_, network, err := net.ParseCIDR(v)
		if err != nil {
			LogComponentAlert("config", "unable to parse BannedNetwork %s, %s", v, err)
			os.Exit(1)
		}
		config.parsedBannedNets = append(config.parsedBannedNets, network)
	}

	thisMaster.Service.MOTD = config.MOTD
	thisMaster.Service.MasterID = config.ID
	thisMaster.Service.CommonName = config.Hostname

	thisMaster.BannedService.MOTD = config.BannedMessage
	thisMaster.BannedService.MasterID = config.ID
	thisMaster.BannedService.CommonName = config.Hostname

	thisMaster.Options.MaxServerPacketSize = config.MaxPacketSize
}
