package balance

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func mustDate(s string) time.Time {
	t, e := Date(s)
	if e != nil {
		panic(e)
	}
	return t
}

var testToday = mustDate("2026-10-26")

func fixture() State {
	return State{Version: 1, Users: []User{{ID: "a", Started: "2026-07-01", Initial: 192000}, {ID: "b", Started: "2026-07-01"}}, Categories: []Category{{ID: "food", UserID: "a", Name: "Продукты", Kind: "expense"}, {ID: "salary", UserID: "a", Name: "Зарплата", Kind: "income"}, {ID: "foreign", UserID: "b", Name: "Продукты", Kind: "expense"}}, Operations: []Operation{{ID: "history", UserID: "a", CategoryID: "food", Kind: "expense", Date: "2026-07-15", Amount: 92000, Token: "history", Created: testToday}}, Occurrences: []Occurrence{{ID: "oct", UserID: "a", CategoryID: "food", Kind: "expense", Name: "Интернет", Date: "2026-10-30", Amount: 10000, Status: "pending"}, {ID: "nov", UserID: "a", CategoryID: "food", Kind: "expense", Name: "Интернет", Date: "2026-11-30", Amount: 10000, Status: "pending"}}}
}
func TestMoney(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want Money
	}{{"0,01", 1}, {"999999999.99", MaxMoney}, {"-0,01", -1}, {"1.1", 110}, {" 42 ", 4200}} {
		m, e := ParseMoney(tc.in)
		if e != nil || m != tc.want {
			t.Fatalf("%q=%d, %v", tc.in, m, e)
		}
	}
	for _, in := range []string{"", "1.001", "1e3", "NaN", "+1", "--1", "1.", "1000000000", "1,2,3", ".01"} {
		if _, e := ParseMoney(in); e == nil {
			t.Errorf("accepted %q", in)
		}
	}
	if Money(-1).String() != "-0,01" {
		t.Fatal("format")
	}
}
func TestForecastExamples(t *testing.T) {
	s := fixture()
	f := s.Forecast("a", testToday)
	if !f.Available || f.EndBalance != 85000 || f.NextExpense != 40000 || f.HistoryDays != 92 || f.Months != 3 || f.RemainingDays != 5 {
		t.Fatalf("F01/F02: %+v", f)
	}
	s.User("a").Started = "2026-10-15"
	if s.Forecast("a", testToday).Available {
		t.Fatal("F03 missing history")
	}
	s = fixture()
	s.User("a").Started = "2026-09-01"
	s.User("a").Initial = 100000
	s.Operations = nil
	f = s.Forecast("a", testToday)
	if f.EndBalance != 90000 || f.HistoryDays != 30 {
		t.Fatalf("F04 %+v", f)
	}
	s = fixture()
	f = s.Forecast("a", mustDate("2026-10-31"))
	if f.RemainingDays != 0 || f.EndBalance != 90000 {
		t.Fatalf("F05 %+v", f)
	}
	s = fixture()
	s.User("a").Initial = 102000
	if f = s.Forecast("a", testToday); f.EndBalance != -5000 {
		t.Fatalf("F06 %+v", f)
	}
	s = fixture()
	before := s.Forecast("a", testToday)
	if e := s.Resolve("a", "oct", "complete", "", "2026-10-26", 10000, testToday); e != nil {
		t.Fatal(e)
	}
	after := s.Forecast("a", testToday)
	if after.EndBalance != before.EndBalance || s.Balance("a") != 90000 {
		t.Fatalf("F07 before %+v after %+v", before, after)
	}
	if e := s.Resolve("a", "oct", "complete", "", "2026-10-26", 10000, testToday); e != nil {
		t.Fatal(e)
	}
	if len(s.Operations) != 2 {
		t.Fatal("F08 duplicate")
	}
	id := s.Occurrence("a", "oct").OperationID
	if e := s.DeleteOperation("a", id, true); e != nil {
		t.Fatal(e)
	}
	if s.Forecast("a", testToday).EndBalance != 85000 || s.Occurrence("a", "oct").Status != "pending" {
		t.Fatal("F15 deletion must reopen")
	}
	s = fixture()
	s.Occurrence("a", "oct").Date = "2026-09-30"
	if s.Forecast("a", testToday).EndBalance != 85000 {
		t.Fatal("F12 overdue omitted")
	}
	if e := s.Resolve("a", "oct", "skip", "", "", 0, testToday); e != nil {
		t.Fatal(e)
	}
	if s.Forecast("a", testToday).EndBalance != 95000 {
		t.Fatal("F12 skipped counted")
	}
	s = fixture()
	s.User("a").Started = "2026-07-15"
	if s.Forecast("a", testToday).Months != 2 {
		t.Fatal("F13 partial month")
	}
	s.User("a").Started = "2026-06-01"
	if s.Forecast("a", testToday).From != "2026-07-01" {
		t.Fatal("F14 last three")
	}
	s = fixture()
	s.Operations[0].Amount += 1000
	f = s.Forecast("a", testToday)
	if f.EndBalance != 83946 || f.NextExpense != 40326 {
		t.Fatalf("F16 %+v", f)
	}
	s = fixture()
	if e := s.Resolve("a", "oct", "move", "", "2026-11-01", 0, testToday); e != nil {
		t.Fatal(e)
	}
	f = s.Forecast("a", testToday)
	if f.Payments != 0 || f.NextPayments != 20000 {
		t.Fatalf("F17 %+v", f)
	}
}
func TestForecastRoundingAndIsolation(t *testing.T) {
	s := fixture()
	s.User("a").Started = "2026-08-01"
	s.Operations = []Operation{{ID: "a1", UserID: "a", Kind: "expense", Date: "2026-08-01", Amount: 10000}, {ID: "a2", UserID: "a", Kind: "expense", Date: "2026-09-01", Amount: 20000}, {ID: "b", UserID: "b", Kind: "expense", Date: "2026-09-01", Amount: 999999}}
	f := s.Forecast("a", testToday)
	if f.HistoryExpense != 30000 || f.HistoryDays != 61 || f.NextExpense != 24754 {
		t.Fatalf("F09 %+v", f)
	}
	s.Occurrences = append(s.Occurrences, Occurrence{ID: "linked", UserID: "a", OperationID: "a1", Status: "done"})
	f = s.Forecast("a", testToday)
	if f.HistoryExpense != 20000 {
		t.Fatal("regular cost in history")
	}
	if roundFraction(bigInt(-5), 2) != -3 {
		t.Fatal("negative half rounding")
	}
}
func TestOperationsAndCategories(t *testing.T) {
	s := fixture()
	op := Operation{CategoryID: "food", Kind: "expense", Date: "2026-10-26", Amount: 2000, Token: "one"}
	id, e := s.SaveOperation("a", op, testToday)
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.SaveOperation("a", op, testToday)
	if e != nil || again != id || len(s.Operations) != 2 {
		t.Fatal("duplicate form")
	}
	v := *s.Operation("a", id)
	v.Amount = 2500
	if _, e = s.SaveOperation("a", v, testToday); e != nil || s.Balance("a") != 97500 {
		t.Fatal("edit delta", e)
	}
	v.Amount = 0
	if _, e = s.SaveOperation("a", v, testToday); e == nil || s.Operation("a", id).Amount != 2500 {
		t.Fatal("invalid edit")
	}
	v.Amount = 100
	v.CategoryID = "foreign"
	if _, e = s.SaveOperation("a", v, testToday); e == nil {
		t.Fatal("foreign category")
	}
	if e = s.DeleteOperation("b", id, true); e == nil {
		t.Fatal("foreign delete")
	}
	if e = s.DeleteOperation("a", id, false); e == nil {
		t.Fatal("confirmation")
	}
	if e = s.SetBudget("a", "2026-08-01", 1000, true, testToday); e == nil {
		t.Fatal("start not fixed")
	}
	if e = s.SetBudget("a", "2026-07-01", 1000, false, testToday); e == nil {
		t.Fatal("initial confirm")
	}
	if e = s.DeleteOperation("a", id, true); e != nil {
		t.Fatal(e)
	}
	if e = s.SetBudget("a", "2026-08-01", 1000, true, testToday); e == nil {
		t.Fatal("start unlocked after delete")
	}
	c := *s.Category("a", "food")
	c.Archived = true
	if e = s.SaveCategory("a", c); e != nil {
		t.Fatal(e)
	}
	op.Token = "new"
	if _, e = s.SaveOperation("a", op, testToday); e == nil {
		t.Fatal("archived category")
	}
	if e = s.SaveCategory("a", Category{Name: "ПРОДУКТЫ", Kind: "expense"}); e == nil {
		t.Fatal("duplicate name")
	}
}
func TestFiltersLimits(t *testing.T) {
	s := fixture()
	for i := 0; i < 25; i++ {
		_, e := s.SaveOperation("a", Operation{CategoryID: "food", Kind: "expense", Amount: 100, Date: "2026-10-26", Token: fmt.Sprint(i)}, testToday)
		if e != nil {
			t.Fatal(e)
		}
	}
	rows, _, exp, total, e := s.Filter("a", Filter{From: "2026-10-26", To: "2026-10-26", Page: 2})
	if e != nil || len(rows) != 5 || total != 25 || exp != 2500 {
		t.Fatal("pagination sums", len(rows), total, exp, e)
	}
	if _, _, _, _, e = s.Filter("a", Filter{From: "2026-10-27", To: "2026-10-26"}); e == nil {
		t.Fatal("reverse dates")
	}
	if e = s.SetLimit("a", "food", "2026-10", 2500, false); e != nil {
		t.Fatal(e)
	}
	if e = s.SetLimit("a", "food", "2026-10", 2000, false); e != nil || len(s.Limits) != 1 {
		t.Fatal("limit unique")
	}
	if e = s.SetLimit("a", "foreign", "2026-10", 100, false); e == nil {
		t.Fatal("foreign limit")
	}
}
func TestCalendarAndRuleEdits(t *testing.T) {
	for _, tc := range []struct{ month, want string }{{"2028-02-01", "2028-02-29"}, {"2027-02-01", "2027-02-28"}} {
		if got := due(Rule{Day: 31}, mustDate(tc.month)); got != tc.want {
			t.Fatalf("F10/F11 %s", got)
		}
	}
	s := fixture()
	if e := s.SaveRule("a", Rule{Name: "Internet", Active: true, Kind: "expense", CategoryID: "food", Day: 30, Start: "2026-07-01", Amount: 10000}, testToday); e != nil {
		t.Fatal(e)
	}
	r := s.Rules[0]
	old := s.Occurrence("a", r.ID+":2026-09")
	oldAmount := old.Amount
	r.Amount = 12000
	if e := s.SaveRule("a", r, testToday); e != nil {
		t.Fatal(e)
	}
	if s.Occurrence("a", r.ID+":2026-09").Amount != oldAmount || s.Occurrence("a", r.ID+":2026-11").Amount != 12000 {
		t.Fatal("rule rewrote history")
	}
}
func TestStoreTransactionsAndBackups(t *testing.T) {
	path := filepath.Join(t.TempDir(), "balance.json")
	store, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = Open(path); e == nil {
		t.Fatal("second process allowed")
	}
	if e = store.Update(func(st *State) error { *st = fixture(); return nil }); e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Go(func() {
			if e := store.Update(func(st *State) error { return st.Resolve("a", "oct", "complete", "", "2026-10-26", 10000, testToday) }); e != nil {
				t.Error(e)
			}
		})
	}
	wg.Wait()
	st := store.Snapshot()
	if len(st.Operations) != 2 {
		t.Fatal("parallel duplicate")
	}
	sentinel := errors.New("rollback")
	if e = store.Update(func(st *State) error { st.Users[0].Initial = 0; return sentinel }); e != sentinel {
		t.Fatal(e)
	}
	if store.Snapshot().Users[0].Initial != 192000 {
		t.Fatal("rollback failed")
	}
	// Fail writes by moving the containing directory. Memory must remain unchanged.
	dir := filepath.Dir(path)
	moved := dir + "-moved"
	if e = os.Rename(dir, moved); e != nil {
		t.Fatal(e)
	}
	if e = store.Update(func(st *State) error { st.Users[0].Initial = 1; return nil }); !errors.Is(e, ErrStorage) {
		t.Fatal("write failure not detected", e)
	}
	if store.Snapshot().Users[0].Initial != 192000 {
		t.Fatal("write failure committed memory")
	}
	if store.Ready() == nil {
		t.Fatal("health with missing storage")
	}
	if e = os.Rename(moved, dir); e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 9; i++ {
		if e = store.Backup(testToday.AddDate(0, 0, i)); e != nil {
			t.Fatal(e)
		}
	}
	files, _ := filepath.Glob(path + ".backups/*.json")
	if len(files) != 7 {
		t.Fatal("retention", len(files))
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0600 {
		t.Fatal("permissions")
	}
	store.Close()
	restored, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	if len(restored.Snapshot().Operations) != 2 {
		t.Fatal("restart data")
	}
	restored.Close()
	if e = Restore(path, files[0]); e != nil {
		t.Fatal(e)
	}
}
func TestPasswordsAndLockout(t *testing.T) {
	store, e := Open(filepath.Join(t.TempDir(), "store.json"))
	if e != nil {
		t.Fatal(e)
	}
	defer store.Close()
	if _, e = store.Register("not-an-email", "longpassword123"); e == nil {
		t.Fatal("email")
	}
	if _, e = store.Register("user@example.by", "short"); e == nil {
		t.Fatal("short")
	}
	uid, e := store.Register("user@example.by", "longpassword123")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = store.Register("USER@example.by", "longpassword123"); e == nil {
		t.Fatal("duplicate email")
	}
	for i := 0; i < 5; i++ {
		if _, e = store.Authenticate("user@example.by", "wrong", testToday); e == nil {
			t.Fatal("bad password")
		}
	}
	if _, e = store.Authenticate("user@example.by", "longpassword123", testToday); e == nil {
		t.Fatal("not locked")
	}
	if got, e := store.Authenticate("user@example.by", "longpassword123", testToday.Add(16*time.Minute)); e != nil || got != uid {
		t.Fatal("unlock", e)
	}
	data, _ := os.ReadFile(store.path)
	if strings.Contains(string(data), "longpassword123") {
		t.Fatal("plaintext password")
	}
}

func bigInt(n int64) *big.Int { return big.NewInt(n) }

func TestOneToOneLinksAndRuleHistory(t *testing.T) {
	s := fixture()
	if e := s.Resolve("a", "oct", "link", "history", "", 0, testToday); e != nil {
		t.Fatal(e)
	}
	if e := s.Resolve("a", "nov", "link", "history", "", 0, testToday); e == nil {
		t.Fatal("one operation linked twice")
	}
	if s.Forecast("a", testToday).HistoryExpense != 0 {
		t.Fatal("linked operation in daily spending")
	}
	if e := s.Resolve("b", "oct", "skip", "", "", 0, testToday); e == nil {
		t.Fatal("foreign payment")
	}
	if e := s.SaveRule("a", Rule{Name: "Internet", Active: true, Kind: "expense", CategoryID: "food", Day: 30, Start: "2026-07-01", Amount: 10000}, testToday); e != nil {
		t.Fatal(e)
	}
	r := s.Rules[0]
	oct := r.ID + ":2026-10"
	nov := r.ID + ":2026-11"
	if e := s.Resolve("a", oct, "move", "", "2026-11-04", 0, testToday); e != nil {
		t.Fatal(e)
	}
	r.Day = 28
	r.Amount = 15000
	if e := s.SaveRule("a", r, testToday); e != nil {
		t.Fatal(e)
	}
	if s.Occurrence("a", oct).Date != "2026-11-04" || s.Occurrence("a", oct).Amount != 10000 {
		t.Fatal("individual move overwritten")
	}
	if s.Occurrence("a", nov).Date != "2026-11-28" || s.Occurrence("a", nov).Amount != 15000 {
		t.Fatal("future not updated")
	}
	if s.Occurrence("a", r.ID+":2026-09").Amount != 10000 {
		t.Fatal("past changed")
	}
	r.Day = 5
	if e := s.SaveRule("a", r, testToday); e != nil {
		t.Fatal(e)
	}
	for _, o := range s.Occurrences {
		if o.RuleID == r.ID && o.Date == "2026-10-05" {
			t.Fatal("retroactive occurrence from edit")
		}
	}
}

func TestDeletedOperationFormCannotRecreateRecord(t *testing.T) {
	s := fixture()
	op := Operation{CategoryID: "food", Kind: "expense", Date: "2026-10-26", Amount: 100, Token: "old-form"}
	id, e := s.SaveOperation("a", op, testToday)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.DeleteOperation("a", id, true); e != nil {
		t.Fatal(e)
	}
	if _, e = s.SaveOperation("a", op, testToday); e != nil {
		t.Fatal(e)
	}
	if len(s.Operations) != 1 {
		t.Fatal("browser back + replay recreated a deleted operation")
	}
}
