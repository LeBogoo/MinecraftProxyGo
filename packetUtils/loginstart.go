package packetUtils

import (
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
