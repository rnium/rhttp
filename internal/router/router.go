package router

import (
	"fmt"
	"slices"

	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
)

const MethodGet = "GET" // methods are case sensitive, RFC 9110#section-9.1-5
const MethodPost = "POST"
const MethodPut = "PUT"
const MethodPatch = "PATCH"
const MethodDelete = "DELETE"

type Handler func(*request.Request) *response.Response

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

func (r *Router) getView(target string) (*View, request.Params) {
	node, params := r.findTrailerNode(target)
	if node == nil || node.view == nil {
		return nil, params
	}

	return node.view, params
}

func (r *Router) GetHandler(request *request.Request) Handler {
	rl := request.RequestLine
	path, query_params := parseTarget(rl.Target)
	view, params := r.getView(path)
	request.SetAllParams(params, query_params)
	if view == nil {
		return NewErrorHandler(response.StatusNotFound)
	}
	if !slices.Contains(view.methods, rl.Method) {
		return NewErrorHandler(response.StatusMethodNotAllowed)
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
