package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestIndex(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		userId   int
		userName string
		wantCode int
		location string
	}{
		{
			name:     "Authenticated",
			userId:   1,
			userName: "Testnutzer",
			wantCode: http.StatusPermanentRedirect,
			location: "/users/Testnutzer",
		},
		{
			name:     "Unauthenticated",
			userId:   0,
			userName: "",
			wantCode: http.StatusSeeOther,
			location: "/login",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(noSurf(app.requireAuthentication(app.routesNoMW()))), map[string]any{
				authenticatedUserIdKey:   testCase.userId,
				authenticatedUserNameKey: testCase.userName,
			})))
			defer ts.Close()

			req, err := http.NewRequest("GET", ts.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			response, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			testHelpers.Equal(t, response.StatusCode, testCase.wantCode)
			testHelpers.Equal(t, response.Header.Get("Location"), testCase.location)
		})
	}
}

func TestSignup(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(noSurf(app.routesNoMW())), map[string]any{
		authenticatedUserIdKey: 0,
	})))
	defer ts.Close()
	_, _, body := ts.get(t, "/signup")

	validCSRF := extractCSRFToken(t, body)

	const (
		validName     string = "Neuer Nutzer"
		validPassword string = "Klartext ole"
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
			name:         "Name Taken",
			userName:     "Testnutzer",
			userPassword: validPassword,
			csrfToken:    validCSRF,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  "<label class='error'>Dieser Name ist bereits vergeben.</label>",
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

func TestLogout(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name                  string
		authenticatedUserId   int
		authenticatedUserName string
		wantCode              int
		wantRedirect          string
	}{
		{
			name:                  "Authenticated",
			authenticatedUserId:   1,
			authenticatedUserName: "Testnutzer",
			wantCode:              http.StatusSeeOther,
			wantRedirect:          "/",
		},
		{
			name:                  "Unauthenticated",
			authenticatedUserId:   0,
			authenticatedUserName: "",
			wantCode:              http.StatusSeeOther,
			wantRedirect:          "/login",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(app.requireAuthentication(app.routesNoMW())), map[string]any{
				authenticatedUserIdKey:   testCase.authenticatedUserId,
				authenticatedUserNameKey: testCase.authenticatedUserName,
			})))
			defer ts.Close()

			req, err := http.NewRequest("POST", ts.URL+"/logout", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("session", "fictionalSessionId")
			response, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			testHelpers.Equal(t, response.StatusCode, http.StatusSeeOther)
			testHelpers.Equal(t, response.Header.Get("Location"), testCase.wantRedirect)
			if response.Header.Get("session") == "fictionalSessionId" {
				t.Errorf("session token was not renewed")
			}
		})
	}
}
