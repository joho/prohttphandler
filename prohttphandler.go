// A totally pro http handler. One that isn't particularly pro at all,
// it was just handy and I needed it in multiple projects.
package prohttphandler

import (
	"net/http"
	"strings"
)

// The actuall http handler
type ProHttpHandler struct {
	fsHandler             http.Handler
	exactMatchHandleFuncs map[string]func(http.ResponseWriter, *http.Request)
}

// ServeHTTP will do one of three things:
//   * It will check if there is a handler func registered for the *exact* path of the request
//   * If there is no exact match, it will try and find a file in the static asset path to serve (no dir listings though)
//   * If neither of the above, it will 404 in an ugly way
func (h *ProHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleFunc, handleFuncPresent := h.exactMatchHandleFuncs[r.URL.Path]
	if handleFuncPresent == true {
		handleFunc(w, r)
	} else if !strings.HasSuffix(r.URL.Path, "/") {
		h.fsHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// Registers a http handler func for the path you specify
func (h *ProHttpHandler) ExactMatchFunc(path string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	h.exactMatchHandleFuncs[path] = handlerFunc
}

// Create a ProHttpHandler
// You need to give it a path to use as the static assets root
func New(publicStaticAssetPath string) (h *ProHttpHandler) {
	h = new(ProHttpHandler)
	h.exactMatchHandleFuncs = map[string]func(http.ResponseWriter, *http.Request){}
	h.fsHandler = http.FileServer(http.Dir(publicStaticAssetPath))
	return
}
