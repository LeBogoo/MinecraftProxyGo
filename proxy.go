package main

import (
	"bufio"
	"fmt"
	"net"

	"minecraftproxy/packetUtils"
	"minecraftproxy/packetUtils/utils"
)

func handleConnection(conn net.Conn) {

	state := 0

	reader := bufio.NewReader(conn)

	defer conn.Close()

	for {
		packet, _ := packetUtils.ParsePacket(reader)
		fmt.Println("*************************")
		fmt.Println("PacketLength:", packet.GetPacketLength())
		fmt.Println("PacketId:", packet.GetPacketId())
		fmt.Println("State:", state)

		if state == 0 && packet.GetPacketId() == 0x00 {
			fmt.Println("------------------------\nHandshakePacket detected\n-                      -")
			handshakePacket, _ := packetUtils.ParseHandshakePacket(packet)
			fmt.Println("ProtocolVersion:", handshakePacket.ProtocolVersion)
			fmt.Println("ServerAddress:", handshakePacket.ServerAddress)
			fmt.Println("ServerPort:", handshakePacket.ServerPort)
			fmt.Println("NextState:", handshakePacket.NextState)

			state = 1

			continue
		}

		if state == 1 && packet.GetPacketId() == 0x00 {
			fmt.Println("------------------------\nStatusRequestPacket detected\n-                      -")
			fmt.Println("No data in StatusRequestPacket")

			packetId := utils.ToVarInt(0x00)
			response := utils.ToString("{\"enforcesSecureChat\":false,\"description\":{\"color\":\"green\",\"text\":\"LeBogo Community Server\"},\"players\":{\"max\":20,\"online\":0},\"version\":{\"name\":\"Paper 1.19.3\",\"protocol\":761}}")
			packetLength := utils.ToVarInt(len(packetId) + len(response))

			var responsePacket = append(packetLength, packetId...)

			responsePacket = append(responsePacket, response...)

			_, err := conn.Write(responsePacket)
			if err != nil {
				fmt.Println(err)
			}

			continue
		}

		if state == 1 && packet.GetPacketId() == 0x01 {
			fmt.Println("------------------------\nPingRequestPacket detected\n-                      -")
			pingPacket, _ := packetUtils.ParsePingRequestPacket(packet)
			fmt.Println("PingPacket:", pingPacket)

			packetId := utils.ToVarInt(0x01)
			payload := utils.ToLong(pingPacket.Payload)
			packetLength := utils.ToVarInt(len(packetId) + len(payload))

			var responsePacket = append(packetLength, packetId...)
			responsePacket = append(responsePacket, payload...)

			_, err := conn.Write(responsePacket)
			if err != nil {
				fmt.Println(err)
			}

			conn.Close()

			return
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
