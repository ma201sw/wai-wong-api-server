package golib

import "io"

type IOClient interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

func (c goLibImpl) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
