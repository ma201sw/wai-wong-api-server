package tokenhelper

import (
	"context"
	"net/http"

	"go-wai-wong/common"
)

func Inject(as Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithTokenHelper(r.Context(), as)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

const ctxKey = "07cd889c-9a76-4a28-ab05-c77dd7c6b5ba"

func WithTokenHelper(ctx context.Context, service Service) context.Context {
	return context.WithValue(ctx, ctxKey, service)
}

func FromContextAs(ctx context.Context, out interface{}) error {
	ctxValueKey := ctx.Value(ctxKey)

	if ctxValueKey == nil {
		return common.CtxValueKeyMissingError{CtxKey: ctxKey}
	}

	srv, ok := ctxValueKey.(Service)
	if !ok {
		return common.TypeAssertError{Srv: "tokenhelper", Value: "ctxValueKey"}
	}

	outTypeAssert, outOk := out.(*Service)

	if !outOk {
		return common.TypeAssertError{Srv: "tokenhelper", Value: "out"}
	}

	*outTypeAssert = srv

	return nil
}
