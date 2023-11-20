package httpsrv

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rntrp/mailheap/internal/config"
	"github.com/rntrp/mailheap/internal/rest"
)

func New(ctrl rest.Controller, shutdown chan os.Signal) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", ctrl.Index)
	r.Get("/index.html", ctrl.Index)
	r.Get("/mail/{id}", ctrl.GetEml)
	r.Delete("/mails", ctrl.DeleteMails)
	r.Get("/mails/{id}", ctrl.SeekMails)
	r.Post("/upload", ctrl.UploadMail)
	r.Get("/health", rest.Live)
	if config.IsHTTPEnablePrometheus() {
		r.Handle("/metrics", promhttp.Handler())
	}
	if config.IsHTTPEnableShutdown() {
		r.Post("/shutdown", shutdownFn(shutdown))
	}
	return &http.Server{Addr: config.GetHTTPTCPAddress(), Handler: r}
}

func shutdownFn(sig chan os.Signal) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		slog.Info("Shutdown endpoint call")
		w.WriteHeader(http.StatusAccepted)
		go func() { sig <- os.Interrupt }()
	}
}
