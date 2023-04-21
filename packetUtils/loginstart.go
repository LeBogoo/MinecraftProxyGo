package packetUtils

import (
	"bytes"
	"math"
	"minecraftproxy/packetUtils/utils"
)

type LoginStartPacket struct {
	Packet
	Name       string
	PlayerUUID string
}

func ParseLoginStartPacket(packet Packet) (LoginStartPacket, error) {
	reader := packet.GetReader()

	nameLength, _ := utils.ReadVarInt(reader)
	nameLength = int(math.Min(float64(nameLength), float64(16)))

	name, _ := utils.ReadString(reader, nameLength)

	hasPlayerUUID, _ := utils.ReadBool(reader)

	playerUUID := ""
	if hasPlayerUUID {
		playerUUID, _ = utils.ReadUUID(reader)
	}

	return LoginStartPacket{packet, name, playerUUID}, nil

}

func CreateLoginStartPacket(name string, uuid string) LoginStartPacket {
	loginStartPacket := LoginStartPacket{
		Packet: Packet{
			PacketLength: 0,
			PacketId:     0x00,
		},
		Name:       name,
		PlayerUUID: uuid,
	}

	return loginStartPacket
}

func (packet LoginStartPacket) ToBytes() ([]byte, error) {
	packetBuffer := bytes.NewBuffer([]byte{})
	packetBuffer.Write(utils.ToVarInt(packet.PacketId))
	packetBuffer.Write(utils.ToString(packet.Name))
	if packet.PlayerUUID != "" {
		packetBuffer.Write(utils.ToBool(true))
		packetBuffer.Write(utils.ToUUID(packet.PlayerUUID))
	} else {
		packetBuffer.Write(utils.ToBool(false))
	}

	packet.PacketLength = len(packetBuffer.Bytes())
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(utils.ToVarInt(packet.PacketLength))
	buffer.Write(packetBuffer.Bytes())

	return buffer.Bytes(), nil
}
