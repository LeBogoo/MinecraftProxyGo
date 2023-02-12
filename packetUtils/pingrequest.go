package packetUtils

import (
	"minecraftproxy/packetUtils/utils"
)

type PingRequestPacket struct {
	Packet
	Payload int64
}

func ParsePingRequestPacket(packet Packet) (PingRequestPacket, error) {
	reader := packet.GetReader()

	payload, err := utils.ReadLong(reader)

	if err != nil {
		return PingRequestPacket{}, err
	}

	return PingRequestPacket{packet, payload}, nil

}
