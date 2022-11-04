package compose

import "io"

type PrefixWriter struct {
	Out    io.Writer
	Prefix string
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	return w.Out.Write(append([]byte(w.Prefix), p...))
}
