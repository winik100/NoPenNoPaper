package main

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/models/mocks"
	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestHome(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name                string
		authenticatedUserId int
		wantCode            int
		wantTag             string
	}{
		{
			name:                "Authenticated",
			authenticatedUserId: 1,
			wantTag:             "<td><a href='/characters/1'>Otto Hightower</a></td>",
			wantCode:            http.StatusOK,
		},
		{
			name:                "Unauthenticated",
			authenticatedUserId: 0,
			wantTag:             "<p>Um Charaktere zu erstellen oder einzusehen, bitte einloggen.</p>",
			wantCode:            http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(app.routesNoMW()), map[string]any{
				string(authenticatedUserIdContextKey): testCase.authenticatedUserId,
			})))
			defer ts.Close()

			code, _, body := ts.get(t, "/")
			testHelpers.Equal(t, code, testCase.wantCode)
			testHelpers.StringContains(t, body, testCase.wantTag)
		})
	}

}

func TestSignup(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(app.routesNoMW()), map[string]any{
		string(authenticatedUserIdContextKey): 0,
	})))
	defer ts.Close()
	_, _, body := ts.get(t, "/signup")

	validCSRF := extractCSRFToken(t, body)

	const (
		validName     string = "Testnutzer"
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
		name                string
		authenticatedUserId int
		wantCode            int
		wantRedirect        string
	}{
		{
			name:                "Authenticated",
			authenticatedUserId: 1,
			wantCode:            http.StatusSeeOther,
			wantRedirect:        "/",
		},
		{
			name:                "Unauthenticated",
			authenticatedUserId: 0,
			wantCode:            http.StatusSeeOther,
			wantRedirect:        "/login",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(app.requireAuthentication(app.routesNoMW())), map[string]any{
				string(authenticatedUserIdContextKey): testCase.authenticatedUserId,
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

func TestCreateGet(t *testing.T) {
	app := newTestApplication(t)

	wantTag := "<form action='/create' method='POST'>"
	wantTagRedirect := "<a href='/login'>See Other</a>."
	wantContent := []string{
		"<div id='info'>",
		"<div id='attributes'>",
		"<div id='skills'>",
	}

	tests := []struct {
		name                string
		isAuthenticated     bool
		authenticatedUserId int
		wantCode            int
		wantFormTag         string
		wantContent         []string
	}{
		{
			name:                "Authenticated",
			isAuthenticated:     true,
			authenticatedUserId: 1,
			wantCode:            http.StatusOK,
			wantFormTag:         wantTag,
			wantContent:         wantContent,
		},
		{
			name:                "Unauthenticated",
			isAuthenticated:     false,
			authenticatedUserId: 0,
			wantCode:            http.StatusSeeOther,
			wantFormTag:         wantTagRedirect,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(app.requireAuthentication(app.routesNoMW())), map[string]any{
				string(authenticatedUserIdContextKey): testCase.authenticatedUserId,
			})))
			defer ts.Close()

			code, _, body := ts.get(t, "/create")

			testHelpers.Equal(t, code, testCase.wantCode)

			if testCase.isAuthenticated {
				testHelpers.StringContains(t, body, testCase.wantFormTag)
				for _, tag := range testCase.wantContent {
					testHelpers.StringContains(t, body, tag)
				}
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name                string
		info                models.CharacterInfo
		attributes          models.CharacterAttributes
		skills              models.Skills
		customSkills        models.CustomSkills
		authenticatedUserId int
		wantCode            int
	}{
		{
			name:                "Valid Creation",
			info:                mocks.MockCharacter.Info,
			attributes:          mocks.MockCharacter.Attributes,
			skills:              mocks.MockCharacter.Skills,
			customSkills:        mocks.MockCharacter.CustomSkills,
			authenticatedUserId: 1,
			wantCode:            http.StatusSeeOther,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(noSurf(app.requireAuthentication(app.routesNoMW()))), map[string]any{
				string(authenticatedUserIdContextKey): testCase.authenticatedUserId,
			})))
			defer ts.Close()
			_, _, body := ts.get(t, "/create")

			validCSRF := extractCSRFToken(t, body)

			form := url.Values{}
			form.Add("Info.Name", testCase.info.Name)
			form.Add("Info.Profession", testCase.info.Profession)
			form.Add("Info.Age", testCase.info.Age)
			form.Add("Info.Gender", testCase.info.Gender)
			form.Add("Info.Residence", testCase.info.Residence)
			form.Add("Info.Birthplace", testCase.info.Birthplace)

			form.Add("Attributes.ST", strconv.Itoa(testCase.attributes.ST))
			form.Add("Attributes.GE", strconv.Itoa(testCase.attributes.GE))
			form.Add("Attributes.MA", strconv.Itoa(testCase.attributes.MA))
			form.Add("Attributes.KO", strconv.Itoa(testCase.attributes.KO))
			form.Add("Attributes.ER", strconv.Itoa(testCase.attributes.ER))
			form.Add("Attributes.BI", strconv.Itoa(testCase.attributes.BI))
			form.Add("Attributes.GR", strconv.Itoa(testCase.attributes.GR))
			form.Add("Attributes.IN", strconv.Itoa(testCase.attributes.IN))
			form.Add("Attributes.BW", strconv.Itoa(testCase.attributes.BW))

			for i, skill := range testCase.skills.Name {
				form.Add("Skills.Name", skill)
				form.Add("Skills.Value", strconv.Itoa(testCase.skills.Value[i]))
			}

			for i, customSkill := range testCase.customSkills.Name {
				form.Add("CustomSkills.Name", customSkill)
				form.Add("CustomSkills.Value", strconv.Itoa(testCase.customSkills.Value[i]))
			}
			form.Add("csrf_token", validCSRF)

			code, header, _ := ts.postForm(t, "/create", form)

			testHelpers.Equal(t, code, testCase.wantCode)
			testHelpers.Equal(t, header.Get("Location"), "/")
		})
	}
}

func TestViewCharacter(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(noSurf(app.authenticate(app.requireAuthentication(app.routesNoMW()))),
		map[string]any{
			string(authenticatedUserIdContextKey): 1,
		})))
	defer ts.Close()
	wantContent := []string{
		"<div id='info'>",
		"<div id='attributes'>",
		"<div id='skills'>",
		"<div id='items'>",
		"<div id='notes'>",
		"Otto Hightower",
	}

	tests := []struct {
		name        string
		characterId string
		wantCode    int
		wantContent []string
	}{
		{
			name:        "Valid ID",
			characterId: "1",
			wantCode:    http.StatusOK,
			wantContent: wantContent,
		},
		{
			name:        "Nonexistent, valid ID",
			characterId: "2",
			wantCode:    http.StatusNotFound,
		},
		{
			name:        "Empty ID",
			characterId: "",
			wantCode:    http.StatusNotFound,
		},
		{
			name:        "Invalid ID",
			characterId: "test",
			wantCode:    http.StatusNotFound,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {

			code, _, body := ts.get(t, "/characters/"+testCase.characterId)

			testHelpers.Equal(t, code, testCase.wantCode)
			for _, tag := range testCase.wantContent {
				testHelpers.StringContains(t, body, tag)
			}
		})
	}
}

func TestAddItem(t *testing.T) {
	app := newTestApplication(t)

	wantFormTag := "<form action='/characters/1/addItem' method='POST'>"
	wantContent := []string{
		"<td>Hand-Brosche</td>",
		"<td>Brosche der Hand des Königs</td>",
		"<td>1</td>"}

	tests := []struct {
		name        string
		itemName    string
		itemDesc    string
		itemCount   string
		redirect    bool
		wantCode    int
		wantContent []string
	}{
		{
			name:      "Valid Item, Status",
			itemName:  mocks.MockCharacter.Items[0].Name,
			itemDesc:  mocks.MockCharacter.Items[0].Description,
			itemCount: strconv.Itoa(mocks.MockCharacter.Items[0].Count),
			redirect:  false,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:        "Valid Item, Content",
			itemName:    mocks.MockCharacter.Items[0].Name,
			itemDesc:    mocks.MockCharacter.Items[0].Description,
			itemCount:   strconv.Itoa(mocks.MockCharacter.Items[0].Count),
			redirect:    true,
			wantContent: wantContent,
		},
		{
			name:        "Invalid Item",
			itemName:    "",
			itemDesc:    "",
			itemCount:   "",
			redirect:    false,
			wantCode:    http.StatusUnprocessableEntity,
			wantContent: []string{wantFormTag},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(noSurf(app.authenticate(app.requireAuthentication(app.routesNoMW()))),
				map[string]any{
					string(authenticatedUserIdContextKey): 1,
					"characterId":                         1,
				})))
			defer ts.Close()
			_, _, body := ts.get(t, "/characters/1/addItem")

			validCSRF := extractCSRFToken(t, body)

			form := url.Values{}
			form.Add("CharacterId", "1")
			form.Add("Name", testCase.itemName)
			form.Add("Description", testCase.itemDesc)
			form.Add("Count", testCase.itemCount)
			form.Add("csrf_token", validCSRF)

			if testCase.redirect {
				ts.Client().CheckRedirect = nil
			}
			code, header, body := ts.postForm(t, "/characters/1/addItem", form)

			if !testCase.redirect {
				testHelpers.Equal(t, code, testCase.wantCode)
				if testCase.wantCode == http.StatusSeeOther {
					testHelpers.Equal(t, header.Get("Location"), "/characters/1")
				}
			}
			for _, tag := range testCase.wantContent {
				testHelpers.StringContains(t, body, tag)
			}
		})
	}
}
