package packetUtils

import (
	"bytes"
	"encoding/json"
	"minecraftproxy/packetUtils/utils"
)

type StatusRequestPacket struct {
	Packet
}

type StatusResponsePacket struct {
	Packet
	JSONResponse []byte
}

type StatusResponse struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        Description `json:"description"`
	Favicon            string      `json:"favicon"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
}

type Description struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

type Players struct {
	Max    int64    `json:"max"`
	Online int64    `json:"online"`
	Sample []Sample `json:"sample"`
}

type Sample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int64  `json:"protocol"`
}

func ParseStatusRequestPacket(packet Packet) (StatusRequestPacket, error) {
	return StatusRequestPacket{packet}, nil
}

func ParseStatusResponsePacket(packet Packet) (StatusResponsePacket, error) {
	reader := packet.GetReader()

	statusJsonLength, _ := utils.ReadVarInt(reader)
	statusJSON, _ := utils.ReadString(reader, statusJsonLength)

	return StatusResponsePacket{packet, []byte(statusJSON)}, nil
}

func CreateStatusRequestPacket() StatusRequestPacket {
	statusRequestPacket := StatusRequestPacket{
		Packet: Packet{
			PacketLength: 0,
			PacketId:     0x00,
		},
	}

	return statusRequestPacket
}

func CreateStatusResponsePacket(statusResponse StatusResponse) StatusResponsePacket {
	jsonString, _ := json.Marshal(statusResponse)
	jsonLength := len(jsonString)
	varIntLength := utils.ToVarInt(jsonLength)

	statusRequestPacket := StatusResponsePacket{
		Packet: Packet{
			PacketLength: 0,
			PacketId:     0x00,
		},
		JSONResponse: append(varIntLength, jsonString...),
	}

	return statusRequestPacket
}

func (packet StatusRequestPacket) ToBytes() ([]byte, error) {
	packetBuffer := bytes.NewBuffer([]byte{})
	packetBuffer.Write(utils.ToVarInt(packet.PacketId))

	packet.PacketLength = len(packetBuffer.Bytes())
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(utils.ToVarInt(packet.PacketLength))
	buffer.Write(packetBuffer.Bytes())

	return buffer.Bytes(), nil
}

func (packet StatusResponsePacket) ToBytes() ([]byte, error) {
	packetBuffer := bytes.NewBuffer([]byte{})
	packetBuffer.Write(utils.ToVarInt(packet.PacketId))
	packetBuffer.Write(packet.JSONResponse)

	packet.PacketLength = len(packetBuffer.Bytes())
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(utils.ToVarInt(packet.PacketLength))
	buffer.Write(packetBuffer.Bytes())

	return buffer.Bytes(), nil
}
