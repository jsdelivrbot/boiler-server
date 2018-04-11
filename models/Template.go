package models

import (
	"time"
)

type IssuedTemplate struct {
	Uid string `orm:"pk"`
	Name string
	CreateTime time.Time        `orm:"type(datetime);auto_now;index"`
	UpdateTime time.Time         `orm:"type(datetime);auto_now;index"`
	IsDeleted  bool		         `orm:"index"`
	Organization  *Organization  `orm:"rel(fk)"`
}

type IssuedChannelConfigTemplate struct {
	Uid string `orm:"pk"`
	CreateTime time.Time        `orm:"type(datetime);auto_now;index"`
	Parameter *RuntimeParameter  `orm:"rel(fk)"`
	Template  *IssuedTemplate    `orm:"rel(fk)"`
	ChannelType int
	ChannelNumber int
	Status int
	SequenceNumber int
	SwitchStatus int
	Func  *IssuedFunctionCode    `orm:"rel(fk)"`
	Byte  *IssuedByte             `orm:"rel(fk)"`
	BitAddress int
	Modbus int
	Ranges       []*IssuedChannelConfigRangeTemplate   `orm:"reverse(many)"`
}
type IssuedChannelConfigRangeTemplate struct {
	Uid string 	`orm:"pk"`
	Name   string
	CreateTime time.Time        `orm:"type(datetime);auto_now;index"`
	ChannelConfig  *IssuedChannelConfigTemplate   `orm:"rel(fk)"`
	Min   int
	Max   int
	Value int
}

type IssuedCommunicationTemplate struct {
	Uid string `orm:"pk"`
	Template    *IssuedTemplate        `orm:"rel(fk)"`
	BaudRate    *IssuedBaudRate			`orm:"rel(fk)"`
	DataBit     *IssuedDataBit			`orm:"rel(fk)"`
	StopBit     *IssuedStopBit			`orm:"rel(fk)"`
	CheckBit    *IssuedParityBit		`orm:"rel(fk)"`
	CorrespondType  *IssuedCorrespondType 	`orm:"rel(fk)"`
	SubAddress   *IssuedSlaveAddress		`orm:"rel(fk)"`
	HeartBeat   *IssuedHeartbeatPacket		`orm:"rel(fk)"`
}
