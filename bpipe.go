package bpipe

import (
	"bytes"
)

// Bpipe is a bytes.Buffer with a sync.Cond to allow for channel-like behaviour.
type Bpipe struct {
	writeChan       chan []byte
	readRequestChan chan int
	readChan        chan []byte
	closed          chan struct{}
}

type BpipeReader struct {
	// fields
	bpipe *Bpipe
}

type BpipeWriter struct {
	// fields
	bpipe *Bpipe
}

func New() (BpipeReader, BpipeWriter) {

	b := Bpipe{}

	go piper(&b)

	return BpipeReader{bpipe: &b}, BpipeWriter{bpipe: &b}
}

func piper(bpipe *Bpipe) {
	var buf bytes.Buffer
	var reqs []int

	for {
		if len(reqs) > 0 {
			if buf.Len() >= reqs[0] {
				s := make([]byte, reqs[0])
				buf.Read(s)
				reqs = reqs[1:]
			}
		}

		select {
		case p := <-bpipe.writeChan:
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
func (b *BpipeReader) Read(p []byte) (n int, err error) {
	b.bpipe.readRequestChan <- len(p)

	g := <-b.bpipe.readChan

	n = copy(p, g)
	return n, nil
}

// Write writes n bytes from p into the buffer then signals any waiting reader.
func (b *BpipeWriter) Write(p []byte) (n int, err error) {

	b.bpipe.writeChan <- p
	return len(p), nil

}

// Close closes the Bpipe and signals a waiting reader
func (b *Bpipe) Close() error {
	b.closed <- struct{}{}
	return nil
}
