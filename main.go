package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/emersion/go-smtp"
	"github.com/rntrp/mailheap/internal/config"
	"github.com/rntrp/mailheap/internal/httpsrv"
	"github.com/rntrp/mailheap/internal/msg"
	"github.com/rntrp/mailheap/internal/rest"
	"github.com/rntrp/mailheap/internal/smtprecv"
	"github.com/rntrp/mailheap/internal/storage"
)

func main() {
	config.Load()
	slog.Info("📮 Initializing services...")
	rest.InitIndex()
	slog.Info("🗜️ Static web UI resources minified & compressed")
	storage, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("🥞 Database connection established")
	addMailSvc := msg.NewAddMailSvc(storage)
	recv := smtprecv.Init(addMailSvc)
	sig := make(chan os.Signal, 1)
	srv := httpsrv.New(rest.New(storage, addMailSvc), sig)
	shutdown := make(chan error)
	go shutdownMonitor(sig, shutdown, storage, recv, srv)
	slog.Info("🔌 Set up graceful shutdown monitor")
	out := make(chan<- error)
	go startRecv(out, recv)
	go startSrv(out, srv)
	logShutdown(<-shutdown)
	if err := storage.Shutdown(); err != nil {
		slog.Error("DB shutdown failed", "error", err.Error())
	}
}

func startRecv(out chan<- error, recv *smtp.Server) {
	slog.Info("📧 Receiving SMTP connections",
		"domain", recv.Domain,
		"addr", recv.Addr)
	out <- recv.ListenAndServe()
}

func startSrv(out chan<- error, srv *http.Server) {
	slog.Info("🌐 Listening to HTTP connections", "addr", srv.Addr)
	switch {
	case len(srv.Addr) == 0:
		slog.Info("💡 Type http://localhost in your browser for UI")
	case srv.Addr[0] == ':':
		slog.Info("💡 Type http://localhost" + srv.Addr + " in your browser for UI")
	default:
		slog.Info("💡 Type http://" + srv.Addr + " in your browser for UI")
	}
	out <- srv.ListenAndServe()
}

func shutdownMonitor(sig chan os.Signal, out chan error,
	mailStorage storage.MailStorage, switches ...shutdownSwitch) {
	timeout := config.GetShutdownTimeout()
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sigName := (<-sig).String()
	slog.Info("Shutdown signal received", "signal", sigName)
	wg := new(sync.WaitGroup)
	num := len(switches)
	err := make([]error, num+1)
	wg.Add(num)
	for i, s := range switches {
		go func(i int, s shutdownSwitch) {
			defer wg.Done()
			ctx := context.Background()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}
			err[i] = s.Shutdown(ctx)
		}(i, s)
	}
	wg.Wait()
	err[num] = mailStorage.Shutdown()
	out <- errors.Join(err...)
}

func logShutdown(err error) {
	switch err {
	case nil:
		slog.Info("Mailheap was shut down gracefully. Bye.")
	case http.ErrServerClosed:
		slog.Info(err.Error())
	default:
		slog.Error(err.Error())
	}
}

type shutdownSwitch interface {
	Shutdown(ctx context.Context) error
}
