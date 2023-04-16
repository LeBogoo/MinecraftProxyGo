package packetUtils

import (
	"bytes"
	"minecraftproxy/packetUtils/utils"
)

type HandshakePacket struct {
	Packet
	ProtocolVersion int
	ServerAddress   string
	ServerPort      int
	NextState       int
}

func CreateHandshakePacket(protocolVersion int, serverAddress string, serverPort int, nextState int) HandshakePacket {
	handshakePacket := HandshakePacket{
		Packet: Packet{
			PacketLength: 0,
			PacketId:     0x00,
		},
		ProtocolVersion: protocolVersion,
		ServerAddress:   serverAddress,
		ServerPort:      serverPort,
		NextState:       nextState,
	}

	return handshakePacket
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

func (packet HandshakePacket) ToBytes() ([]byte, error) {
	packetBuffer := bytes.NewBuffer([]byte{})
	packetBuffer.Write(utils.ToVarInt(packet.PacketId))
	packetBuffer.Write(utils.ToVarInt(packet.ProtocolVersion))
	packetBuffer.Write(utils.ToString(packet.ServerAddress))
	packetBuffer.Write(utils.ToUShort(packet.ServerPort))
	packetBuffer.Write(utils.ToVarInt(packet.NextState))

	packet.PacketLength = len(packetBuffer.Bytes())
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(utils.ToVarInt(packet.PacketLength))
	buffer.Write(packetBuffer.Bytes())

	return buffer.Bytes(), nil
}
