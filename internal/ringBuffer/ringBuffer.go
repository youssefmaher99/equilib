package ringBuffer

import (
	"fmt"
	"sync/atomic"
)

type RingBuffer struct {
	size    int
	servers []string
	pointer uint32
}

func New(size int) *RingBuffer {
	return &RingBuffer{size: size, servers: make([]string, 0, size)}
}

func (r *RingBuffer) Next() string {
	index := atomic.AddUint32(&r.pointer, 1) - 1
	index = index % uint32(r.size)
	return r.servers[index]
}

// TODO : ring buffer shouldn't be the one populating equilib with the servers
func (r *RingBuffer) Populate(servers []string) {
	fmt.Println("List of loaded addresses")
	fmt.Println("------------------------")
	for i := 0; i < len(servers); i++ {
		fmt.Println(servers[i])
	}
	r.servers = append(r.servers, servers...)
}
