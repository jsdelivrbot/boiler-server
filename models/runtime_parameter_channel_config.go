package models

import (
	"strconv"
)

type RuntimeParameterChannelConfig struct {
	MyUidObject

	Parameter			*RuntimeParameter		`orm:"rel(fk);index"`
	Boiler				*Boiler					`orm:"rel(fk);index;null"`
	Terminal			*Terminal				`orm:"rel(fk);index;null"`

	ChannelType			int32					`orm:"index"`

	ChannelCol			int32					`orm:"index"`
	ChannelNumber		int32					`orm:"index"`


	Length				int32

	Signed				bool
	NegativeThreshold	int32

	Status				int32
	SequenceNumber		int32

	SwitchStatus		int32					`orm:"index"`

	Scale				float32

	Ranges				[]*RuntimeParameterChannelConfigRange	`orm:"reverse(many)"`

	IsDefault			bool					`orm:"index"`
}

const (
	CHANNEL_TYPE_DEFAULT		= 0
	CHANNEL_TYPE_TEMPERATURE	= 1
	CHANNEL_TYPE_ANALOG			= 2
	CHANNEL_TYPE_SWITCH			= 3
	CHANNEL_TYPE_CALCULATE		= 4
	CHANNEL_TYPE_RANGE			= 5
)

const (
	CHANNEL_STATUS_DEFAULT		= 0
	CHANNEL_STATUS_SHOW			= 1
	CHANNEL_STATUS_HIDE			= 2
)

const (
	CHANNEL_SWITCH_STATUS_DEFAULT		= 0
	CHANNEL_SWITCH_STATUS_RUNNING		= 1
	CHANNEL_SWITCH_STATUS_FAULT			= 2
)

func ChannelName(cType, num int32) string {
	var name string
	switch cType {
	case CHANNEL_TYPE_TEMPERATURE:
		name = "Temper" + strconv.FormatInt(int64(num), 10) + "_channel"
	case CHANNEL_TYPE_ANALOG:
		name = "Analog" + strconv.FormatInt(int64(num), 10) + "_channel"
	case CHANNEL_TYPE_SWITCH:
		var numStr string
		if num <= 16 {
			numStr = "in_1_16"
		} else if num <= 32 {
			numStr = "in_17_32"
		} else if num <= 48 {
			numStr = "out_1_16"
		}
		name = "Switch_" + numStr + "_channel"
	case CHANNEL_TYPE_CALCULATE:
		fallthrough
	case CHANNEL_TYPE_RANGE:
		name = "C" + strconv.FormatInt(int64(num), 10) + "_calculate_parm"
	case CHANNEL_TYPE_DEFAULT:
		fallthrough
	default:

	}

	return name
}

type RuntimeParameterChannelConfigRange struct {
	MyUidObject

	ChannelConfig	*RuntimeParameterChannelConfig		`orm:"rel(fk);index"`

	Min				int64			`orm:"index"`
	Max				int64			`orm:"index"`
	Value			int64			`orm:"index"`
}

func (rg *RuntimeParameterChannelConfigRange) TableUnique() [][]string {
	return [][]string {
		[]string{"ChannelConfig", "Min", "Max"},
	}
}