package cassinirpc

import (
	"context"

	"github.com/bufbuild/connect-go"
)

type Service func(context.Context, connect.AnyRequest) (connect.AnyResponse, error)

func (s Service) AsInterceptor() connect.UnaryInterceptorFunc {
	fn := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(
			func(
				ctx context.Context,
				req connect.AnyRequest,
			) (connect.AnyResponse, error) {
				return s(ctx, req)
			})
	}
	return connect.UnaryInterceptorFunc(fn)
}
