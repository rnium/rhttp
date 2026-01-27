package dynamicdata

import (
	"strconv"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func getDripParams(r *rhttp.Request) (duration, delay float64, numbytes, statusCode int) {
	durationParam, _ := r.QParam("duration")
	nbytesParam, _ := r.QParam("numbytes")
	delayParam, _ := r.QParam("delay")
	statusParam, _ := r.QParam("code")
	duration, _ = strconv.ParseFloat(durationParam, 64)
	numbytes, _ = strconv.Atoi(nbytesParam)
	delay, _ = strconv.ParseFloat(delayParam, 64)
	statusCode, err := strconv.Atoi(statusParam)
	if err != nil {
		statusCode = 200
	}
	return
}

func nParam(r *rhttp.Request) int {
	nParam, _ := r.Param("n")
	n, _ := strconv.Atoi(nParam)
	return n
}