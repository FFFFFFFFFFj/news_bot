package models

type Role string

const (
	RoleUser       Role = "user"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "superadmin"
)

type User struct {
	ID       int64  `bson:"_id"`
	Username string `bson:"username"`
	Role     Role   `bson:"role"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}
