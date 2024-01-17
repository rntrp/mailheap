package rest

import "net/http"

func addSecurityHeaders(hdr http.Header) {
	hdr.Add("Content-Security-Policy", "default-src 'self'; "+
		"img-src * data:; "+
		"style-src 'self' 'unsafe-inline';"+
		"frame-ancestors 'none';"+
		"base-uri 'self';")
	hdr.Add("Permissions-Policy", "accelerometer=(), "+
		"ambient-light-sensor=(), "+
		"autoplay=(), "+
		"battery=(), "+
		"camera=(), "+
		"cross-origin-isolated=(), "+
		"display-capture=(), "+
		"document-domain=(), "+
		"encrypted-media=(), "+
		"execution-while-not-rendered=(), "+
		"execution-while-out-of-viewport=(), "+
		"fullscreen=(), "+
		"geolocation=(), "+
		"gyroscope=(), "+
		"keyboard-map=(), "+
		"magnetometer=(), "+
		"microphone=(), "+
		"midi=(), "+
		"navigation-override=(), "+
		"payment=(), "+
		"picture-in-picture=(), "+
		"publickey-credentials-get=(), "+
		"screen-wake-lock=(), "+
		"sync-xhr=(self), "+
		"usb=(), "+
		"web-share=(), "+
		"xr-spatial-tracking=()")
	hdr.Add("Referrer-Policy", "no-referrer")
	hdr.Add("X-Content-Type-Options", "nosniff")
	hdr.Add("X-Frame-Options", "DENY")
}
