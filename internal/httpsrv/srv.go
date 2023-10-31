package httpsrv

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rntrp/mailheap/internal/rest"
)

func New(ctrl rest.Controller) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", ctrl.Index)
	r.Get("/index.html", ctrl.Index)
	r.Get("/mail/{id}", ctrl.GetEml)
	r.Delete("/mails", ctrl.DeleteMails)
	r.Get("/mails/{id}", ctrl.SeekMails)
	r.Post("/upload", ctrl.UploadMail)
	return &http.Server{Addr: ":8080", Handler: r}
}
