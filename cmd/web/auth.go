package main

import (
	"net/http"
	"regexp"
	"slices"

	"github.com/winik100/NoPenNoPaper/internal/core"
)

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated := app.sessionManager.GetBool(r.Context(), isAuthenticatedKey)
	return isAuthenticated
}

func (app *application) isAuthorized(r *http.Request) bool {
	role := app.sessionManager.GetString(r.Context(), roleKey)
	path := r.URL.Path

	requestedUserName := r.PathValue("name")

	if requestedUserName != "" {
		if role != core.RoleGM {
			userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)
			if userName != requestedUserName {
				return false
			}
		}
	}

	return permitted(role, path)
}

func permitted(role string, path string) bool {
	for key, perms := range permissions {
		exp := regexp.MustCompile(key)
		if exp.MatchString(path) {
			return slices.Contains(perms, role)
		}
	}
	return false
}

var permissions = map[string][]string{
	"/":                   {core.RoleAnon, core.RolePlayer, core.RoleGM},
	"/signup":             {core.RoleAnon, core.RolePlayer, core.RoleGM},
	"/login":              {core.RoleAnon, core.RolePlayer, core.RoleGM},
	"/logout":             {core.RolePlayer, core.RoleGM},
	"/characters/\\d+/.*": {core.RolePlayer, core.RoleGM},
	"/users/*/.*":         {core.RolePlayer, core.RoleGM},
}
