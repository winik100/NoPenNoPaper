package main

import (
	"net/http"
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/models/mocks"
	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestUser(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name                  string
		authenticatedUserId   int
		authenticatedUserName string
		wantCode              int
		wantContent           []string
	}{
		{
			name:                  "Authenticated as Player",
			authenticatedUserId:   mocks.MockPlayer.ID,
			authenticatedUserName: mocks.MockPlayer.Name,
			wantContent:           []string{"<td><a href='/characters/1'>Otto Hightower</a></td>"},
			wantCode:              http.StatusOK,
		},
		{
			name:                  "Authenticated as GM",
			authenticatedUserId:   mocks.MockGM.ID,
			authenticatedUserName: mocks.MockGM.Name,
			wantContent:           []string{"<td><a href='/characters/1'>Otto Hightower</a></td>", "<td><a href='/characters/2'>Viserys Targaryen</a></td>"},
			wantCode:              http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockSession(app.authenticate(app.requireAuthentication(app.requireAuthorization(app.routesNoMW()))), map[string]any{
				authenticatedUserIdKey:   testCase.authenticatedUserId,
				authenticatedUserNameKey: testCase.authenticatedUserName,
			})))
			defer ts.Close()

			code, _, body := ts.get(t, "/users/"+testCase.authenticatedUserName)
			testHelpers.Equal(t, code, testCase.wantCode)
			for _, tag := range testCase.wantContent {
				testHelpers.StringContains(t, body, tag)
			}
		})
	}

}
