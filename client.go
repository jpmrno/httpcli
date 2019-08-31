package httpcli

import "net/http"

type Client struct {
	*http.Client
	Middleware
}

func New() *Client {
	return &Client{
		Client:     &http.Client{},
		Middleware: newLayer(),
	}
}

func (cli *Client) Do(req *http.Request) (*http.Response, error) {
	mainHandler := func(c *Context) {
		res, err := cli.Client.Do(req)
		if err != nil {
			c.Abort(err)
		} else {
			c.Response = res
		}
	}
	chain := append(cli.chain(), mainHandler)
	c := NewContext(chain, req)
	c.Next()
	return c.Response, c.Error()
}
