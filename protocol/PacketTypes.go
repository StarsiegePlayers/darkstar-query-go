package protocol

type PacketType int

const RequestAllPackets = 0xff

const (
	PingInfoQuery PacketType = iota + 0x03
	PingInfoResponse
	MasterServerHeartbeat
	MasterServerList
	GameInfoQuery
	GameInfoResponse
)

var packetTypeString = []string{
	"PingInfoQuery",
	"PingInfoResponse",
	"MasterServerHeartbeat",
	"MasterServerList",
	"GameInfoQuery",
	"GameInfoResponse",
}

// String returns the packet type string based on the defined packet types
func (p PacketType) String() string {
	return packetTypeString[int(p)-int(PingInfoQuery)]
}
