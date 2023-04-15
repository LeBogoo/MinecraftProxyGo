package packetUtils

import (
	"bufio"
	"bytes"
	"minecraftproxy/packetUtils/utils"
)

type Packet struct {
	PacketLength int
	PacketId     int
	Reader       *bufio.Reader
}

func (packet Packet) GetPacketLength() int {
	return packet.PacketLength
}

func (packet Packet) GetPacketId() int {
	return packet.PacketId
}

func (packet Packet) GetReader() *bufio.Reader {
	return packet.Reader
}

func ParsePacket(reader *bufio.Reader) (Packet, error) {
	packetLength, err := utils.ReadInt(reader)
	if err != nil {
		return Packet{}, err
	}

	packetId, err := utils.ReadVarInt(reader)
	if err != nil {
		return Packet{}, err
	}

	return Packet{packetLength, packetId, reader}, nil
}

func (packet Packet) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(utils.ToVarInt(packet.PacketLength))
	buffer.Write(utils.ToVarInt(packet.PacketId))

	return buffer.Bytes(), nil
}
