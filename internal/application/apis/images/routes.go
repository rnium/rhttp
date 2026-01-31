package images

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/image", imageHandler)
	r.Get("/image/jpeg", handlerFactory(acceptJpeg))
	r.Get("/image/webp", handlerFactory(acceptWebp))
	r.Get("/image/gif", handlerFactory(acceptGif))
}
