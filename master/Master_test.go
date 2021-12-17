package master

import (
	"bytes"
	"fmt"
	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
	"github.com/stretchr/testify/suite"
	"log"
	"net"
	"testing"
)

type MasterTestSite struct {
	suite.Suite
	Master *Master
}

func (t *MasterTestSite) SetupTest() {
	t.Master = &Master{
		Address:    "localhost.localdomain",
		CommonName: "\\nDummythicc Masterserver Testing",
		MOTD:       "Welcome to a Testing server for Neo's Dummythiccness",
		MOTDJunk:   "dummythicc",
		Servers:    make(map[string]*protocol.Server),
		MasterID:   99,
	}
}

/********************************************************************/
// TestMaster_SendResponse_1Pkt tests if the master server
// produces the correct one packet response

func (t MasterTestSite) TestMaster_SendResponse_1Pkt() {
	response := [][]byte{
		{
			0x10, 0x06, 0x01, 0x01, 0x00, 0x00, 0x00, 0x45, 0x21, 0x5C, 0x6E, 0x44, 0x75, 0x6D, 0x6D, 0x79,
			0x74, 0x68, 0x69, 0x63, 0x63, 0x20, 0x4D, 0x61, 0x73, 0x74, 0x65, 0x72, 0x73, 0x65, 0x72, 0x76,
			0x65, 0x72, 0x20, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6E, 0x67, 0x3E, 0x64, 0x75, 0x6D, 0x6D, 0x79,
			0x74, 0x68, 0x69, 0x63, 0x63, 0x57, 0x65, 0x6C, 0x63, 0x6F, 0x6D, 0x65, 0x20, 0x74, 0x6F, 0x20,
			0x61, 0x20, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6E, 0x67, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
			0x20, 0x66, 0x6F, 0x72, 0x20, 0x4E, 0x65, 0x6F, 0x27, 0x73, 0x20, 0x44, 0x75, 0x6D, 0x6D, 0x79,
			0x74, 0x68, 0x69, 0x63, 0x63, 0x6E, 0x65, 0x73, 0x73, 0x00, 0x39, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x49, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x4C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x4E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x50, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x51, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x52,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x53, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x54, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x55, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x56, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x57, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x58, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x59, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x5C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x5E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x60, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x61, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x62,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x63, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x64, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x65, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x66, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x67, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x68, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x69, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x6C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x6E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x70, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x71, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x72,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x73, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x74, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x75, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x76, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x77, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x78, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x79, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x7C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x7E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x80, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x81, 0x71,
		},
	}
	t.Master.Servers = generateAddresses(57)
	t.Master.MasterID = 69
	options := newOptions()

	server, err := net.ListenPacket("udp", "127.0.0.1:42071")
	if err != nil {
		t.T().Fatalf("Error duing listen - %s", err)
	}
	defer server.Close()

	resultPipe := make(chan []byte, 1)
	go listeningClient(resultPipe, server)

	t.Master.SendResponse(&server, server.LocalAddr().(*net.UDPAddr), options)
	packets := make([][]byte, 0)

	for i := 0; i < 1; i++ {
		r := <-resultPipe
		packets = append(packets, r)
		t.Assert().Equal(response[i], packets[i])
	}
}

/********************************************************************/
// TestMaster_SendResponse_2Pkts tests if the master server
// produces the correct spanned packet response

func (t MasterTestSite) TestMaster_SendResponse_2Pkts() {
	response := [][]byte{
		{
			0x10, 0x06, 0x01, 0x02, 0x00, 0x00, 0x00, 0x63, 0x21, 0x5C, 0x6E, 0x44, 0x75, 0x6D, 0x6D, 0x79,
			0x74, 0x68, 0x69, 0x63, 0x63, 0x20, 0x4D, 0x61, 0x73, 0x74, 0x65, 0x72, 0x73, 0x65, 0x72, 0x76,
			0x65, 0x72, 0x20, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6E, 0x67, 0x3E, 0x64, 0x75, 0x6D, 0x6D, 0x79,
			0x74, 0x68, 0x69, 0x63, 0x63, 0x57, 0x65, 0x6C, 0x63, 0x6F, 0x6D, 0x65, 0x20, 0x74, 0x6F, 0x20,
			0x61, 0x20, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6E, 0x67, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
			0x20, 0x66, 0x6F, 0x72, 0x20, 0x4E, 0x65, 0x6F, 0x27, 0x73, 0x20, 0x44, 0x75, 0x6D, 0x6D, 0x79,
			0x74, 0x68, 0x69, 0x63, 0x63, 0x6E, 0x65, 0x73, 0x73, 0x00, 0x39, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x49, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x4C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x4E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x50, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x51, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x52,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x53, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x54, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x55, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x56, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x57, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x58, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x59, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x5C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x5E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x5F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x60, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x61, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x62,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x63, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x64, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x65, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x66, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x67, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x68, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x69, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x6C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x6E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x6F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x70, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x71, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x72,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x73, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x74, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x75, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x76, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x77, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x78, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x79, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7B, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x7C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7D, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x7E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x7F, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x80, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x81, 0x71,
		},

		{
			0x10, 0x06, 0x02, 0x02, 0x00, 0x00, 0x00, 0x63, 0x47, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x82, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x83, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x84, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x85, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x86, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x87, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x88, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x89,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x8A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x8B, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x8C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x8D, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x8E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x8F, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0x90, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x91, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x92, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0x93, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x94, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0x95, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x96, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0x97, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x98, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x99,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x9A, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x9B, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0x9C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x9D, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0x9E, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x9F, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0xA0, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xA1, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xA2, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0xA3, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xA4, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0xA5, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xA6, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0xA7, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xA8, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xA9,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xAA, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xAB, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0xAC, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xAD, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0xAE, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xAF, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0xB0, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xB1, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xB2, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0xB3, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xB4, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0xB5, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xB6, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0xB7, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xB8, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xB9,
			0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xBA, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xBB, 0x71, 0x06,
			0x7F, 0x00, 0x00, 0x01, 0xBC, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xBD, 0x71, 0x06, 0x7F, 0x00,
			0x00, 0x01, 0xBE, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xBF, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01,
			0xC0, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xC1, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xC2, 0x71,
			0x06, 0x7F, 0x00, 0x00, 0x01, 0xC3, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xC4, 0x71, 0x06, 0x7F,
			0x00, 0x00, 0x01, 0xC5, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xC6, 0x71, 0x06, 0x7F, 0x00, 0x00,
			0x01, 0xC7, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0xC8, 0x71,
		},
	}
	t.Master.Servers = generateAddresses(128)
	options := newOptions()

	server, err := net.ListenPacket("udp", "127.0.0.1:42070")
	if err != nil {
		t.T().Fatalf("Error duing listen - %s", err)
	}
	defer server.Close()

	resultPipe := make(chan []byte, 1)
	go listeningClient(resultPipe, server)

	t.Master.SendResponse(&server, server.LocalAddr().(*net.UDPAddr), options)

	packets := make([][]byte, 0)

	for i := 0; i < 2; i++ {
		r := <-resultPipe
		packets = append(packets, r)
		t.Assert().Equal(response[i], packets[i])
	}
}

/********************************************************************/
// TestMaster_MarshalPacket_NoMOTD tests if we can correctly create a packet without an MOTD (ala MiniMaster)

func (t MasterTestSite) TestMaster_MarshalPacket_NoMOTD_Header() {
	t.Master = &Master{
		Address:    "localhost.localdomain",
		CommonName: "MiniMaster",
		MOTD:       "",
		Servers:    protocol.NewServersMapFromList([]string{}),
		MasterID:   99,
	}
	options := newOptions()
	options.PacketKey = 0x0300

	response := [][]byte{
		{
			0x10, 0x06, 0x01, 0x01, 0x03, 0x00, 0x00, 0x63, 0x0A, 0x4D, 0x69, 0x6E, 0x69, 0x4D, 0x61, 0x73,
			0x74, 0x65, 0x72,
		},
	}

	server, err := net.ListenPacket("udp", "127.0.0.1:42069")
	if err != nil {
		t.T().Fatalf("Error duing listen - %s", err)
	}
	defer server.Close()

	resultPipe := make(chan []byte, 1)
	go listeningClient(resultPipe, server)

	t.Master.SendResponse(&server, server.LocalAddr().(*net.UDPAddr), options)

	packets := make([][]byte, 0)

	for i := 0; i < 1; i++ {
		r := <-resultPipe
		packets = append(packets, r)
		t.Assert().Equal(response[i], packets[i])
	}
}

/********************************************************************/
// TestMaster_UnmarshalBinary tests if we can correctly parse
// packets generated by other master servers

func (t MasterTestSite) TestMaster_UnmarshalBinary() {
	t.Master = &Master{
		Address:    "localhost.localdomain",
		CommonName: "Master2.Starsiege.pw",
		MOTD:       "Join the Starsiege Discord at discord.gg\\TNm4s2p",
		Servers:    protocol.NewServersMapFromList([]string{"154.0.175.219:29008", "184.89.64.182:29001", "154.0.175.219:29010", "154.0.175.219:29004", "154.0.175.219:29003", "154.0.175.219:29001", "154.0.175.219:29005", "154.0.175.219:29007", "154.0.175.219:29002", "154.0.175.219:29012", "154.0.175.219:29006", "192.155.86.254:29009", "154.0.175.219:29011", "192.155.86.254:29010", "192.155.86.254:29008", "192.155.86.254:29007", "192.155.86.254:29006", "192.155.86.254:29005", "192.155.86.254:29004", "192.155.86.254:29003", "192.155.86.254:29002", "192.155.86.254:29001", "96.126.117.157:29004", "96.126.117.157:29003", "96.126.117.157:29002", "96.126.117.157:29001", "154.0.175.219:29009"}),
		MasterID:   0x02,
	}
	response := []byte{
		0x10, 0x06, 0x01, 0x01, 0x45, 0x00, 0x00, 0x02, 0x16, 0x5C, 0x6E, 0x4D, 0x61, 0x73, 0x74, 0x65,
		0x72, 0x32, 0x2E, 0x53, 0x74, 0x61, 0x72, 0x73, 0x69, 0x65, 0x67, 0x65, 0x2E, 0x70, 0x77, 0x3B,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x4A, 0x6F, 0x69, 0x6E, 0x20, 0x74,
		0x68, 0x65, 0x20, 0x53, 0x74, 0x61, 0x72, 0x73, 0x69, 0x65, 0x67, 0x65, 0x20, 0x44, 0x69, 0x73,
		0x63, 0x6F, 0x72, 0x64, 0x5C, 0x6E, 0x61, 0x74, 0x20, 0x64, 0x69, 0x73, 0x63, 0x6F, 0x72, 0x64,
		0x2E, 0x67, 0x67, 0x5C, 0x54, 0x4E, 0x6D, 0x34, 0x73, 0x32, 0x70, 0x00, 0x1F, 0x06, 0x9A, 0x00,
		0xAF, 0xDB, 0x50, 0x71, 0x06, 0xB8, 0x59, 0x40, 0xB6, 0x49, 0x71, 0x06, 0x9A, 0x00, 0xAF, 0xDB,
		0x52, 0x71, 0x06, 0x9A, 0x00, 0xAF, 0xDB, 0x4C, 0x71, 0x06, 0x9A, 0x00, 0xAF, 0xDB, 0x4B, 0x71,
		0x06, 0x9A, 0x00, 0xAF, 0xDB, 0x49, 0x71, 0x06, 0x9A, 0x00, 0xAF, 0xDB, 0x4D, 0x71, 0x06, 0x9A,
		0x00, 0xAF, 0xDB, 0x4F, 0x71, 0x06, 0x9A, 0x00, 0xAF, 0xDB, 0x4A, 0x71, 0x06, 0x9A, 0x00, 0xAF,
		0xDB, 0x54, 0x71, 0x06, 0x9A, 0x00, 0xAF, 0xDB, 0x4E, 0x71, 0x06, 0xC0, 0x9B, 0x56, 0xFE, 0x51,
		0x71, 0x06, 0x9A, 0x00, 0xAF, 0xDB, 0x53, 0x71, 0x06, 0xC0, 0x9B, 0x56, 0xFE, 0x52, 0x71, 0x06,
		0xC0, 0x9B, 0x56, 0xFE, 0x50, 0x71, 0x06, 0xC0, 0x9B, 0x56, 0xFE, 0x4F, 0x71, 0x06, 0xC0, 0x9B,
		0x56, 0xFE, 0x4E, 0x71, 0x06, 0xC0, 0x9B, 0x56, 0xFE, 0x4D, 0x71, 0x06, 0xC0, 0x9B, 0x56, 0xFE,
		0x4C, 0x71, 0x06, 0xC0, 0x9B, 0x56, 0xFE, 0x4B, 0x71, 0x06, 0xC0, 0x9B, 0x56, 0xFE, 0x4A, 0x71,
		0x06, 0xC0, 0x9B, 0x56, 0xFE, 0x49, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4C, 0x71, 0x06, 0x60,
		0x7E, 0x75, 0x9D, 0x4C, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4B, 0x71, 0x06, 0x60, 0x7E, 0x75,
		0x9D, 0x4B, 0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x4A, 0x71, 0x06, 0x60, 0x7E, 0x75, 0x9D, 0x4A,
		0x71, 0x06, 0x7F, 0x00, 0x00, 0x01, 0x49, 0x71, 0x06, 0x60, 0x7E, 0x75, 0x9D, 0x49, 0x71, 0x06,
		0x9A, 0x00, 0xAF, 0xDB, 0x51, 0x71,
	}

	p := new(protocol.Packet)
	err := p.UnmarshalBinary(response)
	t.Assert().Nil(err)

	t.T().Logf("Type %s - Key %d - ID %d - Number %d - Total %d", p.Type, p.Key, p.ID, p.Number, p.Total)

	m := NewMaster()
	err = m.UnmarshalBinary(p)
	t.Assert().Nil(err)

	t.Assert().Equal(t.Master.CommonName, m.CommonName)
	t.Assert().Equal(t.Master.MOTD, m.MOTD)
	t.Assert().Equal(t.Master.id, m.id)

	for k := range t.Master.Servers {
		t.Assert().Equal(t.Master.Servers[k].String(), m.Servers[k].String())
	}
}

/********************************************************************/
// TestMaster_UnmarshalBinary_Minimaster tests if we can correctly parse
// packets generated by other master servers (minimaster)

func (t MasterTestSite) TestMaster_UnmarshalBinary_Minimaster() {
	t.Master = &Master{
		Address:    "localhost.localdomain",
		CommonName: "MiniMaster",
		MOTD:       "",
		Servers:    protocol.NewServersMapFromList([]string{"184.61.79.247:29002", "73.20.252.165:29002", "65.29.146.31:29003", "96.126.117.157:29002", "96.126.117.157:29004", "96.126.117.157:29005", "96.126.117.157:29007", "65.29.146.31:29001", "73.20.252.165:29001", "65.29.146.31:29002", "96.126.117.157:29001", "96.126.117.157:29003", "96.126.117.157:29006", "146.115.168.117:29001", "184.61.79.247:29001", "184.61.79.247:29003"}),
		MasterID:   99,
	}
	response := []byte{
		0x10, 0x06, 0x01, 0x01, 0x03, 0x00, 0x00, 0x63, 0x0A, 0x4D, 0x69, 0x6E, 0x69, 0x4D, 0x61, 0x73,
		0x74, 0x65, 0x72, 0x00, 0x00, 0x10, 0x06, 0xB8, 0x3D, 0x4F, 0xF7, 0x49, 0x71, 0x06, 0xB8, 0x3D,
		0x4F, 0xF7, 0x4A, 0x71, 0x06, 0xB8, 0x3D, 0x4F, 0xF7, 0x4B, 0x71, 0x06, 0x49, 0x14, 0xFC, 0xA5,
		0x49, 0x71, 0x06, 0x49, 0x14, 0xFC, 0xA5, 0x4A, 0x71, 0x06, 0x41, 0x1D, 0x92, 0x1F, 0x4A, 0x71,
		0x06, 0x60, 0x7E, 0x75, 0x9D, 0x49, 0x71, 0x06, 0x60, 0x7E, 0x75, 0x9D, 0x4A, 0x71, 0x06, 0x60,
		0x7E, 0x75, 0x9D, 0x4B, 0x71, 0x06, 0x60, 0x7E, 0x75, 0x9D, 0x4C, 0x71, 0x06, 0x60, 0x7E, 0x75,
		0x9D, 0x4D, 0x71, 0x06, 0x60, 0x7E, 0x75, 0x9D, 0x4E, 0x71, 0x06, 0x60, 0x7E, 0x75, 0x9D, 0x4F,
		0x71, 0x06, 0x41, 0x1D, 0x92, 0x1F, 0x49, 0x71, 0x06, 0x41, 0x1D, 0x92, 0x1F, 0x4B, 0x71, 0x06,
		0x92, 0x73, 0xA8, 0x75, 0x49, 0x71,
	}

	p := new(protocol.Packet)
	err := p.UnmarshalBinary(response)
	t.Assert().Nil(err)

	t.T().Logf("Type %s - Key %d - ID %d - Number %d - Total %d", p.Type, p.Key, p.ID, p.Number, p.Total)

	m := NewMaster()
	err = m.UnmarshalBinary(p)
	t.Assert().Nil(err)

	t.Assert().Equal(t.Master.CommonName, m.CommonName)
	t.Assert().Equal(t.Master.MOTD, m.MOTD)
	t.Assert().Equal(t.Master.id, m.id)

	output := make([]string, 0)
	for k, _ := range m.Servers {
		output = append(output, k)
	}
	t.T().Log(output)

	for k := range t.Master.Servers {
		t.Assert().Equal(t.Master.Servers[k].String(), m.Servers[k].String())
	}
}

/********************************************************************/
// Utility Functions

func generateAddresses(count int) map[string]*protocol.Server {
	startPort := 29001
	output := make(map[string]*protocol.Server)
	for i := 0; i < count; i++ {
		addrPort := fmt.Sprintf("127.0.0.1:%d", startPort+i)
		thisServer, _ := protocol.NewServerFromString(addrPort, 300)
		output[addrPort] = thisServer
	}
	return output
}

func newOptions() *protocol.Options {
	return &protocol.Options{
		Search:              nil,
		Timeout:             0,
		Debug:               false,
		MaxServerPacketSize: 512,
	}
}

func listeningClient(resultPipe chan []byte, conn net.PacketConn) {
	defer close(resultPipe)
	for {
		data := make([]byte, protocol.MaxPacketSize)
		n, _, err := conn.ReadFrom(data)
		if err != nil {
			resultPipe <- []byte{}
			log.Fatalln(err)
		}

		pkt := protocol.NewPacket()
		err = pkt.UnmarshalBinary(data[:n])
		if err != nil {
			log.Fatalln(err)
		}

		data = bytes.TrimRight(data, "\x00")
		resultPipe <- data

		if pkt.Number == pkt.Total || pkt.Number == 0xFF {
			break
		}
	}
}

/********************************************************************/
// Entrypoint

func TestMasterTestSuite(t *testing.T) {
	suite.Run(t, new(MasterTestSite))
}
