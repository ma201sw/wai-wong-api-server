package golib

import (
	"context"
	"net/http"

	"go-wai-wong/common"
)

func Inject(as Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithGoLib(r.Context(), as)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

const ctxKey = "4486309c-9d1d-4442-bfdd-205f6eec426d"

func WithGoLib(ctx context.Context, service Service) context.Context {
	return context.WithValue(ctx, ctxKey, service)
}

func FromContextAs(ctx context.Context, out interface{}) error {
	ctxValueKey := ctx.Value(ctxKey)

	if ctxValueKey == nil {
		return common.CtxValueKeyMissingError{CtxKey: ctxKey}
	}

	srv, ok := ctxValueKey.(Service)
	if !ok {
		return common.TypeAssertError{Srv: "golib", Value: "ctxValueKey"}
	}

	outTypeAssert, outOk := out.(*Service)

	if !outOk {
		return common.TypeAssertError{Srv: "golib", Value: "out"}
	}

	*outTypeAssert = srv

	return nil
}
