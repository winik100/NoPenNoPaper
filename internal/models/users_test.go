package models

import (
	"errors"
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestInsert(t *testing.T) {
	db := newTestDB(t)

	u := UserModel{db}

	tests := []struct {
		name          string
		playerName    string
		plainPassword string
		wantId        int
	}{
		{
			name:          "Valid Insert",
			playerName:    "test",
			plainPassword: "testpw",
			wantId:        2,
		},
		{
			name:          "Already exists",
			playerName:    "testgm",
			plainPassword: "testpwgm",
			wantId:        0,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			userId, err := u.Insert(testCase.playerName, testCase.plainPassword)
			if err != nil {
				if !errors.Is(err, ErrNameTaken) {
					t.Fatal(err)
				}
			}
			testHelpers.Equal(t, userId, testCase.wantId)
		})
	}
}

func TestGet(t *testing.T) {
	db := newTestDB(t)

	u := UserModel{db}

	tests := []struct {
		name          string
		playerName    string
		plainPassword string
		wantUser      core.User
	}{
		{
			name:          "Valid User",
			playerName:    "testgm",
			plainPassword: "testpwgm",
			wantUser:      core.User{ID: 1, Name: "testgm", Role: "gm"},
		},
		{
			name:          "Nonexistent User",
			playerName:    "test",
			plainPassword: "testpw",
			wantUser:      core.User{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			user, err := u.Get(testCase.playerName)
			if err != nil && !errors.Is(err, ErrNoRecord) {
				t.Fatal(err)
			}
			testHelpers.Equal(t, user.ID, testCase.wantUser.ID)
			testHelpers.Equal(t, user.Role, testCase.wantUser.Role)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	db := newTestDB(t)

	u := UserModel{db}

	tests := []struct {
		name          string
		playerName    string
		plainPassword string
		wantId        int
	}{
		{
			name:          "Valid Authentication",
			playerName:    "testgm",
			plainPassword: "testpwgm",
			wantId:        1,
		},
		{
			name:          "Invalid Authentication",
			playerName:    "test",
			plainPassword: "testpw",
			wantId:        0,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			userId, err := u.Authenticate(testCase.playerName, testCase.plainPassword)
			if err != nil {
				if !errors.Is(err, ErrInvalidCredentials) {
					t.Fatal(err)
				}
			}
			testHelpers.Equal(t, userId, testCase.wantId)
		})
	}
}

func TestExists(t *testing.T) {
	db := newTestDB(t)

	u := UserModel{db}

	tests := []struct {
		name          string
		playerName    string
		plainPassword string
		want          bool
	}{
		{
			name:       "Existing User",
			playerName: "testgm",
			want:       true,
		},
		{
			name:       "Nonexistent User",
			playerName: "test",
			want:       false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			exists, err := u.Exists(testCase.playerName)
			if err != nil {
				t.Fatal(err)
			}
			testHelpers.Equal(t, exists, testCase.want)
		})
	}
}
