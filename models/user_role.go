package models

type UserRole struct {
	MyUidObject

	RoleId		int32		`orm:"index"`
}

const (
	USER_ROLE_SYSTEM = 0
	USER_ROLE_SUPERADMIN = 1
	USER_ROLE_ADMIN = 2
	USER_ROLE_SERVICE = 3
	USER_ROLE_SUPERVISOR = 4
	USER_ROLE_ORG_ADMIN = 10
	USER_ROLE_ORG_USER = 11
	USER_ROLE_USER = 20
)