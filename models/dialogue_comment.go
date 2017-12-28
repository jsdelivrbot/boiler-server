package models

type DialogueComment struct {
	MyUidObject

	Dialogue	*Dialogue	`orm:"rel(fk);index"`
	From		*User		`orm:"rel(fk);index"`
	To		*User 		`orm:"rel(fk);index;null"`

	Content		string
	Attachment	string
}
