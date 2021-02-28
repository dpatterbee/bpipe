package bpipe

import (
	"bytes"
	"io"
	"sync"
)

// Bpipe is a bytes.Buffer with a sync.Cond to allow for channel-like behaviour.
type Bpipe struct {
	writeChan       chan []byte
	readRequestChan chan int
	readChan        chan []byte
	closed          chan struct{}
}

type bpipeReader struct {
	// fields
	bpipe *Bpipe
}

type bpipeWriter struct {
	// fields
	bpipe *Bpipe
}

func New() (bpipeReader, bpipeWriter) {

	b := &Bpipe{}

	go piper(&b)

	return bpipeReader{bpipe: &b}, bpipeWriter{bpipe: &b}
}

func piper(bpipe *Bpipe) {
	var buf bytes.Buffer
	var reqs []int

	for !bpipe.closed {
		if len(reqs) > 0 {
			if buf.Len() >= reqs[0] {
				s := make([]byte, reqs[0])
				buf.Read(s)
				reqs = reqs[1:]
			}
		}

		select {
		case p := <-bpipe.WriteChan:
			buf.Write(p)
		case p := <-bpipe.readRequestChan:

			reqs = append(reqs, p)
		case <-bpipe.closed:
			return
		}
	}
}

// Read waits for either b to be closed or to contain enough data to fill p then reads n bytes into p and signals another waiting reader.
// The read will wait indefinitely if no further writes are made and the bpipe is never closed.
func (b *Bpipe) Read(p []byte) (n int, err error) {
	b.c.L.Lock()
	defer b.c.L.Unlock()

	defer b.c.Signal()

	for b.buf.Len() < len(p) && !b.pipeClosed {

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
