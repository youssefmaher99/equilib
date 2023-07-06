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
	return &pool{maxIdleConnsPerHost: 5, mu: &sync.RWMutex{}}
}

func (p *pool) Get(host string) (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	hostConnection, ok := p.connections[host]
	if !ok {
		return p.OpenNewTcpConnection(host)
	}

	conn, err := hostConnection.Get()
	if err != nil {
		return p.OpenNewTcpConnection(host)
	}

	err = p.ConnectionIsValid(conn)
	if err != nil {

		// discard old connection
		p.Discard(conn, host)

		// create new connection
		conn, err := p.OpenNewTcpConnection(host)
		if err != nil {
			panic(err)
		}

		// added to the pool
		p.connections[host].Put(conn)

		return conn, nil
	}
	p.totalAvailableWarmConnections--
	return conn, nil
}
func (p *pool) Release(conn net.Conn, host string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.connections[host].Put(conn)
	p.totalAvailableWarmConnections++
	return err
}

func (p *pool) Discard(conn net.Conn, host string) {
	fmt.Println("Connection discarded")
	p.mu.Lock()
	defer p.mu.Unlock()

	for idx, connVal := range p.connections[host].connections {
		if connVal == conn {
			p.connections[host].connections = removeSliceElement(p.connections[host].connections, idx)
			break
		}
	}
}

func (p *pool) Clear() error {
	for key := range p.connections {
		delete(p.connections, key)
	}
	return nil
}

func (p *pool) ConnectionIsValid(conn net.Conn) error {
	_, err := conn.Write([]byte("1"))
	return err
}

func (p *pool) Populate(hosts []string) error {
	p.connections = make(map[string]*hostConnections, len(hosts)*p.maxIdleConnsPerHost)
	// TODO : allow method to accept pools option that define number of connections per host
	// and better delegate Populating to New()

	// create host_connections struct & default number of connections for each host
	for _, host := range hosts {
		h_c := hostConnections{make([]net.Conn, 0, p.maxIdleConnsPerHost)}
		for i := 0; i < p.maxIdleConnsPerHost; i++ {
			conn, err := net.Dial("tcp", host)
			if err != nil {
				panic(err)
			}
			h_c.connections = append(h_c.connections, conn)
		}
		p.connections[host] = &h_c
		p.totalAvailableWarmConnections++
	}
	fmt.Println("----- Pool populated succesfuly -----")
	return nil
}
func (p *pool) OpenNewTcpConnection(host string) (net.Conn, error) {
	// fmt.Println("********************** Opening new tcp connection **********************")
	// TODO : custom error if host server is down
	return net.Dial("tcp", host)
}

func removeSliceElement(arr []net.Conn, index int) []net.Conn {
	arr[index], arr[len(arr)-1] = arr[len(arr)-1], arr[index]
	return arr
}
