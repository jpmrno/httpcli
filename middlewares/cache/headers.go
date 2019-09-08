package cache

import (
	"net/http"
	"time"
)

const ( // Headers
	headerDate          = "Date"
	headerCacheControl  = "Cache-Control"
	headerExpires       = "Expires"
	headerLastModified  = "Last-Modified"
	headerAuthorization = "Authorization"
)

func parseExpires(expires string) (time.Duration, bool) {
	if expires == "" {
		return 0, false
	}
	exp, err := http.ParseTime(expires)
	if err != nil {
		return 0, false
	}
	ttl := time.Until(exp)
	if ttl <= 0 {
		return 0, false
	}
	return ttl, true
}
