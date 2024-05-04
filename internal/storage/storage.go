package storage

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/rntrp/mailheap/internal/idsrc"
	"github.com/rntrp/mailheap/internal/model"
	"gorm.io/gorm"
)

type MailStorage interface {
	AddMail(mail model.Mail) error
	CountMails() (int64, error)
	DeleteAllMails() (int64, error)
	DeleteMails(ids ...int64) (int64, error)
	GetMime(id int64) (string, error)
	SeekMails(int64, int) ([]model.Mail, error)
	Shutdown() error
}

type store struct {
	db    *gorm.DB
	idSrc idsrc.IdSrc
}

func New() (MailStorage, error) {
	path := filepath.Join(os.TempDir(), "mailheap.db")
	f, err := os.OpenFile(path, os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	} else if err := f.Close(); err != nil {
		slog.Warn("Closing db file failed:", "error", err.Error())
	}
	db, err := gorm.Open(sqlite.Open(path), new(gorm.Config))
	if err != nil {
		return nil, err
	} else if err := db.AutoMigrate(new(model.Mail)); err != nil {
		return nil, err
	}
	return &store{
		db:    db,
		idSrc: idsrc.New(),
	}, nil
}

func (s *store) AddMail(mail model.Mail) error {
	id, err := s.idSrc.Gen()
	if err != nil {
		return err
	}
	mail.Id = id
	s.db.Create(&mail)
	return nil
}

func (s *store) CountMails() (int64, error) {
	cnt := int64(0)
	err := s.db.Model(new(model.Mail)).Count(&cnt).Error
	return cnt, err
}

func (s *store) DeleteAllMails() (int64, error) {
	tx := s.db.Delete(new(model.Mail), "id>=?", 0)
	return tx.RowsAffected, tx.Error
}

func (s *store) DeleteMails(ids ...int64) (int64, error) {
	tx := s.db.Delete(new(model.Mail), ids)
	return tx.RowsAffected, tx.Error
}

func (s *store) GetMime(id int64) (string, error) {
	m := new(model.Mail)
	err := s.db.Select(model.Mime).First(m, id).Error
	return m.Mime, err
}

func (s *store) SeekMails(afterId int64, limit int) ([]model.Mail, error) {
	mails := make([]model.Mail, 0, limit)
	err := s.db.Order("id DESC").
		Limit(limit).
		Find(&mails, "id<?", afterId).
		Error
	return mails, err
}

func (s *store) Shutdown() error {
	if db, err := s.db.DB(); err != nil {
		return err
	} else {
		return db.Close()
	}
}
