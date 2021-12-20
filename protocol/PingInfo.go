package protocol

import (
	"fmt"
	"time"
)

type PingInfo struct {
	*Packet     `json:"-" csv:"-"`
	GameMode    byte          `csv:"-"` // ??
	GameName    string        `csv:"-"` // es3a
	GameVersion string        `csv:"-"` // V 001.000r
	GameStatus  StatusByte    `csv:"status"`
	PlayerCount byte          `csv:"cur_players"`
	MaxPlayers  byte          `csv:"max_players"`
	Name        string        `csv:"server_name"`
	Ping        time.Duration `csv:"ping"`
}

func (s *PingInfo) String() string {
	return fmt.Sprintf("PingInfoResponse: %s [%s] (%s) Players: %d/%d",
		s.Name, s.Ping, s.GameStatus, s.PlayerCount, s.MaxPlayers)
}

func (s *PingInfo) MarshalBinary() ([]byte, error) {
	// TODO
	return s.Packet.MarshalBinary()
}

func (s *PingInfo) UnmarshalBinary(data []byte) error {
	err := s.Packet.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	p := s.Packet
	s.GameMode = p.Total
	s.PlayerCount = byte(p.ID & 0xff)
	s.MaxPlayers = byte((p.ID >> 8) & 0xff)

	s.GameName, p.Data = string(p.Data[0:4]), p.Data[4:]
	s.GameStatus, p.Data = StatusByte(p.Data[0]), p.Data[1:]
	s.GameVersion, p.Data = string(p.Data[0:10]), p.Data[10:]

	s.Name = string(p.Data[:Clen(p.Data)])

	return nil
}
