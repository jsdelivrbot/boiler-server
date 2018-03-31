package models
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

type IssuedAnalogue struct {
	Channel  *RuntimeParameterChannelConfig     `orm:"pk;rel(fk)"`
	Function   *IssuedFunctionCode               `orm:"rel(fk)"`
	Byte       *IssuedByte                     `orm:"rel(fk)"`
	Modbus   int
}

type IssuedSwitch struct {
	Channel  *RuntimeParameterChannelConfig     `orm:"pk;rel(fk)"`
	Function   *IssuedFunctionCode             `orm:"rel(fk)"`
	Modbus    int
	BitAddress  int
}

type IssuedSwitchBurn struct {
	Terminal      *Terminal         `orm:"pk;rel(fk)"`
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