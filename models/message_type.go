package models

type MessageType struct {
	MyUidObject

	TypeId		int32		`orm:"index"`
	From		string
}
