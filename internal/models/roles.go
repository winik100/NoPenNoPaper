package models

const (
	RoleAnon   string = "anonymous"
	RolePlayer string = "player"
	RoleGM     string = "gm"
)

var Permissions = map[string][]string{
	RoleAnon:   {"/", "/signup", "/login"},
	RolePlayer: {"/", "/logout", "/create", "/player"},
	RoleGM:     {"/", "/logout", "/create", "/player"},
}
