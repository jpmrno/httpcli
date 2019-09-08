package cache

import (
	"strings"
	"time"
)

const ( // Directives
	// REQ The client is unwilling to accept a response whose age is greater than the specified number of seconds.
	// RES The response is to be considered stale after its age is greater than the specified number of seconds.
	DirectiveMaxAge = "max-age"
	// REQ The client is willing to accept a response whose freshness lifetime is no less than its current age plus the specified time in seconds.
	// The client wants a response that will still be fresh for at least the specified number of seconds.
	DirectiveMinFresh = "min-fresh"
	// REQ/RES The cache MUST NOT store any part of either this request or any response to it.
	DirectiveNoStore = "no-store"
	// REQ/RES The cache MUST NOT use a stored response to satisfy the request.
	DirectiveNoCache = "no-cache"
	// RES The cache MAY store the response, even if the response would normally be non-cacheable.
	DirectivePublic = "public"
	// REQ The client only wishes to obtain a stored response. The cache SHOULD respond using a stored response or a 504 (Gateway Timeout)
	//DirectiveOnlyIfCached = "only-if-cached" // TODO: Add support
	// REQ/RES The cache MUST NOT transform the payload // TODO: No hace falta probablemente
	//directiveNoTransform = "no-transform"
)

type cacheControl map[string]string

func (cc cacheControl) Contains(key string) bool {
	_, ok := cc[key]
	return ok
}

func (cc cacheControl) Seconds(key string) (time.Duration, bool) {
	val := cc[key]
	if val == "" {
		return 0, false
	}
	seconds, err := time.ParseDuration(val + "s")
	if err == nil {
		return 0, false
	}
	return seconds, true
}

func parseCacheControl(ccHeader string) cacheControl {
	cc := cacheControl{}
	for _, directive := range strings.Split(ccHeader, ",") {
		if directive != "" {
			keyval := strings.SplitN(directive, "=", 1)
			key := strings.TrimSpace(keyval[0])
			if len(keyval) > 1 {
				cc[key] = strings.TrimSpace(keyval[1])
			} else {
				cc[key] = ""
			}
		}
	}
	return cc
}
