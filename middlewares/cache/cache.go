package cache

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type Cache interface {
	Get(key string) (value interface{}, ok bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}

func key(req *http.Request) string {
	return req.Method + " " + req.URL.String()
}

func saveResponse(cache Cache, req *http.Request, res *http.Response, ttl time.Duration) {
	resCopy := *res
	resCopy.Body = ioutil.NopCloser(res.Body)
	respBytes, err := httputil.DumpResponse(&resCopy, true)
	if err == nil {
		cache.Set(key(req), respBytes, ttl)
	}
}

func getResponse(cache Cache, req *http.Request) (*http.Response, error) {
	cached, ok := cache.Get(key(req))
	if !ok || cached == nil {
		return nil, nil
	}
	cachedResponse, ok := cached.([]byte)
	if !ok {
		return nil, errors.New("invalid cached response type")
	}
	buf := bytes.NewBuffer(cachedResponse)
	res, err := http.ReadResponse(bufio.NewReader(buf), req)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't serve response")
	}
	return res, nil
}

func removeResponse(cache Cache, req *http.Request) {
	cache.Delete(key(req))
}
