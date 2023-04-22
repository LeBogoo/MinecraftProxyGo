package main

import (
	"bufio"
	"fmt"
	"net"

	"minecraftproxy/config"
	"minecraftproxy/networking"
	"minecraftproxy/packetUtils"
)

func handleConnection(conn net.Conn, globalState *config.State) {
	state := 0
	config := globalState.Config
	playing := false
	var username, uuid string

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

			if err != nil && !globalState.Starting {
				response = packetUtils.CreateStatusResponsePacket(config.StatusResponses.Offline)
			} else if err != nil && globalState.Starting {
				response = packetUtils.CreateStatusResponsePacket(config.StatusResponses.Starting)
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
			username = loginStartPacket.Name
			uuid = loginStartPacket.PlayerUUID

			_, err := networking.StatusPingServer(&config.Server)

			var response []byte

			if err != nil && !globalState.Starting { // server is offline and not starting. tell the player that the server is now starting
				globalState.Starting = true

				fmt.Println("Waking up server...")
				err := networking.WakeOnLAN(config.WakeOnLan.Mac)
				if err != nil {
					fmt.Println("Error waking up server:", err)
					globalState.Config.StatusResponses.Offline.Description.Text = "Error starting up server. Please contact an administrator."
					globalState.Config.StatusResponses.Offline.Players.Max = -1

					config.DisconnectMessages.NowStarting.Text = "Error starting up server. Please contact an administrator."
					config.DisconnectMessages.NowStarting.Color = "red"
					globalState.Starting = false
				}

				disconnectPacket := packetUtils.CreateLoginDisconnectPacket(config.DisconnectMessages.NowStarting)
				response, _ = disconnectPacket.ToBytes()
				conn.Write(response)
			} else if err != nil && globalState.Starting { // if there is an error (offline) and it is starting tell the player
				disconnectPacket := packetUtils.CreateLoginDisconnectPacket(config.DisconnectMessages.Starting)
				response, _ = disconnectPacket.ToBytes()
				conn.Write(response)
			} else { // no error, server is online
				fmt.Println("Server is online. Start proxying...")
				globalState.Starting = false
				playing = true
				break
			}
			conn.Close()
			fmt.Println("Connection closed by server (LoginStartPacket)")
			break
		}

		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}

		if playing {
			break
		}
	}

	if playing {
		fmt.Println("Start proxying...")
		server, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}

		fmt.Println(username, uuid)
		handshakePacket := packetUtils.CreateHandshakePacket(config.Server.Protocol, config.Server.Host, config.Server.Port, 2)
		bytes, _ := handshakePacket.ToBytes()
		server.Write(bytes)

		loginStartPacket := packetUtils.CreateLoginStartPacket(username, uuid)
		bytes, _ = loginStartPacket.ToBytes()
		server.Write(bytes)

		networking.StartProxying(&conn, &server)
	}
}

func main() {
	globalState := config.State{
		Starting: false,
		Config:   config.LoadConfig("config.json"),
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", globalState.Config.ProxyPort))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()

	fmt.Printf("Proxy is now listening on port %d\n-------------------------------------\n", globalState.Config.ProxyPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		fmt.Println("~~~~~~~~~\nClient connected\n~~~~~~~~~")
		go handleConnection(conn, &globalState)
	}
}
