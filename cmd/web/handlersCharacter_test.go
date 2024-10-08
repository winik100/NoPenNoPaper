package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/internal/models/mocks"
	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestCharacter(t *testing.T) {
	app := newTestApplication(t)
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
			characterId: "69",
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
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(noSurf(app.authenticate(app.requireAuthentication(app.routesNoMW()))),
				map[string]any{
					authenticatedUserIdKey:   1,
					authenticatedUserNameKey: "Testnutzer",
					characterIdKey:           testCase.characterId,
				})))
			defer ts.Close()

			code, _, body := ts.get(t, "/characters/"+testCase.characterId)

			testHelpers.Equal(t, code, testCase.wantCode)
			for _, tag := range testCase.wantContent {
				testHelpers.StringContains(t, body, tag)
			}
		})
	}
}

func TestCreateCharacterGet(t *testing.T) {
	app := newTestApplication(t)

	wantTag := "<form action='/create' method='POST'>"
	wantTagRedirect := "<a href='/login'>See Other</a>."
	wantContent := []string{
		"<div id='info'>",
		"<div id='attributes'>",
		"<div id='skills'>",
	}

	tests := []struct {
		name                  string
		isAuthenticated       bool
		authenticatedUserId   int
		authenticatedUserName string
		wantCode              int
		wantFormTag           string
		wantContent           []string
	}{
		{
			name:                  "Authenticated",
			isAuthenticated:       true,
			authenticatedUserId:   1,
			authenticatedUserName: "Testnutzer",
			wantCode:              http.StatusOK,
			wantFormTag:           wantTag,
			wantContent:           wantContent,
		},
		{
			name:                  "Unauthenticated",
			isAuthenticated:       false,
			authenticatedUserId:   0,
			authenticatedUserName: "",
			wantCode:              http.StatusSeeOther,
			wantFormTag:           wantTagRedirect,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(app.requireAuthentication(app.routesNoMW())), map[string]any{
				authenticatedUserIdKey:   testCase.authenticatedUserId,
				authenticatedUserNameKey: testCase.authenticatedUserName,
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

func TestCreateCharacterPost(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name                  string
		info                  core.CharacterInfo
		attributes            core.CharacterAttributes
		skills                core.Skills
		customSkills          core.CustomSkills
		authenticatedUserId   int
		authenticatedUserName string
		wantCode              int
	}{
		{
			name:                  "Valid Creation",
			info:                  mocks.MockCharacterOtto.Info,
			attributes:            mocks.MockCharacterOtto.Attributes,
			skills:                mocks.MockCharacterOtto.Skills,
			customSkills:          mocks.MockCharacterOtto.CustomSkills,
			authenticatedUserId:   1,
			authenticatedUserName: "Testnutzer",
			wantCode:              http.StatusSeeOther,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(noSurf(app.requireAuthentication(app.routesNoMW()))), map[string]any{
				authenticatedUserIdKey:   testCase.authenticatedUserId,
				authenticatedUserNameKey: testCase.authenticatedUserName,
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

func TestAddItem(t *testing.T) {
	app := newTestApplication(t)

	wantFormTag := "<form action='/characters/1/addItem' method='POST'>"
	wantContent := []string{
		`<form id="deleteItem" hx-post="/characters/1/deleteItem" hx-target="#item1" hx-swap="outerHTML">`,
		`<input type="hidden" name="ItemId" Value="1">`,
		`Hand-Brosche   <button type="submit">entfernen</button>`,
		"<td>Brosche der Hand des Königs</td>"}

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
			itemName:  mocks.MockCharacterOtto.Items.Name[0],
			itemDesc:  mocks.MockCharacterOtto.Items.Description[0],
			itemCount: strconv.Itoa(mocks.MockCharacterOtto.Items.Count[0]),
			redirect:  false,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:        "Valid Item, Content",
			itemName:    mocks.MockCharacterOtto.Items.Name[0],
			itemDesc:    mocks.MockCharacterOtto.Items.Description[0],
			itemCount:   strconv.Itoa(mocks.MockCharacterOtto.Items.Count[0]),
			redirect:    true,
			wantContent: wantContent,
		},
		{
			name:        "Empty Name",
			itemName:    "",
			itemDesc:    mocks.MockCharacterOtto.Items.Description[0],
			itemCount:   strconv.Itoa(mocks.MockCharacterOtto.Items.Count[0]),
			redirect:    false,
			wantCode:    http.StatusUnprocessableEntity,
			wantContent: []string{wantFormTag},
		},
		{
			name:        "Empty Description",
			itemName:    mocks.MockCharacterOtto.Items.Name[0],
			itemDesc:    "",
			itemCount:   strconv.Itoa(mocks.MockCharacterOtto.Items.Count[0]),
			redirect:    false,
			wantCode:    http.StatusUnprocessableEntity,
			wantContent: []string{wantFormTag},
		},
		{
			name:        "Name length > 50",
			itemName:    strings.Repeat(".", 51),
			itemDesc:    mocks.MockCharacterOtto.Items.Description[0],
			itemCount:   strconv.Itoa(mocks.MockCharacterOtto.Items.Count[0]),
			redirect:    false,
			wantCode:    http.StatusUnprocessableEntity,
			wantContent: []string{wantFormTag},
		},
		{
			name:        "Description length > 255",
			itemName:    mocks.MockCharacterOtto.Items.Name[0],
			itemDesc:    strings.Repeat(".", 256),
			itemCount:   strconv.Itoa(mocks.MockCharacterOtto.Items.Count[0]),
			redirect:    false,
			wantCode:    http.StatusUnprocessableEntity,
			wantContent: []string{wantFormTag},
		},
		{
			name:        "Item count < 1",
			itemName:    mocks.MockCharacterOtto.Items.Name[0],
			itemDesc:    strings.Repeat(".", 256),
			itemCount:   "0",
			redirect:    false,
			wantCode:    http.StatusUnprocessableEntity,
			wantContent: []string{wantFormTag},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(noSurf(app.authenticate(app.requireAuthentication(app.routesNoMW()))),
				map[string]any{
					authenticatedUserIdKey:   1,
					authenticatedUserNameKey: "Testnutzer",
					"characterId":            1,
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

func TestDeleteItem(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(noSurf(app.authenticate(app.requireAuthentication(app.routesNoMW()))),
		map[string]any{
			authenticatedUserIdKey:   1,
			authenticatedUserNameKey: "Testnutzer",
			"characterId":            1,
		})))
	defer ts.Close()
	_, _, body := ts.get(t, "/characters/1")
	testHelpers.StringContains(t, body, "Hand-Brosche")
	validCSRF := extractCSRFToken(t, body)

	tests := []struct {
		name     string
		itemId   int
		wantCode int
	}{
		{
			name:     "Successful Deletion",
			itemId:   1,
			wantCode: http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("ItemId", strconv.Itoa(testCase.itemId))
			form.Add("csrf_token", validCSRF)

			code, _, _ := ts.postForm(t, "/characters/1/deleteItem", form)
			testHelpers.Equal(t, code, testCase.wantCode)

			_, _, body := ts.get(t, "/characters/1")
			if strings.Contains(body, "Hand-Brosche") {
				t.Errorf("item was not removed")
			}
		})
	}
}

func TestAddNote(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(noSurf(app.authenticate(app.requireAuthentication(app.routesNoMW()))),
		map[string]any{
			authenticatedUserIdKey:   1,
			authenticatedUserNameKey: "Testnutzer",
		})))
	defer ts.Close()

	_, _, body := ts.get(t, "/characters/1/addNote")

	validCSRF := extractCSRFToken(t, body)

	wantContent := []string{
		`<form id="deleteNote" hx-post="/characters/1/deleteNote" hx-target="this" hx-swap="outerHTML">`,
		`<input type="hidden" name="NoteId" value="2">`,
		`<li>Dies ist eine gültige Notiz.    <button type="submit">löschen</button></li>`,
	}

	tests := []struct {
		name        string
		text        string
		wantCode    int
		wantContent []string
	}{
		{
			name:        "Valid Note",
			text:        "Dies ist eine gültige Notiz.",
			wantCode:    http.StatusOK,
			wantContent: wantContent,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("CharacterId", "1")
			form.Add("Text", testCase.text)
			form.Add("csrf_token", validCSRF)

			code, _, body := ts.postForm(t, "/characters/1/addNote", form)

			testHelpers.Equal(t, code, testCase.wantCode)
			for _, tag := range testCase.wantContent {
				testHelpers.StringContains(t, body, tag)
			}
		})
	}
}

func TestDeleteNote(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(noSurf(app.authenticate(app.requireAuthentication(app.routesNoMW()))),
		map[string]any{
			authenticatedUserIdKey:   1,
			authenticatedUserNameKey: "Testnutzer",
			"characterId":            1,
		})))
	defer ts.Close()
	_, _, body := ts.get(t, "/characters/1")
	testHelpers.StringContains(t, body, "Aegon ist blöde.")
	validCSRF := extractCSRFToken(t, body)

	tests := []struct {
		name     string
		noteId   int
		wantCode int
	}{
		{
			name:     "Successful Deletion",
			noteId:   1,
			wantCode: http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("NoteId", strconv.Itoa(testCase.noteId))
			form.Add("csrf_token", validCSRF)

			code, _, _ := ts.postForm(t, "/characters/1/deleteNote", form)
			testHelpers.Equal(t, code, testCase.wantCode)

			_, _, body := ts.get(t, "/characters/1")
			if strings.Contains(body, "Aegon ist blöde.") {
				t.Errorf("note was not removed")
			}
			testHelpers.StringContains(t, body, "Viserys war viel besser.")
		})
	}
}
