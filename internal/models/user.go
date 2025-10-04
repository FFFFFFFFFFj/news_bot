package models

type Role string

const (
	RoleUser       Role = "user"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "superadmin"
)

//Helper functions for checking roles by string
func IsAdmin(role string) bool {
	return role == string(RoleAdmin) || role == string(RoleSuperAdmin)
}

func IsSuperAdmin(role string) bool {
	return role == string(RoleSuperAdmin)
}
