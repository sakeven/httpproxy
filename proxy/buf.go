package proxy

import (
	"sync"
)

// BufferPool holds a pool of buffer
type BufferPool struct {
	pool *sync.Pool
	size int
}

// NewBufferPool creates a new buffer pool.
func NewBufferPool(size int) *BufferPool {
	return &BufferPool{pool: &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, size)
			return buf
		}},
		size: size,
	}
}

// Get gets a buffer from pool.
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

// Put puts buffer back pool.
func (bp *BufferPool) Put(buf []byte) {
	bp.pool.Put(buf)
}

var defaultBufferSize = 4096
var leakyBuf = NewBufferPool(defaultBufferSize)

// GetBuf gets a buffer from default pool.
func GetBuf() []byte {
	return leakyBuf.Get()
}

// PutBuf puts a buffer back default pool.
func PutBuf(buf []byte) {
	leakyBuf.Put(buf)
}
