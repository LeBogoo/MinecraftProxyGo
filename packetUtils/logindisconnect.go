package packetUtils

import (
	"bytes"
	"encoding/json"
	"minecraftproxy/minecraft"
	"minecraftproxy/packetUtils/utils"
)

type LoginDisconnectPacket struct {
	Packet
	Reason string
}

func CreateLoginDisconnectPacket(reason minecraft.Text) LoginDisconnectPacket {
	jsonText, _ := json.Marshal(reason)

	return LoginDisconnectPacket{
		Packet: Packet{
			PacketLength: 0,
			PacketId:     0x00,
		},
		Reason: string(jsonText),
	}
}

func (packet LoginDisconnectPacket) ToBytes() ([]byte, error) {
	packetBuffer := bytes.NewBuffer([]byte{})
	packetBuffer.Write(utils.ToVarInt(packet.PacketId))
	packetBuffer.Write(utils.ToString(packet.Reason))

	packet.PacketLength = len(packetBuffer.Bytes())
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(utils.ToVarInt(packet.PacketLength))
	buffer.Write(packetBuffer.Bytes())

	return buffer.Bytes(), nil
}
