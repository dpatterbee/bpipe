package bpipe

import (
	"testing"
)

func TestAll(t *testing.T) {
	pr, pw := New()

	a := []byte{1, 2, 3, 4, 5, 6, 7, 7, 8}

	n, _ := pw.Write(a)

	if len(a) != n {
		t.Errorf("Wanted %v, got %v", len(a), n)
	}

	bra := make([]byte, len(a))

	pr.Read(bra)

	t.Log(bra)
}

func TestWrite(t *testing.T) {
	_, pw := New()

	a := []byte{1, 2, 3, 4}
	want := len(a)

	n, _ := pw.Write(a)

	if want != n {
		t.Errorf("Wanted %v, got %v", want, n)
	}

}

func TestRead(t *testing.T) {
	pr, pw := New()

	a := []byte("This is some data")
	writeLen := len(a)

	pw.Write(a)

	n := 8
	readNSlice := make([]byte, n)

	returnedN, _ := pr.Read(readNSlice)

	if returnedN != n {
		t.Errorf("Wanted %v, got %v", n, returnedN)
	}

	n = 100
	readNSlice2 := make([]byte, n)

	returnedN2, _ := pr.Read(readNSlice2)

	remaining := writeLen - returnedN

	if returnedN2 != remaining {
		t.Errorf("Wanted %v, got %v", remaining, returnedN2)
	}
}

func TestReadRest(t *testing.T) {

	pr, pw := New()

	a := []byte("My name is chuck")

	writeLen, _ := pw.Write(a)

	re := make([]byte, 100)

	pr.Close()

	pr.Read(re)

}
