// A totally pro http handler. One that isn't particularly pro at all,
// it was just handy and I needed it in multiple projects.
package prohttphandler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// The actuall http handler
type ProHttpHandler struct {
	fsHandler             http.Handler
	exactMatchHandleFuncs map[string]func(http.ResponseWriter, *http.Request)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

// ServeHTTP will do one of three things:
//   * It will check if there is a handler func registered for the *exact* path of the request
//   * If there is no exact match, it will try and find a file in the static asset path to serve (no dir listings though)
//   * If neither of the above, it will 404 in an ugly way
// If the client sends an Accept-Encoding of gzip, then they'll get it gzipped.
// Keep an eye out for weird Content-Type issues if you're gzipping, as this will try and auto-guess
// the content type and may not be 100% accurate.
func (h *ProHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.handleRequest(w, r)
		return
	}
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
	h.handleRequest(gzr, r)
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

func (h *ProHttpHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
	handleFunc, handleFuncPresent := h.exactMatchHandleFuncs[r.URL.Path]
	if handleFuncPresent == true {
		handleFunc(w, r)
	} else if !strings.HasSuffix(r.URL.Path, "/") {
		h.fsHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
