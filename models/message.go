package models

import (
	"time"
)

type Message struct {
	MyUidObject

	Formatter		*MessageFormatter	`orm:"rel(fk)"`
	Header			string				`orm:"size(20)"`
	ContentLength	int32
	Content			string
	SerialNumber	int32
	Status			string				`orm:"size(20)"`
	Identifier		string				`orm:"size(20)"`
	ResponseCode	string				`orm:"size(20)"`
	ServerDate		time.Time			`orm:"type(datetime);auto_now_add;index"`
	Password		string				`orm:"size(60)"`
	Crc				string				`orm:"size(20)"`
	Tailer			string				`orm:"size(20)"`
	Body			string
}
