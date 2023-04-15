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
	Text string `json:"text"`
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

	statusRequestPacket.PacketLength = int(len(statusRequestPacket.JSONResponse) + 1)
	return statusRequestPacket
}

func (packet StatusResponsePacket) ToBytes() ([]byte, error) {
	parentBytes, err := packet.Packet.ToBytes()

	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(parentBytes)
	buffer.Write(packet.JSONResponse)

	return buffer.Bytes(), nil
}
