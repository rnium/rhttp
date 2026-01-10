package request

type Params map[string]string

func NewParams() Params {
	return make(Params)
}

func (r *Request) SetAllParams(params Params, query_params Params) {
	r.params = params
	r.query_params = query_params
}

func (r *Request) Param(name string) (string, bool) {
	value, ok := r.params[name]
	return value, ok
}

func (r *Request) QParam(name string) (string, bool) {
	value, ok := r.query_params[name]
	return value, ok
}

func (r *Request) QParamForEach(f func(name, value string)) {
	for name, val := range r.query_params {
		f(name, val)
	}
}
