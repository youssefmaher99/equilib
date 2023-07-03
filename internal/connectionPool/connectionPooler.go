package connectionPool

import "net"

type ConnectionPooler interface {
	Get(host string) (net.Conn, error)
	Release(conn net.Conn) error
	Discard(conn net.Conn) error
	Clear() error
}
