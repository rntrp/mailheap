package smtprecv

import "io"

type readerDecorator struct {
	delegate io.Reader
	length   uint64
	err      error
}

func (r *readerDecorator) Read(p []byte) (int, error) {
	n, err := r.delegate.Read(p)
	if err != nil {
		r.length += uint64(n)
	} else {
		r.err = err
	}
	return n, err
}
