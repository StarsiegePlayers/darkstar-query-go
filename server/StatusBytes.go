package server

import (
	"encoding/json"
	"strings"
)

type StatusBit int
type StatusByte byte

const (
	Protected StatusBit = 1 << iota
	Dedicated
	AllowOldClients
	Started
	Dynamix
	WON
	Unknown2
	Unknown3
)

var statusBitString = []string{
	"Protected",
	"Dedicated",
	"AllowOldClients",
	"Started",
	"Dynamix",
	"WON",
	"Unknown2",
	"Unknown3",
}

type StatusByteStruct struct {
	Protected       bool
	Dedicated       bool
	AllowOldClients bool
	Started         bool
	Dynamix         bool
	WON             bool
	Unknown2        bool
	Unknown3        bool
}

func (s StatusByteStruct) MarshalBinary() byte {
	output := byte(0)
	if s.Protected {
		output |= byte(Protected)
	}
	if s.Dedicated {
		output |= byte(Dedicated)
	}
	if s.AllowOldClients {
		output |= byte(AllowOldClients)
	}
	if s.Started {
		output |= byte(Started)
	}
	if s.Dynamix {
		output |= byte(Dynamix)
	}
	if s.WON {
		output |= byte(WON)
	}
	if s.Unknown2 {
		output |= byte(Unknown2)
	}
	if s.Unknown3 {
		output |= byte(Unknown3)
	}
	return output
}

func (s StatusByte) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Struct())
}

func (s StatusByte) Struct() StatusByteStruct {
	output := StatusByteStruct{}
	statusArr := make(map[StatusBit]bool)
	for i := 0; i < 8; i++ {
		bit := StatusBit(1 << i)
		if int(s)&int(bit) != 0 {
			statusArr[bit] = true
		} else {
			statusArr[bit] = false
		}
	}

	for k, v := range statusArr {
		switch k {
		case Protected:
			output.Protected = v
		case Dedicated:
			output.Dedicated = v
		case AllowOldClients:
			output.AllowOldClients = v
		case Started:
			output.Started = v
		case Dynamix:
			output.Dynamix = v
		case WON:
			output.WON = v
		case Unknown2:
			output.Unknown2 = v
		case Unknown3:
			output.Unknown3 = v
		}
	}
	return output
}

func (s StatusByte) StringSlice() []string {
	statusArr := make([]string, 0)
	for i := 0; i < 8; i++ {
		bit := StatusBit(1 << i)
		if int(s)&int(bit) != 0 {
			statusArr = append(statusArr, statusBitString[i])
		}
	}
	return statusArr
}

func (s StatusByte) String() string {
	statusArr := s.StringSlice()
	var out strings.Builder
	for k, v := range statusArr {
		if k == len(statusArr)-1 {
			out.WriteString(v)
			break
		}
		out.WriteString(v + " | ")
	}
	return out.String()
}
