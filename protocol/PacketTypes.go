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

var packetTypeString = map[PacketType]string{
	PingInfoQuery:         "PingInfoQuery",
	PingInfoResponse:      "PingInfoResponse",
	MasterServerHeartbeat: "MasterServerHeartbeat",
	MasterServerList:      "MasterServerList",
	GameInfoQuery:         "GameInfoQuery",
	GameInfoResponse:      "GameInfoResponse",
}

func (p PacketType) String() string {
	return packetTypeString[PacketType(int(p)-int(PingInfoQuery))]
}
