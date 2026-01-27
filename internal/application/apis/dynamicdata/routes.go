package dynamicdata

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/base64/:value", decodeBase64Handler)
	r.Get("/uuid", uuidGenHandler)
	r.Get("/bytes/:n", bytesHandler)
	r.Get("/delay/:delay", delayHandler)
	r.Post("/delay/:delay", delayHandler)
	r.Put("/delay/:delay", delayHandler)
	r.Patch("/delay/:delay", delayHandler)
	r.Delete("/delay/:delay", delayHandler)
	r.Get("/drip", dripHandler)
	r.Get("/stream/:n", streamHandler)
}
