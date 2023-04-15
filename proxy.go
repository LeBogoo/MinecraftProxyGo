package main

import (
	"bufio"
	"fmt"
	"net"

	"minecraftproxy/packetUtils"
)

func handleConnection(conn net.Conn) {

	state := 0

	reader := bufio.NewReader(conn)

	defer conn.Close()

	for {
		packet, err := packetUtils.ParsePacket(reader)
		fmt.Println("*************************")
		fmt.Println("PacketLength:", packet.GetPacketLength())
		fmt.Println("PacketId:", packet.GetPacketId())
		fmt.Println("State:", state)

		if (state == 0) && (packet.GetPacketId() == 0x00) {
			handshakePacket, _ := packetUtils.ParseHandshakePacket(packet)
			state = handshakePacket.NextState
		}

		if (state == 1) && (packet.GetPacketId() == 0x00) {
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

		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
	}
}

func main() {
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
		go handleConnection(conn)
	}
}
