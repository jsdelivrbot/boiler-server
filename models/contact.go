package models

type Contact struct {
	MyUidObject

	ContactId		int32		`orm:"index"`

	MyContact
}
