package httpsrv

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func logged() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			uuid := uuid.NewString()
			slog.Info("HTTP request", "uuid", uuid, "method", r.Method,
				"url", r.URL, "protocol", r.Proto, "host", r.Host,
				"remote", r.RemoteAddr, "length", r.ContentLength)
			d := &rwDecorator{delegate: w, status: http.StatusOK}
			defer func() {
				elapsed := time.Since(start)
				slog.Info("HTTP response", "uuid", uuid, "status", d.status,
					"length", d.length, "elapsed", elapsed)
			}()
			next.ServeHTTP(d, r)
		})
	}
}

type rwDecorator struct {
	delegate http.ResponseWriter
	status   int
	length   uint64
}

func (i *rwDecorator) Header() http.Header {
	return i.delegate.Header()
}

func (i *rwDecorator) Write(bytes []byte) (int, error) {
	length, err := i.delegate.Write(bytes)
	if err == nil {
		i.length += uint64(length)
	}
	return length, err
}

func (i *rwDecorator) WriteHeader(statusCode int) {
	i.status = statusCode
	i.delegate.WriteHeader(statusCode)
}
