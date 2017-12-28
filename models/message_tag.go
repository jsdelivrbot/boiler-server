package models

type MessageTag struct {
	MyUidObject

	TagId		int32		`orm:"index`
	DataType	string
	Column		string
	Content		string
	Length		int32
}