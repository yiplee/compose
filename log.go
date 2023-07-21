package compose

import "io"

type PrefixWriter struct {
	Out    io.Writer
	Prefix string
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	n, err = w.Out.Write(append([]byte(w.Prefix), p...))

	if n > len(p) {
		n = len(p)
	}

	return
}
