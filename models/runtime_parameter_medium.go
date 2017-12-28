package models

type RuntimeParameterMedium struct {
	MyIdObject

	Parameters	[]*RuntimeParameter		`orm:"reverse(many)"`
}
