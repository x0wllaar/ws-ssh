package client

import "io"

type connectedReadWriter struct {
	r io.Reader
	w io.Writer
}

func (rw *connectedReadWriter) Read(p []byte) (n int, err error) {
	return rw.r.Read(p)
}

func (rw *connectedReadWriter) Write(p []byte) (n int, err error) {
	return rw.w.Write(p)
}

func newConnectedReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	rw := connectedReadWriter{r, w}
	return &rw
}
