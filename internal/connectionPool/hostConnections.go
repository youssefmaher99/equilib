package connectionPool

import (
	"fmt"
	"net"
)

type hostConnections struct {
	connections []net.Conn
}

func (h *hostConnections) Get() (net.Conn, error) {
	if len(h.connections) > 0 {
		conn := h.connections[len(h.connections)-1]
		h.connections = h.connections[:len(h.connections)-1]
		return conn, nil
	}
	return nil, fmt.Errorf("all connections are in use")
}

func (h *hostConnections) Put(conn net.Conn) error {
	h.connections = append(h.connections, conn)
	return nil
}
