package httpsrv

import (
	"log/slog"
	"net/http"
	"os"

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
	return &http.Server{Addr: config.GetHTTPTCPAddress(), Handler: logged()(r)}
}

func logged() func(next http.Handler) http.Handler {
	return httplog.Handler(httplog.NewLogger(config.GetLogServiceName(), httplog.Options{
		LogLevel:           httplog.LevelByName(config.GetLogLevel()),
		LevelFieldName:     config.GetLogLevelFieldName(),
		MessageFieldName:   config.GetLogMessageFieldName(),
		JSON:               config.IsLogJSON(),
		Concise:            config.IsLogConcise(),
		Tags:               config.GetLogTags(),
		RequestHeaders:     config.IsLogRequestHeaders(),
		HideRequestHeaders: config.GetLogHideRequestHeaders(),
		ResponseHeaders:    config.IsLogResponseHeaders(),
		QuietDownRoutes:    config.GetLogQuietDownRoutes(),
		QuietDownPeriod:    config.GetLogQuietDownPeriod(),
		TimeFieldFormat:    config.GetLogTimeFieldFormat(),
		TimeFieldName:      config.GetLogTimeFieldName(),
		SourceFieldName:    config.GetLogSourceFieldName(),
	}))
}

func shutdownFn(sig chan os.Signal) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		slog.Info("Shutdown endpoint call")
		w.WriteHeader(http.StatusAccepted)
		go func() { sig <- os.Interrupt }()
	}
}
