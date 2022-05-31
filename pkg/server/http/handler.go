package http

import (
	"net/http"
	"unsafe"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	PATCH  HttpMethod = "PATCH"
	HEAD   HttpMethod = "HEAD"
)

type HttpFunc func(w http.ResponseWriter, r *http.Request)

type handler struct {
	path   string
	method HttpMethod
	get    func(w http.ResponseWriter, r *http.Request)
	post   func(w http.ResponseWriter, r *http.Request)
	put    func(w http.ResponseWriter, r *http.Request)
	delete func(w http.ResponseWriter, r *http.Request)
	patch  func(w http.ResponseWriter, r *http.Request)
	head   func(w http.ResponseWriter, r *http.Request)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch h.method {
	case GET:
		h.get(w, r)
		break
	case POST:
		h.post(w, r)
		break
	case PUT:
		h.put(w, r)
		break
	case DELETE:
		h.delete(w, r)
		break
	case PATCH:
		h.patch(w, r)
		break
	case HEAD:
		h.head(w, r)
		break
	}
}

var body404 = "404"
var fun404 = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body404)))
}
var body200 = "200"

var handlers = make(map[string]*handler, 16)

func Register(path string, method HttpMethod, f HttpFunc) {
	h := &handler{
		path:   path,
		method: method,
		get:    trueOrDefault(method == GET, f, fun404),
		post:   trueOrDefault(method == POST, f, fun404),
		put:    trueOrDefault(method == PUT, f, fun404),
		delete: trueOrDefault(method == DELETE, f, fun404),
		patch:  trueOrDefault(method == PATCH, f, fun404),
		head:   trueOrDefault(method == HEAD, f, fun404),
	}
	handlers[path] = h
}

func trueOrDefault(b bool, f HttpFunc, def HttpFunc) HttpFunc {
	if b {
		return f
	}
	return def
}
