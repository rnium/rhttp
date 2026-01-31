package cookies

import (
	"time"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func cookiesHandler(r *rhttp.Request) *rhttp.Response {
	return rhttp.ResponseJSON(
		200,
		buildCookiePayload(r.Cookies()),
	)
}

func setCookieHandler(r *rhttp.Request) *rhttp.Response {
	cookies := make(map[string]string)
	r.QParamForEach(func(name, value string) {
		cookies[name] = value
	})

	res := rhttp.ResponseJSON(
		200,
		buildCookiePayload(cookies),
	)
	for name, val := range cookies {
		res.SetCookie(
			&rhttp.Cookie{
				Name:  name,
				Value: val,
			},
		)
	}
	return res
}

func deleteCookieHandler(r *rhttp.Request) *rhttp.Response {
	var cookieNames []string
	r.QParamForEach(func(name, _ string) {
		cookieNames = append(cookieNames, name)
	})

	res := rhttp.Redirect("/cookies", false)

	expires, _ := time.Parse(
		time.RFC1123,
		"Thu, 01 Jan 1970 00:00:00 GMT",
	)
	for _, name := range cookieNames {
		res.SetCookie(
			&rhttp.Cookie{
				Name:    name,
				Value:   "",
				MaxAge:  "0",
				Expires: expires,
			},
		)
	}
	return res
}
