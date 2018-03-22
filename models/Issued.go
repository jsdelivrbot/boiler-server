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
