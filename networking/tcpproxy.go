package networking

import (
	"io"
	"net"
)

func StartProxying(client *net.Conn, server *net.Conn) {
	defer (*client).Close()
	defer (*server).Close()

	go io.Copy(*client, *server)
	io.Copy(*server, *client)
}
