package packetUtils

import (
	"bytes"
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

func (packet PingRequestPacket) ToBytes() ([]byte, error) {
	packetBuffer := bytes.NewBuffer([]byte{})
	packetBuffer.Write(utils.ToVarInt(packet.PacketId))
	packetBuffer.Write(utils.ToLong(packet.Payload))

	packet.PacketLength = len(packetBuffer.Bytes())
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(utils.ToVarInt(packet.PacketLength))
	buffer.Write(packetBuffer.Bytes())

	return buffer.Bytes(), nil
}
