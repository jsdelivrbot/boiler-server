package models
type TermByte struct{
	Id int `orm:"pk"`
	Name string
	Value int
}
type TermFunctionCode struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type BaudRate struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type CorrespondType struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type DateBit struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type HeartbeatPacket struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type ParityBit struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type SlaveAddress struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
type StopBit struct {
	Id int `orm:"pk"`
	Name string
	Value int
}
