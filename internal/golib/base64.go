package golib

import "encoding/base64"

type Base64Client interface {
	StdEncodingDecodeString(s string) ([]byte, error)
}

func (c goLibImpl) StdEncodingDecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
