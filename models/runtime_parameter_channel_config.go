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
	ChannelNumber		int32					`orm:"index"`

	Length				int32

	Signed				bool
	NegativeThreshold	int32

	Status				int32
	SequenceNumber		int32

	Scale				float32

	IsDefault			bool					`orm:"index"`
}

const (
	CHANNEL_TYPE_DEFAULT		= 0
	CHANNEL_TYPE_TEMPERATURE	= 1
	CHANNEL_TYPE_ANALOG			= 2
	CHANNEL_TYPE_SWITCH			= 3
	CHANNEL_TYPE_CALCULATE		= 4
)

const (
	CHANNEL_STATUS_DEFAULT		= 0
	CHANNEL_STATUS_SHOW			= 1
	CHANNEL_STATUS_HIDE			= 2
)

func ChannelName(cType, num int32) string {
	var name string
	switch cType {
	case CHANNEL_TYPE_TEMPERATURE:
		name = "Temper" + strconv.FormatInt(int64(num), 10) + "_channel"
	case CHANNEL_TYPE_ANALOG:
		name = "Analog" + strconv.FormatInt(int64(num), 10) + "_channel"
	case CHANNEL_TYPE_SWITCH:
		name = "Switch_in_" + strconv.FormatInt(int64(num), 10) + "_channel"
	case CHANNEL_TYPE_CALCULATE:
		name = "C" + strconv.FormatInt(int64(num), 10) + "_calculate_parm"
	case CHANNEL_TYPE_DEFAULT:
		fallthrough
	default:

	}

	return name
}