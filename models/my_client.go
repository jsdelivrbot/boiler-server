package models

type MyClient struct {
	IpAddress	string		`orm:"size(20)"`
	MacAddress	string		`orm:"size(20)"`

	DeviceName	string		`orm:"size(100)"`
	DeviceModel	string		`orm:"size(100)"`
}
