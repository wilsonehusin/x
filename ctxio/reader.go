package ctxio

import (
	"context"
	"fmt"
	"io"
)

// ContextReader wraps io.Reader to be aware of context.Context,
// denying Read calls after provided context has been canceled.
type ContextReader interface {
	io.Reader
}

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

// NewReader returns an implementation of ContextReader.
func NewReader(ctx context.Context, r io.Reader) ContextReader {
	return &ctxReader{
		ctx: ctx,
		r:   r,
	}
}

func (cr *ctxReader) Read(p []byte) (n int, err error) {
	select {
	case <-cr.ctx.Done():
		return 0, fmt.Errorf("canceled read: %w", cr.ctx.Err())
	default:
		return cr.r.Read(p)
	}
}
