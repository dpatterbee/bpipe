package bpipe

import (
	"bytes"
	"io"
	"sync"
)

// Bpipe is a bytes.Buffer with a sync.Cond to allow for channel-like behaviour.
type Bpipe struct {
	buf        bytes.Buffer
	c          *sync.Cond
	pipeClosed bool
}

// New creates a new Bpipe
func New() *Bpipe {
	var l sync.Mutex

	return &Bpipe{
		buf:        bytes.Buffer{},
		c:          sync.NewCond(&l),
		pipeClosed: false,
	}
}

// Read waits for either b to be closed or to contain enough data to fill p then reads n bytes into p and signals another waiting reader.
// The read will wait indefinitely if no further writes are made and the bpipe is never closed.
func (b *Bpipe) Read(p []byte) (n int, err error) {
	b.c.L.Lock()
	defer b.c.L.Unlock()

	defer b.c.Signal()

	for b.buf.Len() >= len(p) && !b.pipeClosed {

		b.c.Wait()

	}

	n, err = b.buf.Read(p)

	return
}

// Write writes n bytes from p into the buffer then signals any waiting reader.
func (b *Bpipe) Write(p []byte) (n int, err error) {
	b.c.L.Lock()
	defer b.c.L.Unlock()
	defer b.c.Signal()

	if b.pipeClosed {
		return 0, io.ErrUnexpectedEOF
	}

	n, err = b.buf.Write(p)

	return

}

// Close closes the Bpipe and signals a waiting reader
func (b *Bpipe) Close() error {
	b.c.L.Lock()
	defer b.c.L.Unlock()

	if b.pipeClosed {
		return nil
	}

	b.pipeClosed = true

	defer b.c.Signal()

	return nil
}
