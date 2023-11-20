package rest

import (
	"math"
	"net/http"
	"strconv"

	"github.com/rntrp/mailheap/internal/config"
)

const minValidFileSize = 16

func setupFileSizeChecks(w http.ResponseWriter, r *http.Request) bool {
	clen, err := coerceContentLength(r.Header.Get("Content-Length"))
	if err == nil && clen < minValidFileSize {
		http.Error(w, "http: Content-Length too short for a valid eml file",
			http.StatusBadRequest)
		return false
	}
	maxReqSize := config.GetHTTPMaxRequestSize()
	if maxReqSize >= 0 {
		if err == nil && clen > maxReqSize {
			http.Error(w, "http: Content-Length too large",
				http.StatusRequestEntityTooLarge)
			return false
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxReqSize)
	}
	return true
}

func coerceContentLength(contentLength string) (int64, error) {
	return strconv.ParseInt(contentLength, 10, 64)
}

const maxMemoryBufferSize = int64(math.MaxInt64) - 1

func coerceMemoryBufferSize(memoryBufferSize int64) int64 {
	if memoryBufferSize < 0 || memoryBufferSize > maxMemoryBufferSize {
		return maxMemoryBufferSize
	}
	return memoryBufferSize
}
