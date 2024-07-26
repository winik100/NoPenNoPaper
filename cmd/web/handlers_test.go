package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestHomeLoggedIn(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockAuthentication(app.authenticate(app.routesNoMW()), 1, true, "player")))
	defer ts.Close()
	wantTag := "<td><a href='/characters/1'>Otto Hightower</a></td>"

	code, _, body := ts.get(t, "/")
	testHelpers.Equal(t, code, http.StatusOK)
	testHelpers.StringContains(t, body, wantTag)

}

func TestHomeNotLoggedIn(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockAuthentication(app.authenticate(app.routesNoMW()), 0, false, "anonymous")))
	defer ts.Close()
	wantTag := "<p>Um Charaktere zu erstellen oder einzusehen, bitte einloggen.</p>"

	code, _, body := ts.get(t, "/")
	testHelpers.Equal(t, code, http.StatusOK)
	testHelpers.StringContains(t, body, wantTag)

}

func TestSignup(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	_, _, body := ts.get(t, "/signup")

	validCSRF := extractCSRFToken(t, body)

	const (
		validName     string = "Testnutzer"
		validPassword string = "Klartext"
		formTag       string = "<form action='/signup' method='POST' novalidate>"
	)
	tests := []struct {
		name         string
		userName     string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid Signup",
			userName:     validName,
			userPassword: validPassword,
			csrfToken:    validCSRF,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF",
			userName:     validName,
			userPassword: validPassword,
			csrfToken:    "",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userPassword: validPassword,
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Password",
			userName:     validName,
			userPassword: "",
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Password less than 8 characters",
			userName:     validName,
			userPassword: "test",
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("Name", testCase.userName)
			form.Add("Password", testCase.userPassword)
			form.Add("csrf_token", testCase.csrfToken)

			code, _, body := ts.postForm(t, "/signup", form)

			testHelpers.Equal(t, code, testCase.wantCode)

			if testCase.wantFormTag != "" {
				testHelpers.StringContains(t, body, testCase.wantFormTag)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	_, _, body := ts.get(t, "/login")

	validCSRF := extractCSRFToken(t, body)

	const (
		validName     string = "Testnutzer"
		validPassword string = "Klartext ole"
		formTag       string = "<form action='/login' method='POST' novalidate>"
	)
	tests := []struct {
		name         string
		userName     string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid login",
			userName:     validName,
			userPassword: validPassword,
			csrfToken:    validCSRF,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Wrong Name",
			userName:     "wrongname",
			userPassword: validPassword,
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Wrong Password",
			userName:     validName,
			userPassword: "wrongpassword",
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid CSRF",
			userName:     validName,
			userPassword: validPassword,
			csrfToken:    "",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userPassword: validPassword,
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Password",
			userName:     validName,
			userPassword: "",
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("Name", testCase.userName)
			form.Add("Password", testCase.userPassword)
			form.Add("csrf_token", testCase.csrfToken)

			code, _, _ := ts.postForm(t, "/login", form)

			testHelpers.Equal(t, code, testCase.wantCode)

			if testCase.wantFormTag != "" {
				testHelpers.StringContains(t, body, testCase.wantFormTag)
			}
		})
	}
}
