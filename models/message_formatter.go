package models

type MessageFormatter struct {
	MyUidObject

	FormatterId		int32			`orm:"index"`
	Type			*MessageType	`orm:"rel(fk)"`
	Tag				*MessageTag		`orm:"rel(fk)"`
	Length			int32
	SequenceNumber	int32
	StartPoint		int32
	Default			string
}
