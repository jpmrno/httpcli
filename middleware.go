package httpcli

import (
	"sync"
)

type Middleware interface {
	Use(handlers ...HandlerFunc) Middleware
	chain() HandlersChain
}

type layer struct {
	handlersMtx sync.RWMutex
	handlers    HandlersChain
}

func newLayer() *layer {
	return &layer{}
}

func (s *layer) Use(handlers ...HandlerFunc) Middleware {
	s.handlersMtx.Lock()
	s.handlers = append(s.handlers, handlers...)
	s.handlersMtx.Unlock()
	return s
}

func (s *layer) chain() HandlersChain {
	s.handlersMtx.RLock()
	defer s.handlersMtx.RUnlock()
	return s.handlers
}
