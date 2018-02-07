package models

import "time"

type Message32Bit struct {
	Uid						string		`orm:"pk;type(uuid);size(36);index"`

	CreatedDate				time.Time	`orm:"type(datetime);auto_now_add;index"`

	TerminalCode 			int32		`orm:"index"`
	TerminalSetId			int32		`orm:"index"`
	ServerDate				string		`orm:"size(19)"`
	Version					int32		`orm:"index"`
	SerialNumber			int32

	ErrorCode				string		`orm:"size(2)"`

	Channel1				float32		`orm:"column(channel_1)"`
	Channel2				float32		`orm:"column(channel_2)"`
	Channel3				float32		`orm:"column(channel_3)"`
	Channel4				float32		`orm:"column(channel_4)"`
	Channel5				float32		`orm:"column(channel_5)"`
	Channel6				float32		`orm:"column(channel_6)"`
	Channel7				float32		`orm:"column(channel_7)"`
	Channel8				float32		`orm:"column(channel_8)"`
	Channel9				float32		`orm:"column(channel_9)"`
	Channel10				float32		`orm:"column(channel_10)"`
	Channel11				float32		`orm:"column(channel_11)"`
	Channel12				float32		`orm:"column(channel_12)"`
	Channel13				float32		`orm:"column(channel_13)"`
	Channel14				float32		`orm:"column(channel_14)"`
	Channel15				float32		`orm:"column(channel_15)"`
	Channel16				float32		`orm:"column(channel_16)"`
	Channel17				float32		`orm:"column(channel_17)"`
	Channel18				float32		`orm:"column(channel_18)"`
	Channel19				float32		`orm:"column(channel_19)"`
	Channel20				float32		`orm:"column(channel_20)"`
	Channel21				float32		`orm:"column(channel_21)"`
	Channel22				float32		`orm:"column(channel_22)"`
	Channel23				float32		`orm:"column(channel_23)"`
	Channel24				float32		`orm:"column(channel_24)"`
	Channel25				float32		`orm:"column(channel_25)"`
	Channel26				float32		`orm:"column(channel_26)"`
	Channel27				float32		`orm:"column(channel_27)"`
	Channel28				float32		`orm:"column(channel_28)"`
	Channel29				float32		`orm:"column(channel_29)"`
	Channel30				float32		`orm:"column(channel_30)"`
	Channel31				float32		`orm:"column(channel_31)"`
	Channel32				float32		`orm:"column(channel_32)"`
	Channel33				float32		`orm:"column(channel_33)"`
	Channel34				float32		`orm:"column(channel_34)"`
	Channel35				float32		`orm:"column(channel_35)"`
	Channel36				float32		`orm:"column(channel_36)"`
	Channel37				float32		`orm:"column(channel_37)"`
	Channel38				float32		`orm:"column(channel_38)"`
	Channel39				float32		`orm:"column(channel_39)"`
	Channel40				float32		`orm:"column(channel_40)"`
	Channel41				float32		`orm:"column(channel_41)"`
	Channel42				float32		`orm:"column(channel_42)"`
	Channel43				float32		`orm:"column(channel_43)"`
	Channel44				float32		`orm:"column(channel_44)"`

	Status					int
}

func (msg *Message32Bit) TableName() string {
	return "message_32bit"
}