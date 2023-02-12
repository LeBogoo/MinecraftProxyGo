package packetUtils

import (
	"minecraftproxy/packetUtils/utils"
)

type HandshakePacket struct {
	Packet
	ProtocolVersion int
	ServerAddress   string
	ServerPort      int
	NextState       int
}

func ParseHandshakePacket(packet Packet) (HandshakePacket, error) {
	reader := packet.GetReader()

	protocolVersion, err := utils.ReadVarInt(reader)
	addressLength, err := utils.ReadInt(reader)
	address, err := utils.ReadString(reader, addressLength)
	serverPort, err := utils.ReadUShort(reader)
	nextState, err := utils.ReadInt(reader)

	if err != nil {
		return HandshakePacket{}, err
	}

	return HandshakePacket{packet, protocolVersion, address, serverPort, nextState}, nil

}
