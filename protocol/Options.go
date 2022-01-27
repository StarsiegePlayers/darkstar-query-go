package protocol

import (
	"net"
	"time"
)

type Options struct {
	Timeout       time.Duration
	LocalNetworks []*net.IPNet
	ExternalIP    net.IP
	Debug         bool

	MaxServerPacketSize  uint16
	MaxNetworkPacketSize uint16
}
