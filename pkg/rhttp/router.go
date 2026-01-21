package rhttp

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
	"unicode"

)

const MethodGet = "GET" // methods are case sensitive, RFC 9110#section-9.1-5
const MethodPost = "POST"
const MethodPut = "PUT"
const MethodPatch = "PATCH"
const MethodDelete = "DELETE"

type Handler func(*Request) *Response

type View struct {
	handler Handler
	methods []string
}

func NewView(handler Handler, methods ...string) *View {
	return &View{
		handler: handler,
		methods: methods,
	}
}

type Router struct {
	rootNode *Node
}

func NewRouter() *Router {
	return &Router{
		rootNode: newNode(""),
	}
}

func (r *Router) getView(target string) (*View, Params) {
	node, params := r.findTrailerNode(target)
	if node == nil || node.view == nil {
		return nil, params
	}

	return node.view, params
}

// helpers
var ErrInvalidHttpTarget = fmt.Errorf("invalid http target")

func validateTarget(target string) error {
	for _, c := range target {
		switch {
		case unicode.IsLetter(c) || unicode.IsDigit(c):
		case strings.ContainsRune("*./-_:", c):
		default:
			return ErrInvalidHttpTarget
		}
	}
	return nil
}

func newErrorHandler(statusCode int) Handler {
	return func(r *Request) *Response {
		return ErrorResponseJSON(statusCode)
	}
}

func parseTarget(target string) (string, Params) {
	parts := strings.SplitN(target, "?", 2)
	if len(parts) == 1 {
		return parts[0], nil
	}

	queryParams := NewParams()

	for pair := range strings.SplitSeq(parts[1], "&") {
		if pair == "" {
			continue
		}

		key, val, found := strings.Cut(pair, "=")
		if !found {
			continue
		}

		k, err1 := url.QueryUnescape(key)
		v, err2 := url.QueryUnescape(val)
		if err1 != nil || err2 != nil {
			continue
		}

		queryParams[k] = v
	}
	return parts[0], queryParams
}

func (r *Router) GetHandler(request *Request) Handler {
	rl := request.RequestLine
	path, query_params := parseTarget(rl.Target)
	view, params := r.getView(path)
	request.SetAllParams(params, query_params)
	if view == nil {
		return newErrorHandler(StatusNotFound)
	}
	if !slices.Contains(view.methods, rl.Method) {
		return newErrorHandler(StatusMethodNotAllowed)
	}
	return view.handler
}

func (r *Router) addView(target, method string, handler Handler) {
	err := validateTarget(target)
	if err != nil {
		panic(fmt.Errorf("Error while registering handler for target '%s'. %v", target, err))
	}
	view, _ := r.getView(target)
	if view == nil {
		node := r.insertUrl(target)
		node.view = NewView(handler, method)
	} else if !slices.Contains(view.methods, method) {
		view.methods = append(view.methods, method)
	}
}

func (r *Router) Get(target string, handler Handler) {
	r.addView(target, MethodGet, handler)
}

func (r *Router) Post(target string, handler Handler) {
	r.addView(target, MethodPost, handler)
}

func (r *Router) Put(target string, handler Handler) {
	r.addView(target, MethodPut, handler)
}

func (r *Router) Patch(target string, handler Handler) {
	r.addView(target, MethodPatch, handler)
}

func (r *Router) Delete(target string, handler Handler) {
	r.addView(target, MethodDelete, handler)
}
