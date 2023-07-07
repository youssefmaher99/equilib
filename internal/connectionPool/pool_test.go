package connectionPool

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// All tests requrie server generator

func TestGet(t *testing.T) {
	pool := New(1)
	hosts := []string{"127.0.0.1:5000", "127.0.0.1:5001", "127.0.0.1:5002"}
	err := pool.Populate(hosts)
	assert.Nil(t, err)
	assert.Equal(t, len(hosts)*pool.maxIdleConnsPerHost, pool.totalAvailableWarmConnections)

	conn, err := pool.Get(hosts[0])
	assert.Nil(t, err)
	assert.Equal(t, len(hosts)*pool.maxIdleConnsPerHost-1, pool.totalAvailableWarmConnections)
	var connFromPool net.Conn
	for _, hostConnections := range pool.connections {
		for _, connection := range hostConnections.connections {
			if connection == conn {
				assert.Equal(t, conn, connFromPool)
				break
			}
		}
	}

	// closing connection and try to get it again will cause it to be discarded and a new connection is dialed
	conn.Close()
	err = pool.Release(conn, hosts[0])
	assert.Nil(t, err)

	conn2, err := pool.Get(hosts[0])
	assert.Nil(t, err)
	assert.NotEqual(t, conn, conn2)
}

func TestRelease(t *testing.T) {
	pool := New(1)
	hosts := []string{"127.0.0.1:5000", "127.0.0.1:5001", "127.0.0.1:5002"}
	err := pool.Populate(hosts)
	assert.Nil(t, err)
	assert.Equal(t, len(hosts)*pool.maxIdleConnsPerHost, pool.totalAvailableWarmConnections)

	assert.Equal(t, len(pool.connections[hosts[0]].connections), 1)
	conn, err := pool.Get(hosts[0])
	assert.Nil(t, err)
	assert.Equal(t, len(pool.connections[hosts[0]].connections), 0)

	// closing connection and try to get it again will cause it to be discarded and a new connection is dialed
	err = pool.Release(conn, hosts[0])
	assert.Nil(t, err)
	assert.Equal(t, len(pool.connections[hosts[0]].connections), 1)
}

func TestDiscard(t *testing.T) {
	pool := New(5)
	hosts := []string{"127.0.0.1:5000", "127.0.0.1:5001", "127.0.0.1:5002"}
	err := pool.Populate(hosts)
	assert.Nil(t, err)

	type connWithHost struct {
		conn net.Conn
		host string
	}

	connections := []connWithHost{}
	for host, hostConnections := range pool.connections {
		for _, connection := range hostConnections.connections {
			connections = append(connections, connWithHost{conn: connection, host: host})
		}
	}

	for _, connection := range connections {
		pool.Discard(connection.conn, connection.host)
	}

	for _, host := range hosts {
		assert.Equal(t, len(pool.connections[host].connections), 0)
	}
}

func TestClear(t *testing.T) {
	pool := New(5)
	hosts := []string{"127.0.0.1:5000", "127.0.0.1:5001", "127.0.0.1:5002"}
	err := pool.Populate(hosts)
	assert.Nil(t, err)

	pool.Clear()
	assert.Equal(t, len(pool.connections), 0)
}
