package jsonprovider

import (
	"context"
	"net/http"

	"go-wai-wong/common"
)

func Inject(as Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithJSONProvider(r.Context(), as)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

const ctxKey = "d539f36d-b291-4b01-bf65-6169a6e179eb"

func WithJSONProvider(ctx context.Context, service Service) context.Context {
	return context.WithValue(ctx, ctxKey, service)
}

func FromContextAs(ctx context.Context, out interface{}) error {
	ctxValueKey := ctx.Value(ctxKey)

	if ctxValueKey == nil {
		return common.CtxValueKeyMissingError{CtxKey: ctxKey}
	}

	srv, ok := ctxValueKey.(Service)
	if !ok {
		return common.TypeAssertError{Srv: "jsonprovider", Value: "ctxValueKey"}
	}

	outTypeAssert, outOk := out.(*Service)

	if !outOk {
		return common.TypeAssertError{Srv: "jsonprovider", Value: "out"}
	}

	*outTypeAssert = srv

	return nil
}
