package common

import (
	"fmt"
)

type AudienceError string

func (e AudienceError) Error() string {
	return fmt.Sprintf("incorrect audience, audience found: %v", string(e))
}

type TypeAssertError struct {
	Srv   string
	Value string
}

func (e TypeAssertError) Error() string {
	return fmt.Sprintf("type asser error for service: %v on value: %v", e.Srv, e.Value)
}

type CtxValueKeyMissingError struct {
	CtxKey string
}

func (e CtxValueKeyMissingError) Error() string {
	return fmt.Sprintf("ctx value key missing service: %v", e.CtxKey)
}
