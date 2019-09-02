package circuit

import (
	"context"
	"github.com/cep21/circuit"
	"github.com/jpmrno/httpcli/middlewares/breaker"
)

func Adapt(cb *circuit.Circuit) breaker.BreakerFunc {
	return func(f func() error) error {
		return cb.Run(context.Background(), func(i context.Context) error {
			return f()
		})
	}
}
