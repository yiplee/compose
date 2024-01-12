package compose

import (
	"fmt"
	"io"
	"strings"
)

type PrefixWriter struct {
	Out    io.Writer
	Prefix string
	buffer strings.Builder
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if b == '\n' {
			line := w.buffer.String()
			lineWithPrefix := w.Prefix + line
			_, err := fmt.Fprint(w.Out, lineWithPrefix+"\n")
			if err != nil {
				return n, err
			}
			w.buffer.Reset()
		} else {
			w.buffer.WriteByte(b)
		}
		n++
	}
	return n, nil
}
