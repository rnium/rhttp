package dynamicdata

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/rnium/rhttp/internal/build"
	"github.com/rnium/rhttp/internal/logic"
	"github.com/rnium/rhttp/pkg/rhttp"
)

var DEFAULT_BASE64_DATA = "ckh0dHBiaW4gaXMgYXdlc29tZQ=="

func decodeBase64Handler(r *rhttp.Request) *rhttp.Response {
	value, _ := r.Param("value")
	payload, err := base64.StdEncoding.DecodeString(value)
	statusCode := rhttp.StatusOK
	if err != nil {
		payload = fmt.Appendf(nil, "Incorrect Base64 data try: %s", DEFAULT_BASE64_DATA)
		statusCode = rhttp.StatusBadRequest
	}
	res := rhttp.NewResponse(statusCode, payload)
	_ = res.SetHeader("Content-type", "text/plain")
	return res
}

func uuidGenHandler(_ *rhttp.Request) *rhttp.Response {
	uid := buildUUID()
	data := map[string]string{
		"uuid": uid,
	}
	return rhttp.ResponseJSON(200, data)
}

func bytesHandler(r *rhttp.Request) *rhttp.Response {
	value, _ := r.Param("n")
	n, err := strconv.Atoi(value)
	if err != nil {
		rhttp.ResponseJSON(200, "expected value to be integer")
	}
	n = min(n, 102400)
	data := nRandomBytes(n)
	res := rhttp.NewResponse(200, data)
	_ = res.SetHeader("Content-Type", "application/octet-stream")
	return res
}

func delayHandler(r *rhttp.Request) *rhttp.Response {
	value, _ := r.Param("delay")
	delay, err := strconv.Atoi(value)
	if err != nil {
		rhttp.ResponseJSON(200, "expected value to be integer")
	}
	delay = min(delay, 10)
	time.Sleep(time.Second * time.Duration(delay))
	var payload any
	if logic.IsReadMethod(r.RequestLine.Method) {
		payload = build.BuildReadData(r)
	} else {
		payload = build.BuildWriteData(r)
	}
	return rhttp.ResponseJSON(200, payload)
}

func dripHandler(r *rhttp.Request) *rhttp.Response {
	duration, delay, numbytes, statusCode := getDripParams(r)
	time.Sleep(time.Second * time.Duration(delay))
	dr := newDripReader(numbytes, duration)
	res := rhttp.NewChunkedResponse(statusCode, dr)
	_ = res.SetHeader("Content-Type", "application/octet-stream")
	return res
}

func streamHandler(r *rhttp.Request) *rhttp.Response {
	n := nParam(r)
	cr := newChunkedReader(n, r)
	res := rhttp.NewChunkedResponse(200, cr)
	_ = res.SetHeader("Content-Type", "application/json")
	return res
}
