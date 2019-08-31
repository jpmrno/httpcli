package url

import (
	"github.com/jpmrno/httpcli"
	"net/url"
	"path"
	"strings"
)

var pathParamPrefix = ":"

func RawURL(raw string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		parsed, err := url.Parse(raw)
		if err != nil {
			ctx.Abort(err)
			return
		}
		ctx.Request.URL = parsed
		ctx.Next()
	}
}

func BaseURL(base url.URL) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		ctx.Request.URL.Scheme = base.Scheme
		ctx.Request.URL.Host = base.Host
		ctx.Request.URL.Path = path.Join(base.Path, ctx.Request.URL.Path)
		query := base.Query()
		for k, vs := range ctx.Request.URL.Query() {
			for _, v := range vs {
				query.Add(k, v)
			}
		}
		ctx.Request.URL.RawQuery = query.Encode()
		ctx.Next()
	}
}

func SchemeAndHost(scheme, host string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		ctx.Request.URL.Scheme = scheme
		ctx.Request.URL.Host = host
		ctx.Next()
	}
}

func Path(pth string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		ctx.Request.URL.Path = pth
		ctx.Next()
	}
}

func PathPrefix(pth string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		ctx.Request.URL.Path = path.Join(pth, ctx.Request.URL.Path)
		ctx.Next()
	}
}

func PathSuffix(pth string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		ctx.Request.URL.Path = path.Join(ctx.Request.URL.Path, pth)
		ctx.Next()
	}
}

func PathParam(key, value string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		ctx.Request.URL.Path = replacePathParam(ctx.Request.URL.Path, key, value)
		ctx.Next()
	}
}

func PathParams(params map[string]string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		for k, v := range params {
			ctx.Request.URL.Path = replacePathParam(ctx.Request.URL.Path, k, v)
		}
		ctx.Next()
	}
}

func SetQueryParam(key, value string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		query := ctx.Request.URL.Query()
		query.Set(key, value)
		ctx.Request.URL.RawQuery = query.Encode()
		ctx.Next()
	}
}

func SetQueryParams(params map[string]string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		query := ctx.Request.URL.Query()
		for k, v := range params {
			query.Set(k, v)
		}
		ctx.Request.URL.RawQuery = query.Encode()
		ctx.Next()
	}
}

func AddQueryParam(key, value string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		query := ctx.Request.URL.Query()
		query.Add(key, value)
		ctx.Request.URL.RawQuery = query.Encode()
		ctx.Next()
	}
}

func DelQueryParam(key string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		query := ctx.Request.URL.Query()
		query.Del(key)
		ctx.Request.URL.RawQuery = query.Encode()
		ctx.Next()
	}
}

func DelQueryParams(keys ...string) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		query := ctx.Request.URL.Query()
		for _, k := range keys {
			query.Del(k)
		}
		ctx.Request.URL.RawQuery = query.Encode()
		ctx.Next()
	}
}

func DelAllQueryParams() httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		ctx.Request.URL.RawQuery = ""
		ctx.Next()
	}
}

func replacePathParam(path, k, v string) string {
	return strings.Replace(path, pathParamPrefix+k, v, 1)
}
