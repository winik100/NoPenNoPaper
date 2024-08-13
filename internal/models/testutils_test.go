package models

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "test:testpw@/test_nopennopaper?multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	path, err := filepath.Abs("../sql/testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	script, err := os.ReadFile(path)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	//teardown
	t.Cleanup(func() {
		path, err = filepath.Abs("../sql/testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		script, err := os.ReadFile(path)
		if err != nil {
			db.Close()
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			db.Close()
			t.Fatal(err)
		}
	})

	return db
}
