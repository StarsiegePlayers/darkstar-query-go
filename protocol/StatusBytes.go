package protocol

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
	Reserved1
	Reserved2
)

var statusBitString = map[StatusBit]string{
	Protected:       "Protected",
	Dedicated:       "Dedicated",
	AllowOldClients: "AllowOldClients",
	Started:         "Started",
	Dynamix:         "Dynamix",
	WON:             "WON",
	Reserved1:       "Reserved1",
	Reserved2:       "Reserved2",
}

type StatusByteStruct struct {
	Protected       bool
	Dedicated       bool
	AllowOldClients bool
	Started         bool
	Dynamix         bool
	WON             bool
	Reserved1       bool
	Reserved2       bool
}

func (s StatusByteStruct) MarshalBinary() (output byte) {
	// TODO: there's gotta be a better way to do this.
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

	if s.Reserved1 {
		output |= byte(Reserved1)
	}

	if s.Reserved2 {
		output |= byte(Reserved2)
	}

	return
}

func (s StatusByte) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Struct())
}

func (s StatusByte) Struct() (output StatusByteStruct) {
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
		case Reserved1:
			output.Reserved1 = v
		case Reserved2:
			output.Reserved2 = v
		}
	}

	return
}

func (s StatusByte) StringSlice() (statusArr []string) {
	for i := 0; i < 8; i++ {
		bit := StatusBit(1 << i)
		if int(s)&int(bit) != 0 {
			statusArr = append(statusArr, statusBitString[StatusBit(i)])
		}
	}

	return
}

func (s StatusByte) String() string {
	var (
		out       strings.Builder
		statusArr = s.StringSlice()
	)

	for k, v := range statusArr {
		if k == len(statusArr)-1 {
			out.WriteString(v)
			break
		}

		out.WriteString(v + " | ")
	}

	return out.String()
}
