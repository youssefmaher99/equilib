package connectionPool

import (
	"fmt"
	"net"
	"sync"
)

type pool struct {
	// pool size
	mu                            *sync.RWMutex
	connections                   map[string]*hostConnections
	maxIdleConnsPerHost           int
	totalAvailableWarmConnections int
}

func New() *pool {
	// TODO : change later
	return &pool{maxIdleConnsPerHost: 3, mu: &sync.RWMutex{}}
}

func (p *pool) Get(server string) (net.Conn, error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	hostConnection, ok := p.connections[server]
	if !ok {
		return p.OpenNewTcpConnection(server)
	}

	conn, err := hostConnection.Get()
	if err != nil {
		return p.OpenNewTcpConnection(server)
	}
	p.totalAvailableWarmConnections--
	return conn, nil
}
func (p *pool) Release(conn net.Conn) error {
	return nil
}
func (p *pool) Discard(conn net.Conn) error {
	return nil
}
func (p *pool) Clear() error {
	return nil
}

func (p *pool) Populate(servers []string) error {
	p.connections = make(map[string]*hostConnections, len(servers)*p.maxIdleConnsPerHost)
	// TODO : allow method to accept pools option that define number of connections per host
	// and better delegate Populating to New()

	// create host_connections struct & default number of connections for each server
	for _, server := range servers {
		h_c := hostConnections{make([]net.Conn, 0, p.maxIdleConnsPerHost)}
		for i := 0; i < p.maxIdleConnsPerHost; i++ {
			conn, err := net.Dial("tcp", server)
			if err != nil {
				panic(err)
			}
			h_c.connections = append(h_c.connections, conn)
		}
		p.connections[server] = &h_c
		p.totalAvailableWarmConnections++
	}
	fmt.Println("----- Pool populated succesfuly -----")
	return nil
}
func (p *pool) OpenNewTcpConnection(server string) (net.Conn, error) {
	fmt.Println("********************** New tcp connections **********************")
	// TODO : custom error if server is down
	return net.Dial("tcp", server)
}
