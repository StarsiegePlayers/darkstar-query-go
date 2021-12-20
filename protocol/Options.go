package protocol

import (
	"time"
)

type Options struct {
	Timeout time.Duration
	Debug   bool

	MaxServerPacketSize  uint16
	MaxNetworkPacketSize uint16
}
