package logger

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		statusColor := Green
		if wrapper.status >= 400 && wrapper.status < 500 {
			statusColor = Yellow
		} else if wrapper.status >= 500 {
			statusColor = Red
		}

		fmt.Printf("%s %s%s%s %s -> %s%d %s%s (%v)\n",
			getTimestamp(),
			Cyan, r.Method, Reset,
			r.URL.Path,
			statusColor, wrapper.status, http.StatusText(wrapper.status), Reset,
			time.Since(start).Round(time.Millisecond),
		)
	})
}
