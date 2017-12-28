package models

type OrganizationType struct {
	MyUidObject

	TypeId			int32			`orm:"index"`
	Organizations		[]*Organization		`orm:"reverse(many)"`
}