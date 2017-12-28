package models

type Organization struct {
	MyUidObject

	Type				*OrganizationType	`orm:"rel(fk);index"`

	IsSupervisor		bool
	SuperOrganization	*Organization		`orm:"rel(fk);null;index"`
	SubOrganizations	[]*Organization		`orm:"reverse(many)"`

	Address				*Address			`orm:"rel(fk);null"`
	Contact				*Contact			`orm:"rel(fk);null"`

	ShowBrand			bool
	BrandName			string				`orm:"size(20)"`
	BrandImageUrl		string

	Users				[]*User				`orm:"reverse(many);null"`
	Boilers				[]*Boiler			`orm:"reverse(many);null"`
}
