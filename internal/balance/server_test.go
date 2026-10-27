package balance

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func testServer(t *testing.T) (*Server, *Store, string, string) {
	t.Helper()
	store, e := Open(filepath.Join(t.TempDir(), "data.json"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { store.Close() })
	if e = store.Update(func(st *State) error { *st = fixture(); return nil }); e != nil {
		t.Fatal(e)
	}
	s, e := NewServer(store, Config{Now: func() time.Time { return testToday }})
	if e != nil {
		t.Fatal(e)
	}
	key, csrf := s.sessions.new("a", testToday)
	return s, store, key, csrf
}
func request(s *Server, method, path, key string, data url.Values) *httptest.ResponseRecorder {
	var body *strings.Reader
	if data == nil {
		body = strings.NewReader("")
	} else {
		body = strings.NewReader(data.Encode())
	}
	r := httptest.NewRequest(method, path, body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if key != "" {
		r.AddCookie(&http.Cookie{Name: "balance_session", Value: key})
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	return w
}
func TestAllScreens(t *testing.T) {
	s, store, key, _ := testServer(t)
	for _, path := range []string{"/overview", "/operations", "/categories", "/payments", "/analytics", "/settings", "/operations?edit=history", "/operations?delete=history", "/analytics?from=2026-10-01&to=2026-10-26&compare_from=2026-09-01&compare_to=2026-09-26"} {
		t.Run(path, func(t *testing.T) {
			w := request(s, "GET", path, key, nil)
			if w.Code != 200 || !strings.Contains(w.Body.String(), "</html>") {
				t.Fatalf("status %d: %s", w.Code, w.Body.String())
			}
		})
	}
	if e := store.Update(func(st *State) error { st.Users[0].Started = ""; return nil }); e != nil {
		t.Fatal(e)
	}
	for _, path := range []string{"/overview", "/operations", "/categories", "/payments", "/analytics", "/settings"} {
		w := request(s, "GET", path, key, nil)
		if w.Code != 200 {
			t.Fatalf("empty %s: %d %s", path, w.Code, w.Body.String())
		}
	}
	for _, path := range []string{"/login", "/register", "/static/style.css", "/static/app.js", "/healthz"} {
		if w := request(s, "GET", path, "", nil); w.Code != 200 {
			t.Fatalf("%s %d", path, w.Code)
		}
	}
}
func TestCSRFIsolationAndEscaping(t *testing.T) {
	s, store, key, csrf := testServer(t)
	w := request(s, "POST", "/operation", key, url.Values{"kind": {"expense"}, "amount": {"5"}})
	if w.Code != 403 {
		t.Fatal("missing csrf", w.Code)
	}
	data := url.Values{"csrf": {csrf}, "token": {"form-1"}, "kind": {"expense"}, "amount": {"5,25"}, "date": {"2026-10-26"}, "category": {"food"}, "comment": {"<script>alert(1)</script>"}}
	w = request(s, "POST", "/operation", key, data)
	if w.Code != 303 {
		t.Fatalf("save %d %s", w.Code, w.Body.String())
	}
	w = request(s, "GET", "/operations", key, nil)
	if strings.Contains(w.Body.String(), "<script>alert") || !strings.Contains(w.Body.String(), "&lt;script&gt;") {
		t.Fatal("XSS")
	}
	if e := store.Update(func(st *State) error {
		st.Operations = append(st.Operations, Operation{ID: "secret", UserID: "b", CategoryID: "foreign", Kind: "expense", Date: "2026-10-26", Amount: 777})
		return nil
	}); e != nil {
		t.Fatal(e)
	}
	data.Set("id", "secret")
	data.Set("amount", "1")
	data.Set("category", "food")
	if w = request(s, "POST", "/operation", key, data); w.Code != 422 {
		t.Fatal("foreign edit", w.Code)
	}
	st := store.Snapshot()
	if st.Operation("b", "secret").Amount != 777 {
		t.Fatal("foreign record changed")
	}
	data.Del("id")
	data.Set("token", "form-2")
	data.Set("amount", "0")
	w = request(s, "POST", "/operation", key, data)
	if w.Code != 422 || !strings.Contains(w.Body.String(), "&lt;script&gt;") {
		t.Fatal("error must preserve form", w.Code)
	}
	r := httptest.NewRequest("POST", "/operation", strings.NewReader(data.Encode()))
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	w = httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("cross origin")
	}
	if w = request(s, "GET", "/overview", "", nil); w.Code != 303 {
		t.Fatal("auth required")
	}
	if w = request(s, "GET", "/operations?edit=secret", key, nil); w.Code != 422 {
		t.Fatal("foreign get", w.Code)
	}
	w = request(s, "POST", "/logout", key, url.Values{"csrf": {csrf}})
	if w.Code != 303 {
		t.Fatal("logout")
	}
	if w = request(s, "GET", "/overview", key, nil); w.Code != 303 {
		t.Fatal("session after logout")
	}
}
func TestSessionExpiryAndCookie(t *testing.T) {
	s, _, _, _ := testServer(t)
	s.config.Secure = true
	w := request(s, "GET", "/login", "", nil)
	cookies := w.Result().Cookies()
	if len(cookies) == 0 || !cookies[0].Secure || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteStrictMode {
		t.Fatal("cookie flags")
	}
	key, _ := s.sessions.new("a", testToday.Add(-31*time.Minute))
	if w = request(s, "GET", "/overview", key, nil); w.Code != 303 {
		t.Fatal("expiry")
	}
}
func TestPaymentHTTPReplayAndLimits(t *testing.T) {
	s, store, key, csrf := testServer(t)
	values := url.Values{"csrf": {csrf}, "id": {"oct"}, "action": {"complete"}, "date": {"2026-10-26"}, "amount": {"100"}}
	for i := 0; i < 2; i++ {
		if w := request(s, "POST", "/payment", key, values); w.Code != 303 {
			t.Fatal(w.Code, w.Body.String())
		}
	}
	if len(store.Snapshot().Operations) != 2 {
		t.Fatal("duplicate")
	}
	if w := request(s, "POST", "/limit", key, url.Values{"csrf": {csrf}, "category": {"food"}, "month": {"2026-07"}, "amount": {"920"}}); w.Code != 303 {
		t.Fatal(w.Code)
	}
	w := request(s, "GET", "/categories?month=2026-07", key, nil)
	if !strings.Contains(w.Body.String(), "Лимит достигнут") {
		t.Fatal("equality limit")
	}
}
func ExampleParseMoney() {
	m, _ := ParseMoney("12,50")
	fmt.Println(m.String()) // Output: 12,50
}

func TestHTTPSRegistrationAndSession(t *testing.T) {
	app, _, _, _ := testServer(t)
	app.config.Secure = true
	ts := httptest.NewTLSServer(app)
	defer ts.Close()
	client := ts.Client()
	jar, e := cookiejar.New(nil)
	if e != nil {
		t.Fatal(e)
	}
	client.Jar = jar
	response, e := client.Get(ts.URL + "/register")
	if e != nil {
		t.Fatal(e)
	}
	body, e := io.ReadAll(response.Body)
	response.Body.Close()
	if e != nil {
		t.Fatal(e)
	}
	token := regexp.MustCompile(`name="csrf" value="([^"]+)"`).FindSubmatch(body)
	if len(token) != 2 {
		t.Fatal("registration csrf not rendered")
	}
	response, e = client.PostForm(ts.URL+"/register", url.Values{"csrf": {string(token[1])}, "email": {"tls@example.by"}, "password": {"tls-password-12345"}})
	if e != nil {
		t.Fatal(e)
	}
	body, e = io.ReadAll(response.Body)
	response.Body.Close()
	if e != nil {
		t.Fatal(e)
	}
	if response.StatusCode != 200 || response.Request.URL.Path != "/overview" || !strings.Contains(string(body), "tls@example.by") {
		t.Fatal("secure session registration", response.StatusCode)
	}
	u, _ := url.Parse(ts.URL)
	found := false
	for _, c := range jar.Cookies(u) {
		if c.Name == "balance_session" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing session cookie")
	}
	token = regexp.MustCompile(`name="csrf" value="([^"]+)"`).FindSubmatch(body)
	if len(token) != 2 {
		t.Fatal("session csrf not rendered")
	}
	response, e = client.PostForm(ts.URL+"/logout", url.Values{"csrf": {string(token[1])}})
	if e != nil {
		t.Fatal(e)
	}
	response.Body.Close()
	if response.Request.URL.Path != "/login" {
		t.Fatal("secure logout")
	}
}
