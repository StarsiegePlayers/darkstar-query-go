package protocol

import "time"

type Options struct {
	Search  []string
	Timeout time.Duration
	Debug   bool

	MaxServerPacketSize uint16
}
