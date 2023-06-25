package ringBuffer

import (
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

func (r *RingBuffer) Populate(servers []string) {
	r.servers = append(r.servers, servers...)
}
