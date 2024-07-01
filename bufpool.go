package slogtfmt

import "sync"

const (
	initialBufferSize = 1024
	maxBufferSize     = 16 << 10 // 16384
)

// bufPool is a sync.Pool that provides a pool of byte slices to reduce memory allocations.
// The pool will allocate new byte slices of size initialBufferSize when the pool is empty.
// Byte slices larger than maxBufferSize will not be returned to the pool to reduce peak memory usage.
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, initialBufferSize)
		return &b
	},
}

// allocBuf returns a new byte slice from the bufPool. The byte slice will have an initial capacity of initialBufferSize.
func allocBuf() *[]byte {
	return bufPool.Get().(*[]byte)
}

// freeBuf returns a byte slice back to the bufPool. Byte slices larger than maxBufferSize
// will not be returned to the pool to reduce peak memory usage.
func freeBuf(b *[]byte) {
	// To reduce peak allocation, return only smaller buffers to the pool.
	if cap(*b) > maxBufferSize {
		return
	}
	*b = (*b)[:0]
	bufPool.Put(b)
}
