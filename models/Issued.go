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
type IssuedDateBit struct {
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
	Func int
	Byte        int
	Modbus   int
}

type IssuedSwitch struct {
	Channel  *RuntimeParameterChannelConfig     `orm:"pk;rel(fk)"`
	Func   int
	Modbus    int
	BitAddress  int
}

type IssuedSwitchBurning struct {
	Terminal      *Terminal         `orm:"pk;rel(fk)"`
	Func      int
	Modbus        int
	BitAddress    int
}

type IssuedCommunication struct {
	Terminal *Terminal    `orm:"pk;rel(fk)"`
	BaudRate   int
	DataBit    int
	StopBit    int
	CheckBit   int
	CorrespondType   int
	SubAddress   int
	HeartBeat    int
}