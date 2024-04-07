package httpsrv

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rntrp/mailheap/internal/config"
	"github.com/rntrp/mailheap/internal/rest"
)

func New(ctrl rest.Controller, shutdown chan os.Signal) *http.Server {
	r := http.NewServeMux()
	r.HandleFunc("GET /", ctrl.Index)
	r.HandleFunc("GET /index.html", ctrl.Index)
	r.HandleFunc("GET /favicon.ico", ctrl.IndexFaviconIco)
	r.HandleFunc("GET /favicon.svg", ctrl.IndexFaviconSvg)
	r.HandleFunc("GET /index.css", ctrl.IndexCss)
	r.HandleFunc("GET /index.js", ctrl.IndexJs)
	r.HandleFunc("GET /index.jsmimeparser.min.js", ctrl.IndexJsMimeParser)
	r.HandleFunc("GET /mail/{id}", ctrl.GetEml)
	r.HandleFunc("DELETE /mails", ctrl.DeleteMails)
	r.HandleFunc("GET /mails/{id}", ctrl.SeekMails)
	r.HandleFunc("POST /upload", ctrl.UploadMail)
	r.HandleFunc("GET /health", rest.Live)
	if config.IsHTTPEnablePrometheus() {
		r.Handle("/metrics", promhttp.Handler())
	}
	if config.IsHTTPEnableShutdown() {
		r.HandleFunc("POST /shutdown", shutdownFn(shutdown))
	}
	h := httplog.Handler(httplog.NewLogger("MAILHEAP", httplog.Options{
		Concise:         true,
		JSON:            false,
		RequestHeaders:  false,
		TimeFieldFormat: time.RFC3339,
	}))(r)
	return &http.Server{Addr: config.GetHTTPTCPAddress(), Handler: h}
}

func shutdownFn(sig chan os.Signal) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		slog.Info("Shutdown endpoint call")
		w.WriteHeader(http.StatusAccepted)
		go func() { sig <- os.Interrupt }()
	}
}
