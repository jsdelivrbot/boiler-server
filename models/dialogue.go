package models

type Dialogue struct {
	MyUidObject

	DialogueId		int64			`orm:"index"`
	Status			int			`orm:"index"`
	Comments		[]*DialogueComment	`orm:"reverse(many)"`
}

const (
	DIALOGUE_STATUS_DEFAULT = 0
	DIALOGUE_STATUS_NEW = 1
	DIALOGUE_STATUS_PENDING = 2
	DIALOGUE_STATUS_ANSWERED = 3
	DIALOGUE_STATUS_CLOSED = 4
)