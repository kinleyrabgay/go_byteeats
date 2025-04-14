package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// WrapResponseWriter wraps http.ResponseWriter to capture the status code
type WrapResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func NewWrapResponseWriter(w http.ResponseWriter, protoMajor int) *WrapResponseWriter {
	return &WrapResponseWriter{ResponseWriter: w}
}

func (w *WrapResponseWriter) Status() int {
	return w.status
}

func (w *WrapResponseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *WrapResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func (w *WrapResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("underlying ResponseWriter does not implement http.Hijacker")
	}
	return h.Hijack()
}

func (w *WrapResponseWriter) Flush() {
	f, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		f.Flush()
	}
}
