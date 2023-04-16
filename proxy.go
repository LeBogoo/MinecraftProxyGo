package main

import (
	"bufio"
	"fmt"
	"net"

	"minecraftproxy/config"
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
			response := packetUtils.CreateStatusResponsePacket(
				packetUtils.StatusResponse{
					Version: packetUtils.Version{
						Name:     "1.19.3",
						Protocol: 761,
					},
					Players: packetUtils.Players{
						Max:    100,
						Online: 0,
						Sample: []packetUtils.Sample{},
					},
					Description: packetUtils.Description{
						Text: "This is not a Minecraft server, it is just a program written in Go.",
					},
					Favicon:            "",
					EnforcesSecureChat: false,
				},
			)

			bytes, err := response.ToBytes()

			if err != nil {
				fmt.Println("Error converting packet to bytes:", err)
				break
			}

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

			disconnectPacket := packetUtils.CreateLoginDisconnectPacket("{\"text\":\"Starting...\n\",\"bold\":true,\"color\":\"#00ff00\",\"extra\":[{\"color\":\"white\",\"bold\":false,\"text\":\"The server was offline, but is now starting. Please wait a few econds and try connecting again!\"}]}")
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
