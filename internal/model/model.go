package model

import "time"

var BasicMail = []string{"id", "created", "date", "subject", "from", "to", "cc", "bcc", "size"}

const Id = "id"
const Mime = "mime"

type Mail struct {
	Id      int64     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Created time.Time `gorm:"index" json:"created"`
	Date    time.Time `gorm:"index" json:"date"`
	Subject string    `gorm:"text" json:"subject"`
	From    string    `gorm:"text" json:"from"`
	To      string    `gorm:"text" json:"to"`
	Cc      string    `gorm:"text" json:"cc"`
	Bcc     string    `gorm:"text" json:"bcc"`
	Size    int32     `gorm:"index" json:"size"`
	Mime    string    `gorm:"text" json:"mime,omitempty"`
}
