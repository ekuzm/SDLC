package balance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"
)

// Opt-in measurement on the developer's machine, not a substitute for the SRS
// control environment (2 vCPU / 2 GB). Includes rendering, snapshot and HTTP handling.
func TestLocalPerformance(t *testing.T) {
	if os.Getenv("BALANCE_PERF") != "1" {
		t.Skip("BALANCE_PERF=1 enables a 1000-request local measurement")
	}
	app, store, _, _ := testServer(t)
	keys := make([]string, 10)
	if e := store.Update(func(st *State) error {
		*st = State{Version: 1}
		for u := 0; u < 10; u++ {
			uid := fmt.Sprint(u)
			st.Users = append(st.Users, User{ID: uid, Started: "2026-06-01"})
			st.Categories = append(st.Categories, Category{ID: uid + "-food", UserID: uid, Kind: "expense", Name: "Продукты"})
			for i := 0; i < 10000; i++ {
				st.Operations = append(st.Operations, Operation{ID: fmt.Sprintf("%s-%d", uid, i), UserID: uid, CategoryID: uid + "-food", Kind: "expense", Date: fmt.Sprintf("2026-%02d-%02d", 7+i%4, 1+i%26), Amount: 100, Created: testToday, Token: fmt.Sprint(i)})
			}
			for i := 0; i < 100; i++ {
				st.Rules = append(st.Rules, Rule{ID: fmt.Sprintf("r%s-%d", uid, i), UserID: uid, CategoryID: uid + "-food", Name: "Платёж", Kind: "expense", Day: 28, Start: "2026-07-01", Amount: 100, Active: true})
			}
			st.Materialize(uid, testToday)
		}
		return nil
	}); e != nil {
		t.Fatal(e)
	}
	for i := range keys {
		keys[i], _ = app.sessions.new(fmt.Sprint(i), testToday)
	}
	for _, path := range []string{"/operations", "/analytics"} {
		for i := 0; i < 10; i++ {
			if w := request(app, "GET", path, keys[i], nil); w.Code != 200 {
				t.Fatal(w.Code)
			}
		}
		timings := make(chan time.Duration, 1000)
		var wg sync.WaitGroup
		for _, key := range keys {
			wg.Go(func() {
				for i := 0; i < 100; i++ {
					start := time.Now()
					req := httptest.NewRequest("GET", path, nil)
					req.AddCookie(&http.Cookie{Name: "balance_session", Value: key})
					w := httptest.NewRecorder()
					app.ServeHTTP(w, req)
					timings <- time.Since(start)
					if w.Code != 200 {
						t.Errorf("status %d", w.Code)
					}
				}
			})
		}
		wg.Wait()
		close(timings)
		var rows []time.Duration
		for elapsed := range timings {
			rows = append(rows, elapsed)
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i] < rows[j] })
		p95 := rows[949]
		t.Logf("%s p95=%s, max=%s; 1000 requests / 10 concurrent / 10000 operations per user / 100 rules per user; Go=%s CPUs=%d", path, p95, rows[999], runtime.Version(), runtime.NumCPU())
		target := 2 * time.Second
		if path == "/analytics" {
			target = time.Second
		}
		if p95 > target {
			t.Errorf("local p95 above target: %v", p95)
		}
	}
}
