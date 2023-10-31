package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/mail"
	"strings"
	"time"

	"github.com/rntrp/mailheap/internal/model"
	"github.com/rntrp/mailheap/internal/storage"
)

type StoreMailSvc interface {
	StoreMail(r io.Reader) error
}

func NewAddMailSvc(storage storage.MailStorage) StoreMailSvc {
	return &svc{storage: storage}
}

type svc struct {
	storage storage.MailStorage
}

func (s svc) StoreMail(r io.Reader) error {
	if mail, err := readMail(r); err != nil {
		return err
	} else {
		return s.storage.AddMail(mail)
	}
}

func readMail(r io.Reader) (model.Mail, error) {
	m := model.Mail{}
	b, err := io.ReadAll(r)
	if err != nil {
		return m, err
	}
	msg, err := mail.ReadMessage(bytes.NewReader(b))
	if err != nil {
		return m, err
	}
	date, err := msg.Header.Date()
	if err != nil {
		return m, err
	}
	wd := new(mime.WordDecoder)
	subject, err := wd.DecodeHeader(msg.Header.Get("Subject"))
	if err != nil {
		return m, err
	}
	to, err := address2json(msg, "To")
	if err != nil {
		return m, err
	}
	from, err := address2json(msg, "From")
	if err != nil {
		return m, err
	}
	cc, err := address2json(msg, "Cc")
	if err != nil {
		return m, err
	}
	bcc, err := address2json(msg, "Bcc")
	if err != nil {
		return m, err
	}
	m.Created = time.Now()
	m.Date = date
	m.Subject = subject
	m.From = from
	m.To = to
	m.Cc = cc
	m.Bcc = bcc
	m.Size = int32(len(b))
	m.Mime = string(b)
	return m, nil
}

func address2json(msg *mail.Message, hdr string) (string, error) {
	if len(msg.Header.Get(hdr)) == 0 {
		return "[]", nil
	}
	list, err := msg.Header.AddressList(hdr)
	if err != nil {
		return "", err
	}
	s := make([]string, len(list))
	for i := 0; i < len(list); i++ {
		if len(list[i].Name) == 0 {
			s[i] = list[i].Address
		} else {
			s[i] = fmt.Sprintf("%v <%v>", list[i].Name, list[i].Address)
		}
	}
	builder := new(strings.Builder)
	enc := json.NewEncoder(builder)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(s); err != nil {
		return "", err
	}
	return strings.TrimSpace(builder.String()), nil
}
