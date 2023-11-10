package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rntrp/mailheap/internal/idsrc"
	"github.com/rntrp/mailheap/internal/model"
	"xorm.io/xorm"
)

type MailStorage interface {
	AddMail(mail model.Mail) error
	CountMails() (int64, error)
	DeleteAllMails() (int64, error)
	DeleteMails(ids ...int64) (int64, error)
	GetMime(id int64) (string, error)
	SeekMails(int64, int) ([]model.Mail, error)
}

type store struct {
	engine *xorm.Engine
	idSrc  idsrc.IdSrc
}

func New() (MailStorage, error) {
	path := filepath.Join(os.TempDir(), "mimedump.db")
	f, err := os.OpenFile(path, os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	} else if f.Close() != nil {
		log.Println("Closing db file failed:", err)
	}
	engine, err := xorm.NewEngine("sqlite", f.Name())
	if err != nil {
		return nil, err
	} else if err := engine.Sync(new(model.Mail)); err != nil {
		return nil, err
	}
	return &store{
		engine: engine,
		idSrc:  idsrc.New(),
	}, nil
}

func (s *store) AddMail(mail model.Mail) error {
	id, err := s.idSrc.Gen()
	if err != nil {
		return err
	}
	mail.Id = id
	rows, err := s.engine.InsertOne(&mail)
	if err == nil && rows != 1 {
		log.Println("Affected rows should be 1, but was", rows)
	}
	return err
}

func (s *store) CountMails() (int64, error) {
	return s.engine.Count(new(model.Mail))
}

func (s *store) DeleteAllMails() (int64, error) {
	return s.engine.Truncate(new(model.Mail))
}

func (s *store) DeleteMails(ids ...int64) (int64, error) {
	return s.engine.ID(ids).Delete(new(model.Mail))
}

func (s *store) GetMime(id int64) (string, error) {
	m := new(model.Mail)
	if found, err := s.engine.
		Cols(model.Mime).
		ID(id).
		Get(m); err != nil {
		return "", err
	} else if !found {
		return "", nil
	}
	return m.Mime, nil
}

func (s *store) SeekMails(afterId int64, limit int) ([]model.Mail, error) {
	m := new(model.Mail)
	rows, err := s.engine.
		Cols(model.BasicMail...).
		Where(model.Id+"<?", afterId).
		Desc(model.Id).
		Limit(limit).
		Rows(m)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mails := make([]model.Mail, 0, limit)
	for rows.Next() {
		if err := rows.Scan(m); err != nil {
			return nil, err
		}
		mails = append(mails, *m)
	}
	return mails, nil
}
