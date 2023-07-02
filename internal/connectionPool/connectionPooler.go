package connectionPool

import "net"

type connectionPooler interface {
	Get() (net.Conn, error)
	Release(conn net.Conn) error
	Discard(conn net.Conn) error
	Clear() error
}
