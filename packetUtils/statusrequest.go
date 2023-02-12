package packetUtils

type StatusRequestPacket struct {
	Packet
}

type StatusResponsePacket struct {
	Packet
	JSONResponse string
}

func ParseStatusRequestPacket(packet Packet) (StatusRequestPacket, error) {
	return StatusRequestPacket{packet}, nil
}
