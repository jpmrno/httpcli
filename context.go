package httpcli

import (
	"github.com/pkg/errors"
	"net/http"
)

const notAbortedIndex = -1

type Context struct {
	//
	Request  *http.Request
	Response *http.Response
	error    error
	// Abort
	abortedIndex int
	// Handlers
	handlers     HandlersChain
	currentIndex int
}

func NewContext(chain HandlersChain, req *http.Request) *Context {
	return &Context{
		Request:      req,
		Response:     nil,
		error:        nil,
		abortedIndex: notAbortedIndex,
		handlers:     chain,
		currentIndex: 0,
	}
}

func (c *Context) Error() error {
	return c.error
}

func (c *Context) Next() {
	if c.currentIndex == c.abortedIndex {
		panic(errors.New("aborted context"))
	}

	currentReq := new(http.Request)
	*currentReq = *c.Request
	currentIndex := c.currentIndex

	c.abortedIndex = notAbortedIndex
	c.error = nil
	c.Response = nil
	for c.currentIndex < len(c.handlers) && c.error == nil {
		c.handlers[c.currentIndex](c)
		c.currentIndex++
	}

	c.Request = currentReq
	c.currentIndex = currentIndex
}

func (c *Context) Stop(res *http.Response) {
	if res == nil {
		panic("res is nil")
	}
	c.abortedIndex = c.currentIndex
	c.Response = res
}

func (c *Context) Abort(err error) {
	if err == nil {
		panic("err is nil")
	}
	c.abortedIndex = c.currentIndex
	c.error = err
}
