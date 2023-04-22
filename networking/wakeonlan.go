package networking

import (
	"net"
)

func WakeOnLAN(macAddress string) error {
	macBytes, err := net.ParseMAC(macAddress)
	if err != nil {
		return err
	}

	// Construct the WoL packet
	packet := []byte{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		packet = append(packet, macBytes...)
	}

	// Send the packet to the network broadcast address
	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(packet)
	if err != nil {
		return err
	}
	return nil
}
