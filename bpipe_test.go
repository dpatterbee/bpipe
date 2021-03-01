package bpipe

import (
	"io"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	bpipe := New()

	a := []byte{1, 2, 3, 4, 5, 6, 7, 7, 8}

	n, _ := bpipe.Write(a)

	if len(a) != n {
		t.Errorf("Wanted %v, got %v", len(a), n)
	}

	bra := make([]byte, len(a))

	bpipe.Read(bra)

	t.Log(bra)
}

func TestWrite(t *testing.T) {
	bpipe := New()

	a := []byte{1, 2, 3, 4}
	want := len(a)

	n, _ := bpipe.Write(a)

	if want != n {
		t.Errorf("Wanted %v, got %v", want, n)
	}

}

func TestRead(t *testing.T) {
	bpipe := New()

	a := []byte("This is some data")
	writeLen := len(a)

	bpipe.Write(a)

	n := 8
	readNSlice := make([]byte, n)

	returnedN, _ := bpipe.Read(readNSlice)

	if returnedN != n {
		t.Errorf("Wanted %v, got %v", n, returnedN)
	}

	remaining := writeLen - returnedN
	readNSlice2 := make([]byte, remaining)

	returnedN2, _ := bpipe.Read(readNSlice2)

	if returnedN2 != remaining {
		t.Errorf("Wanted %v, got %v", remaining, returnedN2)
	}
}

func TestReadRest(t *testing.T) {

	bpipe := New()

	a := []byte("My name is chuck")

	writeLen, _ := bpipe.Write(a)

	re := make([]byte, 100)

	bpipe.Close()

	n, _ := bpipe.Read(re)

	if n != writeLen {
		t.Errorf("Wanted %v, got %v", writeLen, n)
	}
}

func TestManyClose(t *testing.T) {
	bpipe := New()

	for i := 0; i < 10; i++ {
		bpipe.Close()
	}
}

func TestPipeClosedError(t *testing.T) {
	bpipe := New()

	bpipe.Close()

	_, err := bpipe.Write([]byte("hello"))

	want := io.ErrClosedPipe

	if err != want {
		t.Errorf("Wanted %v, got %v", want, err)
	}
}

func TestReadClosedPipe(t *testing.T) {
	bpipe := New()

	bpipe.Close()

	s := make([]byte, 10)
	_, err := bpipe.Read(s)

	want := io.EOF

	if err != want {
		t.Errorf("Wanted %v, got %v", want, err)
	}
}

func TestReadWait(t *testing.T) {
	bpipe := New()

	go func() {
		time.Sleep(2 * time.Second)
		bpipe.Write([]byte("Hello"))
	}()

	s := make([]byte, 5)

	bpipe.Read(s)

}
