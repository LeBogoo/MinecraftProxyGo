package networking

import (
	"bufio"
	"encoding/json"
	"fmt"
	"minecraftproxy/config"
	"minecraftproxy/packetUtils"
	"net"
	"time"
)

func StatusPingServer(server *config.Server) (packetUtils.StatusResponse, error) {
	d := net.Dialer{Timeout: 250 * time.Millisecond}
	conn, err := d.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		return packetUtils.StatusResponse{}, err
	}

	defer conn.Close()

	handshakePacket := packetUtils.CreateHandshakePacket(server.Protocol, server.Host, server.Port, 1)
	handshakeBytes, _ := handshakePacket.ToBytes()

	conn.Write(handshakeBytes)

	statusRequestPacket := packetUtils.CreateStatusRequestPacket()
	statusRequestBytes, _ := statusRequestPacket.ToBytes()

	conn.Write(statusRequestBytes)

	reader := bufio.NewReader(conn)
	packet, _ := packetUtils.ParsePacket(reader)

	responsePacket, _ := packetUtils.ParseStatusResponsePacket(packet)

	var statusResponse packetUtils.StatusResponse
	json.Unmarshal(responsePacket.JSONResponse, &statusResponse)

	pingRequestPacket := packetUtils.CreatePingRequestPacket()
	pingRequestBytes, _ := pingRequestPacket.ToBytes()

	conn.Write(pingRequestBytes)

	packet, _ = packetUtils.ParsePacket(reader)
	packetUtils.ParsePingRequestPacket(packet)

	// Throw away the ping response packet. It is only implemented to complete the ping sequence.

	return statusResponse, nil
}
