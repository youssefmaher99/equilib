package connectionPool

import "net"

type ConnectionPooler interface {
	Get(host string) (net.Conn, error)
	Release(conn net.Conn, host string) error
	Discard(conn net.Conn, host string)
	Clear() error
}
