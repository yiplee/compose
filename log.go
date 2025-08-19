package compose

import (
	"bytes"
	"io"
)

type PrefixWriter struct {
	Out    io.Writer
	Prefix []byte
	buf    bytes.Buffer
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	// Reuse buffer to avoid allocations
	w.buf.Reset()
	w.buf.Write(w.Prefix)
	w.buf.Write(p)
	
	// Write the combined data
	_, err = w.Out.Write(w.buf.Bytes())
	if err != nil {
		return 0, err
	}
	
	return len(p), nil
}

// NewPrefixWriter creates a new PrefixWriter with optimized memory usage
func NewPrefixWriter(out io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{
		Out:    out,
		Prefix: []byte(prefix),
		buf:    bytes.Buffer{},
	}
}
