package cache

import (
	"net/http"
	"time"
)

func getDate(res *http.Response) *time.Time {
	date, err := http.ParseTime(res.Header.Get(headerDate))
	if err != nil {
		return nil
	}
	return &date
}

func getAge(res *http.Response) time.Duration {
	date := getDate(res)
	if date == nil {
		return 0
	}
	age := time.Since(*date)
	if age < 0 {
		return 0
	}
	return age
}

func getTTL(res *http.Response) time.Duration {
	expires, err := http.ParseTime(res.Header.Get(headerExpires))
	if err != nil {
		return 0
	}
	return expires.Sub(time.Now())
}

func updateHeaders(res *http.Response, elapsedTime time.Duration, ttl time.Duration) {
	var resDate time.Time
	if dateVal := getDate(res); dateVal != nil {
		resDate = *dateVal
	} else {
		resDate = time.Now().Add(-elapsedTime)
		res.Header.Set(headerDate, resDate.Format(http.TimeFormat))
	}
	res.Header.Set(headerExpires, resDate.Add(ttl).Format(http.TimeFormat))
}
