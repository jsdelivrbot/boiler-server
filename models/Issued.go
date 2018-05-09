package models

import "time"

type IssuedByte struct{
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedFunctionCode struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedBaudRate struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedCorrespondType struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedDataBit struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedHeartbeatPacket struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedParityBit struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedSlaveAddress struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type IssuedStopBit struct {
	Id int `orm:"pk"`
	Name string
	Value int
}

type IssuedAnalogueSwitch struct {
	Channel  *RuntimeParameterChannelConfig     `orm:"pk;rel(fk)"`
	CreateTime time.Time                     `orm:"type(datetime);auto_now;index"`
	Function   *IssuedFunctionCode               `orm:"rel(fk)"`
	Byte       *IssuedByte                     `orm:"rel(fk)"`
	Modbus   int
	BitAddress  int
}

type IssuedSwitchDefault struct {
	Uid        string 				`orm:"pk"`
	Terminal      *Terminal         `orm:"rel(fk)"`
	CreateTime    time.Time             `orm:"type(datetime);auto_now;index"`
	ChannelType   int
	ChannelNumber int
	Function      *IssuedFunctionCode      `orm:"rel(fk)"`
	Modbus        int
	BitAddress    int
}

type IssuedCommunication struct {
	Terminal *Terminal    `orm:"pk;rel(fk)"`
	BaudRate   *IssuedBaudRate   `orm:"rel(fk)"`
	DataBit    *IssuedDataBit     `orm:"rel(fk)"`
	StopBit    *IssuedStopBit      `orm:"rel(fk)"`
	CheckBit   *IssuedParityBit       `orm:"rel(fk)"`
	CorrespondType   *IssuedCorrespondType      `orm:"rel(fk)"`
	SubAddress  *IssuedSlaveAddress        `orm:"rel(fk)"`
	HeartBeat   *IssuedHeartbeatPacket        `orm:"rel(fk)"`
}

type IssuedErrorCode struct {
	Id int `orm:"pk"`
	Remark  string
	Value string
}

type IssuedPlcAlarm struct {
	Uid string `orm:"pk"`
	Sn string
	CreateTime time.Time         `orm:"type(datetime);auto_now;index"`
	Ver int
	ChannelNumber int
	Err   *IssuedErrorCode      `orm:"rel(fk)"`
}
type IssuedBoilerStatus struct {
	Boiler *Boiler 		`orm:"pk;rel(fk)"`
	CreateTime time.Time 	`orm:"type(datetime);auto_now;index"`
	UpdateTime time.Time	`orm:"type(datetime);auto_now;index"`
	Status bool
}

type IssuedWeekInformationLog struct {
	Uid string       `orm:"pk"`
	Boiler     *Boiler 		`orm:"rel(fk)"`
	BoilerName string
	StartDate time.Time      `orm:"type(datetime);auto_now;index"`
	EndDate time.Time        `orm:"type(datetime);auto_now;index"`
	CreateTime time.Time      `orm:"type(datetime);auto_now;index"`
	AlarmCount int
	ParameterId int
	ParameterAlarmCount int
	Description string
	BoilerRuntime float64
}
