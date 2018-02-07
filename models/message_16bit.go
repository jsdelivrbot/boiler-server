package models

import "time"

type Message16Bit struct {
	Uid						string		`orm:"pk;type(uuid);size(36);index"`

	CreatedDate				time.Time	`orm:"type(datetime);auto_now_add;index"`

	TerminalCode 			int32		`orm:"index"`
	TerminalSetId			int32		`orm:"index"`
	ServerDate				string		`orm:"size(19)"`
	MessageType				string		`orm:"type(char);size(2)"`
	Version					int32		`orm:"index"`
	SerialNumber			int32

	ErrorCode				string		`orm:"size(2)"`

	ChannelTemperature1		int32		`orm:"column(channel_temperature_1)"`
	ChannelTemperature2		int32		`orm:"column(channel_temperature_2)"`
	ChannelTemperature3		int32		`orm:"column(channel_temperature_3)"`
	ChannelTemperature4		int32		`orm:"column(channel_temperature_4)"`
	ChannelTemperature5		int32		`orm:"column(channel_temperature_5)"`
	ChannelTemperature6		int32		`orm:"column(channel_temperature_6)"`
	ChannelTemperature7		int32		`orm:"column(channel_temperature_7)"`
	ChannelTemperature8		int32		`orm:"column(channel_temperature_8)"`
	ChannelTemperature9		int32		`orm:"column(channel_temperature_9)"`
	ChannelTemperature10	int32		`orm:"column(channel_temperature_10)"`
	ChannelTemperature11	int32		`orm:"column(channel_temperature_11)"`
	ChannelTemperature12	int32		`orm:"column(channel_temperature_12)"`

	ChannelAnalog1			int32		`orm:"column(channel_analog_1)"`
	ChannelAnalog2			int32		`orm:"column(channel_analog_2)"`
	ChannelAnalog3			int32		`orm:"column(channel_analog_3)"`
	ChannelAnalog4			int32		`orm:"column(channel_analog_4)"`
	ChannelAnalog5			int32		`orm:"column(channel_analog_5)"`
	ChannelAnalog6			int32		`orm:"column(channel_analog_6)"`
	ChannelAnalog7			int32		`orm:"column(channel_analog_7)"`
	ChannelAnalog8			int32		`orm:"column(channel_analog_8)"`
	ChannelAnalog9			int32		`orm:"column(channel_analog_9)"`
	ChannelAnalog10			int32		`orm:"column(channel_analog_10)"`
	ChannelAnalog11			int32		`orm:"column(channel_analog_11)"`
	ChannelAnalog12			int32		`orm:"column(channel_analog_12)"`

	ChannelSwitchIn1To16	int32		`orm:"column(channel_switch_in_1_16)"`
	ChannelSwitchIn17To32	int32		`orm:"column(channel_switch_in_17_32)"`
	ChannelSwitchOut1To16	int32		`orm:"column(channel_switch_out_1_16)"`

	ChannelCalculate1		float32		`orm:"column(channel_calculate_1)"`
	ChannelCalculate2		float32		`orm:"column(channel_calculate_2)"`
	ChannelCalculate3		float32		`orm:"column(channel_calculate_3)"`
	ChannelCalculate4		float32		`orm:"column(channel_calculate_4)"`
	ChannelCalculate5		float32		`orm:"column(channel_calculate_5)"`
	ChannelCalculate6		float32		`orm:"column(channel_calculate_6)"`
	ChannelCalculate7		float32		`orm:"column(channel_calculate_7)"`
	ChannelCalculate8		float32		`orm:"column(channel_calculate_8)"`
	ChannelCalculate9		float32		`orm:"column(channel_calculate_9)"`
	ChannelCalculate10		float32		`orm:"column(channel_calculate_10)"`
	ChannelCalculate11		float32		`orm:"column(channel_calculate_11)"`
	ChannelCalculate12		float32		`orm:"column(channel_calculate_12)"`

	Reserved1				int32		`orm:"column(reserved_1)"`
	Reserved2				int32		`orm:"column(reserved_2)"`
	Reserved3				int32		`orm:"column(reserved_3)"`
	Reserved4				int32		`orm:"column(reserved_4)"`

	Status					int
}

func (msg *Message16Bit) TableName() string {
	return "message_16bit"
}