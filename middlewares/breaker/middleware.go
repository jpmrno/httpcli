package breaker

import (
	"github.com/jpmrno/httpcli"
)

type BreakerFunc func(func() error) error

func Using(f BreakerFunc) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		err := f(func() error {
			ctx.Next()
			return ctx.Error()
		})
		if err != nil {
			ctx.Abort(err)
		}
	}
}
