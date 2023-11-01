package main

import (
	"log"
	"log/slog"

	"github.com/rntrp/mailheap/internal/httpsrv"
	"github.com/rntrp/mailheap/internal/msg"
	"github.com/rntrp/mailheap/internal/rest"
	"github.com/rntrp/mailheap/internal/smtprecv"
	"github.com/rntrp/mailheap/internal/storage"
	_ "modernc.org/sqlite"
)

func main() {
	slog.Info("📮 Initializing services...")
	storage, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("🥞 Database connection established")
	addMailSvc := msg.NewAddMailSvc(storage)
	receiver := smtprecv.Init(addMailSvc)
	defer receiver.Close()
	server := httpsrv.New(rest.New(storage, addMailSvc))
	defer server.Close()
	c := make(chan error)
	go func() { c <- receiver.ListenAndServe() }()
	slog.Info("📧 Receiving SMTP connections", "domain", receiver.Domain, "addr", receiver.Addr)
	go func() { c <- server.ListenAndServe() }()
	slog.Info("🌐 Listening to HTTP connections", "addr", server.Addr)
	uiHint(server.Addr)
	if err := <-c; err != nil {
		log.Println(err)
	}
}

func uiHint(addr string) {
	if len(addr) == 0 {
		slog.Info("💡 Type http://localhost in your browser for UI")
	} else if addr[0] == ':' {
		slog.Info("💡 Type http://localhost" + addr + " in your browser for UI")
	} else {
		slog.Info("💡 Type http://" + addr + " in your browser for UI")
	}
}
