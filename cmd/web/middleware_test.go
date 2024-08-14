package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/winik100/NoPenNoPaper/internal/testHelpers"
)

func TestHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	headers(mockNext).ServeHTTP(rec, r)

	result := rec.Result()

	want := "origin-when-cross-origin"
	actual := result.Header.Get("Referrer-Policy")
	testHelpers.Equal(t, actual, want)

	want = "nosniff"
	actual = result.Header.Get("X-Content-Type-Options")
	testHelpers.Equal(t, actual, want)

	want = "deny"
	actual = result.Header.Get("X-Frame-Options")
	testHelpers.Equal(t, actual, want)

	want = "0"
	actual = result.Header.Get("X-XSS-Protection")
	testHelpers.Equal(t, actual, want)

}

func TestNoSurf(t *testing.T) {
	rec := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	noSurf(mockNext).ServeHTTP(rec, r)

	result := rec.Result()

	cookies := result.Cookies()
	actual := cookies[0]
	testHelpers.Equal(t, actual.HttpOnly, true)
	testHelpers.Equal(t, actual.Secure, true)
	testHelpers.Equal(t, actual.Path, "/")
}
