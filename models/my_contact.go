package models

type MyContact struct {
	PhoneNumber		string		`orm:"size(20)"`
	MobileNumber		string		`orm:"size(20)"`
	Email			string		`orm:"size(60)"`
}
