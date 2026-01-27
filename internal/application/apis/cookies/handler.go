package cookies

import "github.com/rnium/rhttp/pkg/rhttp"

func cookiesHandler(r *rhttp.Request) *rhttp.Response {
	cookies := r.Cookies()
	payload := map[string]any{
		"cookies": cookies,
	}
	return rhttp.ResponseJSON(200, payload)
}

func setCookieHandler(r *rhttp.Request) *rhttp.Response {
	cookies := make(map[string]string)
	r.QParamForEach(func(name, value string) {
		cookies[name] = value
	})
	payload := map[string]any{
		"cookies": cookies,
	}
	res := rhttp.ResponseJSON(200, payload)
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
