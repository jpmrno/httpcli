package gobreaker

import (
	"github.com/jpmrno/httpcli/middlewares/breaker"
	"github.com/sony/gobreaker"
)

func Adapt(cb *gobreaker.CircuitBreaker) breaker.BreakerFunc {
	return func(f func() error) error {
		_, err := cb.Execute(func() (i interface{}, e error) {
			return nil, f()
		})
		return err
	}
}
