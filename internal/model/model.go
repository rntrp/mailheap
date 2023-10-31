package model

import "time"

var BasicMail = []string{"id", "created", "date", "subject", "from", "to", "cc", "bcc", "size"}

const Id = "id"
const Mime = "mime"

type Mail struct {
	Id      int64     `xorm:"pk" json:"id"`
	Created time.Time `xorm:"created index" json:"created"`
	Date    time.Time `xorm:"timestamp index" json:"date"`
	Subject string    `xorm:"mediumtext" json:"subject"`
	From    string    `xorm:"mediumtext" json:"from"`
	To      string    `xorm:"mediumtext" json:"to"`
	Cc      string    `xorm:"mediumtext" json:"cc"`
	Bcc     string    `xorm:"mediumtext" json:"bcc"`
	Size    int32     `xorm:"index" json:"size"`
	Mime    string    `xorm:"longtext" json:"mime,omitempty"`
}
