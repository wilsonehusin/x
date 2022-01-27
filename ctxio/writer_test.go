package ctxio

import (
	"bytes"
	"context"
	"testing"
)

// TestSimpleWrite ensures we have not broken io.Writer contract.
func TestSimpleWrite(t *testing.T) {
	src := []byte(`foo-bar-baz`)
	var dst bytes.Buffer
	w := NewWriter(context.Background(), &dst)
	n, err := w.Write(src)
	if err != nil {
		t.Fatalf("failed to write to buffer: %v", err)
	}
	if len(src) != n {
		t.Fatalf("write was incomplete, expected %v, actual %v", len(src), n)
	}
	if !bytes.Equal(src, dst.Bytes()) {
		t.Fatalf("expected:\n\t%v\nreceived:\n\t%v", src, dst.Bytes())
	}
}

// TestAbortWriteOnCanceledContext ensures Write call returns error
// once context is no longer valid
func TestAbortWriteOnCanceledContext(t *testing.T) {
	src := []byte(`foo-bar-baz`)
	var dst bytes.Buffer

	ctx, cancel := context.WithCancel(context.Background())
	w := NewWriter(ctx, &dst)

	n, err := w.Write(src)
	if err != nil {
		t.Fatalf("write unexpectedly failed: %v", err)
	}
	if len(src) != n {
		t.Fatalf("write was incomplete, expected %v, actual %v", len(src), n)
	}

	cancel()

	n, err = w.Write(src)
	if err == nil {
		t.Fatalf("write unexpectedly succeeded (%v bytes) after context was canceled", n)
	}
	if n != 0 {
		t.Fatalf("n should be empty but has a length %v", n)
	}
	if !bytes.Equal(src, dst.Bytes()) {
		t.Fatalf("expected:\n\t%v\nreceived:\n\t%v", src, dst.Bytes())
	}
}
