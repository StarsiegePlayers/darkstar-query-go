package main

import (
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Configuration struct {
	sync.Mutex

	ColorLogs bool

	ListenPort    uint16
	MaxPacketSize uint16
	MaxBufferSize uint16
	ID            uint16
	ServersPerIP  uint16

	ListenIP  string
	Hostname  string
	MOTD      string
	ServerTTL string

	MaintenanceInterval string
	BannedMessage       string
	BannedNetworks      []string

	parsedBannedNets []*net.IPNet
	localNetworks    []*net.IPNet
	serverTimeout    time.Duration
	maintenanceTimer time.Duration
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
	v.SetDefault("ServerTTL", 5*time.Minute)
	v.SetDefault("MaintenanceInterval", time.Minute)
	v.SetDefault("Hostname", "SlimThiccMaster")
	v.SetDefault("MOTD", "Welcome to Neo's MiniMaster")
	v.SetDefault("BannedMessage", "Welcome to bansville, population: you\\nVisit the discord to appeal!")
	v.SetDefault("ID", 99)
	v.SetDefault("ServersPerIP", 15)
	v.SetDefault("BannedNetworks", []string{"224.0.0.0/4"})
	v.SetDefault("ColorLogs", true)

	v.OnConfigChange(func(in fsnotify.Event) {
		if in.Op == fsnotify.Write {
			LogComponent("config", "configuration change detected, updating...")
			rehashConfig(v)
		}
	})

	rehashConfig(v)

	loggerInit(config.ColorLogs)
	v.WatchConfig()
}

func rehashConfig(v *viper.Viper) {
	err := v.ReadInConfig()
	if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) {
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
	err = v.Unmarshal(&config)
	if err != nil {
		LogComponentAlert("config", "error unmarshalling config [%s]", err)
	}

	config.Lock()

	config.parsedBannedNets = make([]*net.IPNet, 0)
	for _, v := range config.BannedNetworks {
		_, network, err := net.ParseCIDR(v)
		if err != nil {
			LogComponentAlert("config", "unable to parse BannedNetwork %s, %s", v, err)
			os.Exit(1)
		}

		config.parsedBannedNets = append(config.parsedBannedNets, network)
	}

	config.serverTimeout, err = time.ParseDuration(config.ServerTTL)
	if err != nil {
		LogComponentAlert("config", "unable to parse ServerTimeout, defaulting to 5 minutes")

		config.serverTimeout = 5 * time.Minute
	}

	config.maintenanceTimer, err = time.ParseDuration(config.MaintenanceInterval)
	if err != nil {
		LogComponentAlert("config", "unable to parse MaintenanceInterval, defaulting to 60 seconds")

		config.serverTimeout = time.Minute
	}

	thisMaster.Lock()
	thisMaster.Service.MOTD = config.MOTD
	thisMaster.Service.MasterID = config.ID
	thisMaster.Service.CommonName = config.Hostname

	thisMaster.BannedService.MOTD = config.BannedMessage
	thisMaster.BannedService.MasterID = config.ID
	thisMaster.BannedService.CommonName = config.Hostname

	thisMaster.Options.MaxServerPacketSize = config.MaxPacketSize
	thisMaster.Unlock()

	config.localNetworks = generateLocalAddresses()
	config.Unlock()
}

func generateLocalAddresses() (output []*net.IPNet) {
	output = make([]*net.IPNet, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for k := range addrs {
			output = append(output, addrs[k].(*net.IPNet))
		}
	}

	return
}
