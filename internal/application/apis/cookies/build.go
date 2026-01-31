package cookies

func buildCookiePayload(cookies map[string]string) map[string]any {
	payload := map[string]any{
		"cookies": cookies,
	}
	return payload
}