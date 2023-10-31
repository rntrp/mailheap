package rest

import "net/http"

func addSecurityHeaders(hdr http.Header) {
	hdr.Add("Referrer-Policy", "no-referrer")
	hdr.Add("X-Content-Type-Options", "nosniff")
	hdr.Add("X-Frame-Options", "DENY")
}
