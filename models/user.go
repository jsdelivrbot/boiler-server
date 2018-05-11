package models

type User struct {
	MyUidObject

	Supervisor	*User		`orm:"rel(fk);null"`
	Role		*UserRole	`orm:"rel(fk)"`

	Organization	*Organization	`orm:"rel(fk);null"`

	Username	string		`orm:"size(60)"`
	Password	string

	Gender		int
	Picture		string
	Status		int

	Address		*Address	`orm:"rel(fk);null"`
	MyContact

	RegisterIp	string		`orm:"size(20)"`

	Thirds				[]*UserThird	`orm:"reverse(many)"`
	Sessions			[]*UserSession	`orm:"reverse(many)"`
	Logins				[]*UserLogin	`orm:"reverse(many)"`
	SubscribeBoilers	[]*Boiler		`orm:"reverse(many);rel_through(BoilerGo/models.BoilerMessageSubscriber)"`
}

const (
	USER_GENDER_UNKNOWN = 0
	USER_GENDER_MALE = 1
	USER_GENDER_FEMAIL = 2
)

const (
	USER_STATUS_NEW = -1
	USER_STATUS_INACTIVE = 0
	USER_STATUS_NORMAL = 1
	USER_STATUS_BANNED = 2
	USER_STATUS_THIRD = 3
)

func (usr *User) IsAdmin() bool {
	return usr.Role.RoleId <= USER_ROLE_SUPERADMIN
}

func (usr *User) IsOrganizationUser() bool {
	return usr.Role.RoleId == USER_ROLE_ORG_USER || usr.Role.RoleId == USER_ROLE_ORG_ADMIN
}

func (usr *User) IsAllowCreateOrganization() bool {
	return usr.Role.RoleId <= USER_ROLE_SUPERADMIN || usr.Role.RoleId == USER_ROLE_ORG_ADMIN
}

func (usr *User) IsCommonUser() bool {
	return usr.Role.RoleId == USER_ROLE_USER
}