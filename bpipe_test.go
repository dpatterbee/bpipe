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
