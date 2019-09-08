package cache

import (
	"github.com/jpmrno/httpcli"
	"github.com/jpmrno/httpcli/slices"
	"net/http"
	"time"
)

const ( // Other
	MaxTTL     = time.Minute * 5
	DefaultTTL = time.Minute * 1
)

var (
	cacheableMethods     = []string{http.MethodGet, http.MethodHead, http.MethodOptions}
	cacheableStatusCodes = []int{http.StatusOK, http.StatusNonAuthoritativeInfo, http.StatusNoContent, http.StatusPartialContent,
		http.StatusMultipleChoices, http.StatusMovedPermanently, http.StatusNotFound, http.StatusMethodNotAllowed, http.StatusGone,
		http.StatusRequestURITooLong, http.StatusNotImplemented}
)

func Enable(cache Cache) httpcli.HandlerFunc {
	return func(ctx *httpcli.Context) {
		minTTL, maxAge, cacheable := reqIsCacheable(ctx.Request)
		if !cacheable {
			removeResponse(cache, ctx.Request)
			ctx.Next()
			return
		}
		res, _ := getResponse(cache, ctx.Request)
		if res != nil && getTTL(res) >= minTTL && getAge(res) <= maxAge {
			ctx.Stop(res)
			return
		}
		removeResponse(cache, ctx.Request)
		before := time.Now()
		ctx.Next()
		elapsedTime := time.Since(before)
		if ctx.Error() != nil {
			return
		}
		if ttl, ok := resIsCacheable(ctx.Response); ok {
			if ttl > MaxTTL {
				ttl = MaxTTL
			}
			updateHeaders(res, elapsedTime, ttl)
			saveResponse(cache, ctx.Request, ctx.Response, ttl)
		}
		return
	}
}

func reqIsCacheable(req *http.Request) (minTTL time.Duration, maxAge time.Duration, cacheable bool) {
	// The Authorization header does not appear in the request and the request method is cacheable
	if !slices.ContainsString(cacheableMethods, req.Method) || req.Header.Get(headerAuthorization) != "" {
		return
	}
	cacheControl := parseCacheControl(req.Header.Get(headerCacheControl))
	// The "no-store"/"no-cache" directives does not appear in request headers
	if cacheControl.Contains(DirectiveNoStore) || cacheControl.Contains(DirectiveNoCache) {
		return
	}
	// TODO: Doc
	if secs, ok := cacheControl.Seconds(DirectiveMaxAge); ok {
		if secs <= 0 {
			return
		}
		maxAge = secs
	} else {
		maxAge = MaxTTL
	}
	// TODO: Doc
	if secs, ok := cacheControl.Seconds(DirectiveMinFresh); ok {
		minTTL = secs
	}
	cacheable = true
	return
}

func resIsCacheable(res *http.Response) (time.Duration, bool) {
	cacheControl := parseCacheControl(res.Header.Get(headerCacheControl))
	// The "no-store"/"no-cache" directives does not appear in response headers
	if cacheControl.Contains(DirectiveNoStore) || cacheControl.Contains(DirectiveNoCache) {
		return 0, false
	}
	// Contains a max-age response directive
	if ttl, ok := cacheControl.Seconds(DirectiveMaxAge); ok {
		if ttl <= 0 {
			return 0, false
		}
		return ttl, true
	}
	// Contains an Expires header
	if ttl, ok := parseExpires(res.Header.Get(headerExpires)); ok {
		return ttl, true
	}
	// Heuristic freshness lifetime calculation
	// Has a status code that is defined as cacheable by default or contains a public response directive
	if slices.ContainsInt(cacheableStatusCodes, res.StatusCode) || cacheControl.Contains(DirectivePublic) {
		if lastModified := res.Header.Get(headerLastModified); lastModified != "" {
			if lastModifiedDate, err := http.ParseTime(lastModified); err == nil {
				return time.Since(lastModifiedDate) / 10.0, true
			}
		}
		return DefaultTTL, true
	}
	// Contains a Cache Control Extension
	// ?
	return 0, false
}
