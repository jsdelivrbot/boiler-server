package models

import "time"

type Message16Bit struct {
	Id						int64		`orm:"auto;pk;index"`

	CreatedDate				time.Time	`orm:"type(datetime);auto_now_add;index"`

	TerminalCode 			int32		`orm:"index"`
	TerminalSetId			int32		`orm:"index"`
	ServerDate				string		`orm:"size(19)"`
	Version					int32		`orm:"index"`
	SerialNumber			int32

	ErrorCode				string		`orm:"size(2)"`

	ChannelTemperature1		int32
	ChannelTemperature2		int32
	ChannelTemperature3		int32
	ChannelTemperature4		int32
	ChannelTemperature5		int32
	ChannelTemperature6		int32
	ChannelTemperature7		int32
	ChannelTemperature8		int32
	ChannelTemperature9		int32
	ChannelTemperature10	int32
	ChannelTemperature11	int32
	ChannelTemperature12	int32

	ChannelAnalog1			int32
	ChannelAnalog2			int32
	ChannelAnalog3			int32
	ChannelAnalog4			int32
	ChannelAnalog5			int32
	ChannelAnalog6			int32
	ChannelAnalog7			int32
	ChannelAnalog8			int32
	ChannelAnalog9			int32
	ChannelAnalog10			int32
	ChannelAnalog11			int32
	ChannelAnalog12			int32

	ChannelSwitchIn1To16	int32
	ChannelSwitchIn17To32	int32
	ChannelSwitchOut1To16	int32

	ChannelCalculate1		float32
	ChannelCalculate2		float32
	ChannelCalculate3		float32
	ChannelCalculate4		float32
	ChannelCalculate5		float32
	ChannelCalculate6		float32
	ChannelCalculate7		float32
	ChannelCalculate8		float32
	ChannelCalculate9		float32
	ChannelCalculate10		float32
	ChannelCalculate11		float32
	ChannelCalculate12		float32

	Reserved1				int32
	Reserved2				int32
	Reserved3				int32
	Reserved4				int32

	Status					int
}

func (msg *Message16Bit) TableName() string {
	return "message_16bit"
}