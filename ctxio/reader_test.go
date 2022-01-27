package ctxio

import (
	"bytes"
	"context"
	"io"
	"testing"
)

// TestReadAll ensures we have not broken io.Reader contract.
func TestReadAll(t *testing.T) {
	src := []byte(`foo-bar-baz`)
	r := NewReader(context.Background(), bytes.NewBuffer(src))
	dst, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read from buffer: %v", err)
	}
	if !bytes.Equal(src, dst) {
		t.Fatalf("expected:\n\t%v\nreceived:\n\t%v", src, dst)
	}
}

// TestAbortReadOnCanceledContext ensures Read call returns error
// once context is no longer valid.
func TestAbortReadOnCanceledContext(t *testing.T) {
	src := []byte(`foo-bar-baz`)
	dst1 := make([]byte, len(src)/2)
	dst2 := make([]byte, len(src)/2)

	ctx, cancel := context.WithCancel(context.Background())
	r := NewReader(ctx, bytes.NewBuffer(src))

	n1, err := r.Read(dst1)
	if err != nil {
		t.Fatalf("read unexpectedly failed: %v", err)
	}
	if len(dst1) != n1 {
		t.Fatalf("uneven read length: expected %v, actual %v", len(dst1), n1)
	}

	cancel()

	n2, err := r.Read(dst2)
	if err == nil {
		t.Fatalf("read unexpectedly succeeded (%v bytes) after context was cancelled", n2)
	}
	if n2 != 0 {
		t.Fatalf("n2 should be empty but has a length %v", n2)
	}
}
