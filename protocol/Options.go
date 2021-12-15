package protocol

import (
	"time"
)

type Options struct {
	Search  map[string]*Server
	Timeout time.Duration
	Debug   bool

	MaxServerPacketSize uint16
	PacketKey           uint16
}
