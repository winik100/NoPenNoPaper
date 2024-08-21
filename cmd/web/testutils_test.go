package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/schema"
	"github.com/winik100/NoPenNoPaper/internal/models/mocks"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := schema.NewDecoder()
	formDecoder.IgnoreUnknownKeys(true)

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		log:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		characters:     &mocks.CharacterModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

func (app *application) mockSession(next http.Handler, keyValuePairs map[string]any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range keyValuePairs {
			app.sessionManager.Put(r.Context(), key, value)
		}
		next.ServeHTTP(w, r)
	})
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, path string) (int, http.Header, string) {
	response, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return response.StatusCode, response.Header, string(body)
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)

	if len(matches) < 2 {
		t.Fatalf("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}
