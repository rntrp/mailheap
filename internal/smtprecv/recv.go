package smtprecv

import (
	"io"
	"log/slog"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/rntrp/mailheap/internal/config"
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
	slog.Info("SMTP command", "uuid", s.uuid.String(),
		"command", "AUTH PLAIN", "user", username)
	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	if config.IsSMTPAuthRequired() && !s.auth {
		return smtp.ErrAuthRequired
	}
	mode := "US-ASCII"
	if opts.UTF8 {
		mode = "SMTPUTF8"
	}
	body := smtp.Body7Bit
	if len(opts.Body) > 0 {
		body = opts.Body
	}
	slog.Info("SMTP command", "uuid", s.uuid.String(),
		"command", "MAIL FROM", "from", from, "body", body,
		"mode", mode, "size", opts.Size, "envelope", opts.EnvelopeID)
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if config.IsSMTPAuthRequired() && !s.auth {
		return smtp.ErrAuthRequired
	}
	slog.Info("SMTP command", "uuid", s.uuid.String(),
		"command", "RCPT TO", "to", to, "type", opts.OriginalRecipientType,
		"recipient", opts.OriginalRecipient)
	return nil
}

var invalidContent = &smtp.SMTPError{
	Code:         554,
	EnhancedCode: smtp.EnhancedCode{5, 6, 0},
	Message:      "Invalid message content",
}

func (s *session) Data(r io.Reader) error {
	if config.IsSMTPAuthRequired() && !s.auth {
		return smtp.ErrAuthRequired
	}
	start := time.Now()
	slog.Info("SMTP command", "uuid", s.uuid.String(),
		"command", "DATA")
	d := &readerDecorator{delegate: r}
	if err := s.addMailSvc.StoreMail(d); err != nil {
		slog.Error("SMTP: failed to store mail", "uuid", s.uuid,
			"error", err.Error())
		return invalidContent
	}
	elapsed := time.Since(start)
	slog.Info("SMTP command", "uuid", s.uuid.String(),
		"command", "<CR><LF>.<CR><LF>", "length", d.length,
		"elapsed", elapsed)
	return nil
}

func (s *session) Reset() {
	slog.Info("SMTP command", "uuid", s.uuid.String(), "command", "RSET")
}

func (s *session) Logout() error {
	slog.Info("SMTP command", "uuid", s.uuid.String(), "command", "QUIT")
	return nil
}

func Init(addMailSvc msg.StoreMailSvc) *smtp.Server {
	s := smtp.NewServer(&recv{
		username:   config.GetSMTPUsername(),
		password:   config.GetSMTPPassword(),
		addMailSvc: addMailSvc,
	})
	s.Network = config.GetSMTPNetworkType()
	s.Addr = config.GetSMTPAddress()
	s.Domain = config.GetSMTPDomain()
	s.ReadTimeout = config.GetSMTPReadTimeout()
	s.WriteTimeout = config.GetSMTPWriteTimeout()
	s.MaxMessageBytes = config.GetSMTPMaxMessageBytes()
	s.MaxRecipients = int(config.GetSMTPMaxRecipients())
	s.MaxLineLength = int(config.GetSMTPMaxLineLength())
	s.AllowInsecureAuth = config.IsSMTPAllowInsecureAuth()
	s.EnableSMTPUTF8 = config.IsSMTPEnableSMTPUTF8()
	s.EnableBINARYMIME = config.IsSMTPEnableBINARYMIME()
	s.EnableDSN = config.IsSMTPEnableDSN()
	s.EnableREQUIRETLS = config.IsSMTPEnableREQUIRETLS()
	return s
}
