package models

import "time"

type Message32Bit struct {
	Id						int64		`orm:"auto;pk;index"`

	CreatedDate				time.Time	`orm:"type(datetime);auto_now_add;index"`

	TerminalCode 			int32		`orm:"index"`
	TerminalSetId			int32		`orm:"index"`
	ServerDate				string		`orm:"size(19)"`
	Version					int32		`orm:"index"`
	SerialNumber			int32

	ErrorCode				string		`orm:"size(2)"`

	Channel1				int32
	Channel2				int32
	Channel3				int32
	Channel4				int32
	Channel5				int32
	Channel6				int32
	Channel7				int32
	Channel8				int32
	Channel9				int32
	Channel10				int32
	Channel11				int32
	Channel12				int32
	Channel13				int32
	Channel14				int32
	Channel15				int32
	Channel16				int32
	Channel17				int32
	Channel18				int32
	Channel19				int32
	Channel20				int32
	Channel21				int32
	Channel22				int32

	Status					int
}

func (msg *Message32Bit) TableName() string {
	return "message_32bit"
}