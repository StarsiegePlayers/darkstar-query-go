package protocol

import (
	"fmt"
	"math"
	"time"
)

type PingInfo struct {
	GameMode    byte          `csv:"-"` // ??
	PlayerCount byte          `csv:"cur_players"`
	MaxPlayers  byte          `csv:"max_players"`
	GameName    []byte        `csv:"-"` // es3a
	GameVersion []byte        `csv:"-"` // V 001.000r
	Name        []byte        `csv:"server_name"`
	GameStatus  StatusByte    `csv:"status"`
	Ping        time.Duration `csv:"ping"`
	*Packet     `json:"-" csv:"-"`
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
	s.PlayerCount = byte(p.ID & math.MaxUint8)
	s.MaxPlayers = byte((p.ID >> 8) & math.MaxUint8)

	s.GameName, p.Data = p.Data[0:4], p.Data[4:]
	s.GameStatus, p.Data = StatusByte(p.Data[0]), p.Data[1:]
	s.GameVersion, p.Data = p.Data[0:10], p.Data[10:]

	s.Name = p.Data[:Clen(p.Data)]

	return nil
}
