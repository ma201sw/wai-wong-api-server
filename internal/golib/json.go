package golib

import "encoding/json"

type JSONClient interface {
	Unmarshal(data []byte, v interface{}) error
	Marshal(v interface{}) ([]byte, error)
}

func (c goLibImpl) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (c goLibImpl) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
