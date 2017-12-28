package models

type Application struct {
	MyUidObject

	Platform		string		`orm:"size(60);index"`
	App				string		`orm:"size(60);index"`
	Identity		string		`orm:"size(60);index"`

	Domain			string		`orm:"index"`
	Path			string

	AppId			string		`orm:"index"`
	AppSecret		string
	OriginId		string
	ApiToken		string
	AesKey			string
}
