package ctxio

import (
	"context"
	"fmt"
	"io"
)

// ContextWriter wraps io.Writer to be aware of context.Context,
// denying Write calls after provided context has been canceled.
type ContextWriter interface {
	io.Writer
}

type ctxWriter struct {
	ctx context.Context
	w   io.Writer
}

// NewWriter returns an implementation of ContextWriter.
func NewWriter(ctx context.Context, w io.Writer) ContextWriter {
	return &ctxWriter{
		ctx: ctx,
		w:   w,
	}
}

func (cw *ctxWriter) Write(p []byte) (n int, err error) {
	select {
	case <-cw.ctx.Done():
		return 0, fmt.Errorf("canceled write: %w", cw.ctx.Err())
	default:
		return cw.w.Write(p)
	}
}
