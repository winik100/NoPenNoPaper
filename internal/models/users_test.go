package models

import (
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestInsert(t *testing.T) {
	t.SkipNow()
	db := newTestDB(t)

	u := UserModel{db}

	userId, err := u.Insert("test", "testpw")
	if err != nil {
		t.Fatal(err)
	}

	testHelpers.Equal(t, userId, 1)
}
