package compose

import "io"

type PrefixWriter struct {
	Out    io.Writer
	Prefix string
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	if _, err := w.Out.Write([]byte(w.Prefix)); err != nil {
		return 0, err
	}

	return w.Out.Write(p)
}
