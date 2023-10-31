package smtprecv

import (
	"io"
	"log/slog"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/rntrp/mailheap/internal/msg"
)

type recv struct {
	username   string
	password   string
	addMailSvc msg.StoreMailSvc
}

func (b *recv) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &session{
		uuid:       uuid,
		username:   b.username,
		password:   b.password,
		addMailSvc: b.addMailSvc,
	}, nil
}

type session struct {
	uuid       uuid.UUID
	auth       bool
	username   string
	password   string
	addMailSvc msg.StoreMailSvc
}

func (s *session) AuthPlain(username, password string) error {
	if username != s.username || password != s.password {
		return smtp.ErrAuthFailed
	}
	s.auth = true
	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	slog.Info("MAIL FROM", "session", s.uuid.String(), "from", from)
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	slog.Info("RCPT TO", "session", s.uuid.String(), "to", to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	slog.Info("DATA", "session", s.uuid.String())
	return s.addMailSvc.StoreMail(r)
}

func (s *session) Reset() {
	slog.Info("RSET", "session", s.uuid.String())
}

func (s *session) Logout() error {
	slog.Info("QUIT", "session", s.uuid.String())
	return nil
}

func Init(addMailSvc msg.StoreMailSvc) *smtp.Server {
	s := smtp.NewServer(&recv{
		username:   "username",
		password:   "password",
		addMailSvc: addMailSvc,
	})
	s.Addr = ":2525"
	s.Domain = "localhost"
	s.ReadTimeout = 60 * time.Second
	s.WriteTimeout = 60 * time.Second
	s.MaxMessageBytes = 50 * 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	s.EnableSMTPUTF8 = true
	return s
}
