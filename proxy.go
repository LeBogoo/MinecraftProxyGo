package main

import (
	"bufio"
	"fmt"
	"net"

	"minecraftproxy/config"
	"minecraftproxy/networking"
	"minecraftproxy/packetUtils"
)

func handleConnection(conn net.Conn, config *config.Config) {
	state := 0

	reader := bufio.NewReader(conn)

	defer conn.Close()

	for {
		// check if connection is closed
		if conn == nil {
			break
		}

		packet, err := packetUtils.ParsePacket(reader)
		fmt.Println("*************************")
		fmt.Println("PacketLength:", packet.PacketLength)
		fmt.Println("PacketId:", packet.PacketId)
		fmt.Println("State:", state)

		if (state == 0) && (packet.PacketId == 0x00) {
			handshakePacket, _ := packetUtils.ParseHandshakePacket(packet)
			state = handshakePacket.NextState
			continue
		}

		if (state == 1) && (packet.PacketId == 0x00) {
			fmt.Println("StatusRequestPacket")

			statusResponse, err := networking.StatusPingServer(&config.Server)
			var response packetUtils.StatusResponsePacket

			if err != nil {
				response = packetUtils.CreateStatusResponsePacket(config.OfflineStatusResponse)
			} else {
				response = packetUtils.CreateStatusResponsePacket(statusResponse)

			}

			bytes, _ := response.ToBytes()
			conn.Write(bytes)
		}

		if (state == 1) && (packet.PacketId == 0x01) {
			fmt.Println("PingPacket")
			pingPacket, _ := packetUtils.ParsePingRequestPacket(packet)
			fmt.Println("Payload:", pingPacket.Payload)

			response, _ := pingPacket.ToBytes()

			conn.Write(response)
			conn.Close()
			fmt.Println("Connection closed by server")
			break

		}

		if (state == 2) && (packet.PacketId == 0x00) {
			fmt.Println("LoginStartPacket")
			loginStartPacket, _ := packetUtils.ParseLoginStartPacket(packet)
			fmt.Println("Username:", loginStartPacket.Name)
			fmt.Println("UUID:", loginStartPacket.PlayerUUID)

			disconnectPacket := packetUtils.CreateLoginDisconnectPacket("{\"text\":\"Starting...\n\",\"bold\":true,\"color\":\"#00ff00\",\"extra\":[{\"color\":\"white\",\"bold\":false,\"text\":\"" + config.StartingDisconnectMessage + "\"}]}")
			response, _ := disconnectPacket.ToBytes()
			conn.Write(response)
			conn.Close()
			fmt.Println("Connection closed by server (LoginStartPacket)")
			break
		}

		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
	}
}

func main() {
	config := config.LoadConfig("config.json")

	ln, err := net.Listen("tcp", ":25566")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		fmt.Println("~~~~~~~~~\nClient connected\n~~~~~~~~~")
		go handleConnection(conn, &config)
	}
}
