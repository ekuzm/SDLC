package balance

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Money int64

const MaxMoney Money = 99999999999

func ParseMoney(s string) (Money, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", ".")
	negative := strings.HasPrefix(s, "-")
	if negative {
		s = s[1:]
	}
	parts := strings.Split(s, ".")
	if len(parts) > 2 || len(parts[0]) == 0 || len(parts[0]) > 9 {
		return 0, errors.New("Введите сумму до 999999999,99 с точностью до копейки")
	}
	for _, part := range parts {
		for _, c := range part {
			if c < '0' || c > '9' {
				return 0, errors.New("Сумма должна быть числом")
			}
		}
	}
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, errors.New("Некорректная сумма")
	}
	cents := int64(0)
	if len(parts) == 2 {
		if len(parts[1]) < 1 || len(parts[1]) > 2 {
			return 0, errors.New("Допустимо не более двух знаков после запятой")
		}
		p := parts[1]
		if len(p) == 1 {
			p += "0"
		}
		cents, _ = strconv.ParseInt(p, 10, 64)
	}
	m := Money(whole*100 + cents)
	if m > MaxMoney {
		return 0, errors.New("Сумма слишком большая")
	}
	if negative {
		m = -m
	}
	return m, nil
}
func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	return fmt.Sprintf("%s%d,%02d", sign, m/100, m%100)
}
func ID() string { return rand.Text() }
func Date(s string) (time.Time, error) {
	t, e := time.Parse("2006-01-02", s)
	if e != nil || t.Format("2006-01-02") != s {
		return time.Time{}, errors.New("Укажите корректную дату")
	}
	return t, nil
}
func Day(t time.Time) time.Time   { return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC) }
func Month(t time.Time) time.Time { return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC) }
func Days(t time.Time) int        { return Month(t).AddDate(0, 1, -1).Day() }
func validKind(s string) bool     { return s == "income" || s == "expense" }

type User struct {
	ID, Email, Salt, PasswordHash string
	Started                       string
	Initial                       Money
	StartLocked                   bool
	Failed                        int
	FailWindow, BlockedUntil      time.Time
}
type Category struct {
	ID, UserID, Name, Kind string
	Archived               bool
}
type Operation struct {
	ID, UserID, CategoryID, Kind, Date, Comment, Token string
	Amount                                             Money
	Created                                            time.Time
}
type Rule struct {
	ID, UserID, CategoryID, Kind, Name, Start, End, Revised string
	Amount                                                  Money
	Day                                                     int
	Active                                                  bool
}
type Occurrence struct {
	ID, UserID, RuleID, CategoryID, Kind, Name, Date, Status, OperationID string
	Amount                                                                Money
	Moved                                                                 bool
}
type Limit struct {
	UserID, CategoryID, Month string
	Amount                    Money
}
type Receipt struct{ UserID, Token, OperationID string }

type State struct {
	Version     int
	Receipts    []Receipt
	Users       []User
	Categories  []Category
	Operations  []Operation
	Rules       []Rule
	Occurrences []Occurrence
	Limits      []Limit
}

func (s *State) User(id string) *User {
	for i := range s.Users {
		if s.Users[i].ID == id {
			return &s.Users[i]
		}
	}
	return nil
}
func (s *State) Category(uid, id string) *Category {
	for i := range s.Categories {
		v := &s.Categories[i]
		if v.UserID == uid && v.ID == id {
			return v
		}
	}
	return nil
}
func (s *State) Operation(uid, id string) *Operation {
	for i := range s.Operations {
		v := &s.Operations[i]
		if v.UserID == uid && v.ID == id {
			return v
		}
	}
	return nil
}
func (s *State) Rule(uid, id string) *Rule {
	for i := range s.Rules {
		v := &s.Rules[i]
		if v.UserID == uid && v.ID == id {
			return v
		}
	}
	return nil
}
func (s *State) Occurrence(uid, id string) *Occurrence {
	for i := range s.Occurrences {
		v := &s.Occurrences[i]
		if v.UserID == uid && v.ID == id {
			return v
		}
	}
	return nil
}
func (s *State) CategoryName(uid, id string) string {
	if c := s.Category(uid, id); c != nil {
		return c.Name
	}
	return "—"
}
func (s *State) SetBudget(uid, start string, initial Money, confirm bool, today time.Time) error {
	u := s.User(uid)
	if u == nil {
		return errors.New("Бюджет не найден")
	}
	d, e := Date(start)
	if e != nil {
		return e
	}
	if d.After(Day(today)) {
		return errors.New("Начало учёта не может быть в будущем")
	}
	if d.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return errors.New("Дата начала должна быть не раньше 01.01.2000")
	}
	if u.StartLocked && u.Started != start {
		return errors.New("Дата начала фиксируется после первой операции")
	}
	if u.Started != "" && u.Initial != initial && !confirm {
		return errors.New("Подтвердите пересчёт начального остатка")
	}
	u.Started = start
	u.Initial = initial
	return nil
}
func (s *State) SaveCategory(uid string, c Category) error {
	c.Name = strings.TrimSpace(c.Name)
	if utf8.RuneCountInString(c.Name) < 1 || utf8.RuneCountInString(c.Name) > 50 || !validKind(c.Kind) {
		return errors.New("Название: 1–50 символов; выберите тип категории")
	}
	for _, v := range s.Categories {
		if v.UserID == uid && v.ID != c.ID && v.Kind == c.Kind && strings.EqualFold(v.Name, c.Name) {
			return errors.New("Категория с таким названием уже существует")
		}
	}
	if c.ID != "" {
		v := s.Category(uid, c.ID)
		if v == nil {
			return errors.New("Категория не найдена")
		}
		if v.Kind != c.Kind {
			return errors.New("Тип существующей категории изменить нельзя")
		}
		v.Name = c.Name
		v.Archived = c.Archived
		return nil
	}
	c.ID = ID()
	c.UserID = uid
	s.Categories = append(s.Categories, c)
	return nil
}
func (s *State) validateOperation(uid string, o Operation, today time.Time) error {
	u := s.User(uid)
	if u == nil || u.Started == "" {
		return errors.New("Сначала настройте бюджет")
	}
	if o.Amount <= 0 || o.Amount > MaxMoney {
		return errors.New("Сумма должна быть от 0,01 до 999999999,99 BYN")
	}
	d, e := Date(o.Date)
	if e != nil {
		return e
	}
	if o.Date < u.Started || d.After(Day(today)) {
		return errors.New("Дата операции должна быть от начала учёта до сегодняшнего дня")
	}
	if utf8.RuneCountInString(o.Comment) > 250 {
		return errors.New("Комментарий не длиннее 250 символов")
	}
	c := s.Category(uid, o.CategoryID)
	if c == nil || c.Kind != o.Kind {
		return errors.New("Выберите собственную категорию нужного типа")
	}
	if c.Archived {
		old := s.Operation(uid, o.ID)
		if old == nil || old.CategoryID != c.ID || old.Kind != o.Kind {
			return errors.New("Категория архивирована")
		}
	}
	return nil
}
func (s *State) SaveOperation(uid string, o Operation, today time.Time) (string, error) {
	if o.ID == "" {
		if len(o.Token) > 128 {
			return "", errors.New("Некорректный ключ отправки")
		}
		for _, receipt := range s.Receipts {
			if receipt.UserID == uid && receipt.Token == o.Token {
				return receipt.OperationID, nil
			}
		}
		if o.Token == "" {
			return "", errors.New("Отсутствует ключ отправки формы")
		}
		for _, v := range s.Operations {
			if v.UserID == uid && v.Token == o.Token {
				return v.ID, nil
			}
		}
	}
	if e := s.validateOperation(uid, o, today); e != nil {
		return "", e
	}
	if o.ID != "" {
		v := s.Operation(uid, o.ID)
		if v == nil {
			return "", errors.New("Операция не найдена")
		}
		for _, occ := range s.Occurrences {
			if occ.UserID == uid && occ.OperationID == o.ID && occ.Kind != o.Kind {
				return "", errors.New("Сначала удалите связь с плановым платежом, чтобы сменить тип")
			}
		}
		o.UserID = uid
		o.Created = v.Created
		o.Token = v.Token
		*v = o
		return o.ID, nil
	}
	o.ID = ID()
	o.UserID = uid
	o.Created = today
	s.Operations = append(s.Operations, o)
	s.Receipts = append(s.Receipts, Receipt{uid, o.Token, o.ID})
	s.User(uid).StartLocked = true
	return o.ID, nil
}
func (s *State) DeleteOperation(uid, id string, confirmed bool) error {
	if !confirmed {
		return errors.New("Подтвердите удаление")
	}
	for i, o := range s.Operations {
		if o.ID == id && o.UserID == uid {
			for j := range s.Occurrences {
				p := &s.Occurrences[j]
				if p.UserID == uid && p.OperationID == id {
					p.OperationID = ""
					p.Status = "pending"
				}
			}
			s.Operations = append(s.Operations[:i], s.Operations[i+1:]...)
			return nil
		}
	}
	return errors.New("Операция не найдена")
}
func (s *State) SetLimit(uid, category, month string, amount Money, remove bool) error {
	c := s.Category(uid, category)
	if c == nil || c.Kind != "expense" {
		return errors.New("Лимит доступен только для собственной категории расходов")
	}
	if _, e := Date(month + "-01"); e != nil {
		return errors.New("Укажите месяц лимита")
	}
	if !remove && (amount <= 0 || amount > MaxMoney) {
		return errors.New("Лимит должен быть положительным")
	}
	for i, v := range s.Limits {
		if v.UserID == uid && v.CategoryID == category && v.Month == month {
			if remove {
				s.Limits = append(s.Limits[:i], s.Limits[i+1:]...)
			} else {
				s.Limits[i].Amount = amount
			}
			return nil
		}
	}
	if !remove {
		s.Limits = append(s.Limits, Limit{uid, category, month, amount})
	}
	return nil
}
func due(r Rule, m time.Time) string {
	day := r.Day
	if day > Days(m) {
		day = Days(m)
	}
	return time.Date(m.Year(), m.Month(), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

// Materialize stable monthly occurrences. Past and fulfilled occurrences are never rewritten.
func (s *State) Materialize(uid string, today time.Time) {
	end := Month(today).AddDate(0, 2, 0)
	existing := make(map[string]bool, len(s.Occurrences))
	for _, o := range s.Occurrences {
		existing[o.ID] = true
	}
	for _, r := range s.Rules {
		if r.UserID != uid || !r.Active {
			continue
		}
		start, _ := Date(r.Start)
		for m := Month(start); m.Before(end); m = m.AddDate(0, 1, 0) {
			d := due(r, m)
			if d < r.Start || r.End != "" && d > r.End {
				continue
			}
			id := r.ID + ":" + m.Format("2006-01")
			if !existing[id] && (r.Revised == "" || d > r.Revised) {
				s.Occurrences = append(s.Occurrences, Occurrence{ID: id, UserID: uid, RuleID: r.ID, CategoryID: r.CategoryID, Kind: r.Kind, Name: r.Name, Date: d, Amount: r.Amount, Status: "pending"})
				existing[id] = true
			}
		}
	}
}
func (s *State) SaveRule(uid string, r Rule, today time.Time) error {
	u := s.User(uid)
	if u == nil || u.Started == "" {
		return errors.New("Сначала настройте бюджет")
	}
	r.Name = strings.TrimSpace(r.Name)
	c := s.Category(uid, r.CategoryID)
	if !validKind(r.Kind) || c == nil || c.Kind != r.Kind || c.Archived || utf8.RuneCountInString(r.Name) < 1 || utf8.RuneCountInString(r.Name) > 50 || r.Day < 1 || r.Day > 31 || r.Amount <= 0 || r.Amount > MaxMoney {
		return errors.New("Проверьте название (1–50), категорию, положительную сумму и день (1–31)")
	}
	if _, e := Date(r.Start); e != nil {
		return e
	}
	if r.Start < u.Started {
		return errors.New("Правило не может начинаться раньше бюджета")
	}
	if r.End != "" {
		if _, e := Date(r.End); e != nil {
			return e
		}
		if r.End < r.Start {
			return errors.New("Окончание раньше начала")
		}
	}
	if r.ID != "" {
		v := s.Rule(uid, r.ID)
		if v == nil {
			return errors.New("Правило не найдено")
		}
		if r.Start != v.Start {
			return errors.New("Дата начала существующего правила фиксирована")
		}
		r.UserID = uid
		r.Revised = Day(today).Format("2006-01-02")
		*v = r
		kept := s.Occurrences[:0]
		for _, o := range s.Occurrences {
			if o.RuleID == r.ID && o.UserID == uid && o.Date > r.Revised && o.Status == "pending" && !o.Moved {
				if !r.Active {
					continue
				}
				keyMonth := strings.TrimPrefix(o.ID, r.ID+":")
				m, err := Date(keyMonth + "-01")
				if err != nil {
					return err
				}
				newDate := due(r, m)
				if r.End != "" && newDate > r.End {
					continue
				}
				// A schedule edit never introduces an occurrence in the past.
				if newDate > r.Revised {
					o.Date = newDate
					o.Name = r.Name
					o.Kind = r.Kind
					o.CategoryID = r.CategoryID
					o.Amount = r.Amount
				}
			}
			kept = append(kept, o)
		}
		s.Occurrences = kept
	} else {
		r.ID = ID()
		r.UserID = uid
		s.Rules = append(s.Rules, r)
	}
	s.Materialize(uid, today)
	return nil
}
func (s *State) Resolve(uid, id, action, operationID, date string, amount Money, today time.Time) error {
	p := s.Occurrence(uid, id)
	if p == nil {
		return errors.New("Плановый платёж не найден")
	}
	if p.Status == "done" {
		if action == "complete" || action == "link" {
			return nil
		}
		return errors.New("Платёж уже исполнен")
	}
	switch action {
	case "skip":
		p.Status = "skipped"
		return nil
	case "resume":
		p.Status = "pending"
		return nil
	case "move":
		if _, e := Date(date); e != nil {
			return e
		}
		if date < s.User(uid).Started {
			return errors.New("Дата раньше начала учёта")
		}
		p.Date = date
		p.Moved = true
		p.Status = "pending"
		return nil
	case "complete", "link":
		if p.Status != "pending" {
			return errors.New("Сначала верните платёж в ожидаемые")
		}
		var o *Operation
		if action == "link" {
			o = s.Operation(uid, operationID)
			if o == nil || o.Kind != p.Kind {
				return errors.New("Выберите собственную операцию того же типа")
			}
			for _, v := range s.Occurrences {
				if v.OperationID == o.ID {
					return errors.New("Операция уже связана с другим платежом")
				}
			}
		} else {
			op := Operation{CategoryID: p.CategoryID, Kind: p.Kind, Date: date, Amount: amount, Comment: p.Name, Token: "occurrence:" + p.ID + ":" + ID()}
			opID, e := s.SaveOperation(uid, op, today)
			if e != nil {
				return e
			}
			o = s.Operation(uid, opID)
		}
		p.Status = "done"
		p.OperationID = o.ID
		return nil
	default:
		return errors.New("Неизвестное действие")
	}
}
func (s *State) Balance(uid string) Money {
	u := s.User(uid)
	if u == nil {
		return 0
	}
	b := u.Initial
	for _, o := range s.Operations {
		if o.UserID == uid {
			if o.Kind == "income" {
				b += o.Amount
			} else {
				b -= o.Amount
			}
		}
	}
	return b
}
func (s *State) Totals(uid, from, to string) (income, expense Money) {
	for _, o := range s.Operations {
		if o.UserID == uid && o.Date >= from && o.Date <= to {
			if o.Kind == "income" {
				income += o.Amount
			} else {
				expense += o.Amount
			}
		}
	}
	return
}

type Filter struct {
	From, To, Kind, Category string
	Page                     int
}

func (s *State) Filter(uid string, f Filter) ([]Operation, Money, Money, int, error) {
	for _, d := range []string{f.From, f.To} {
		if d != "" {
			if _, e := Date(d); e != nil {
				return nil, 0, 0, 0, e
			}
		}
	}
	if f.From != "" && f.To != "" && f.From > f.To {
		return nil, 0, 0, 0, errors.New("Начало периода позже окончания")
	}
	if f.Kind != "" && !validKind(f.Kind) {
		return nil, 0, 0, 0, errors.New("Некорректный тип")
	}
	rows := []Operation{}
	var income, expense Money
	for _, o := range s.Operations {
		if o.UserID != uid || f.From != "" && o.Date < f.From || f.To != "" && o.Date > f.To || f.Kind != "" && o.Kind != f.Kind || f.Category != "" && o.CategoryID != f.Category {
			continue
		}
		rows = append(rows, o)
		if o.Kind == "income" {
			income += o.Amount
		} else {
			expense += o.Amount
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Date != rows[j].Date {
			return rows[i].Date > rows[j].Date
		}
		if !rows[i].Created.Equal(rows[j].Created) {
			return rows[i].Created.After(rows[j].Created)
		}
		return rows[i].ID > rows[j].ID
	})
	total := len(rows)
	page := f.Page
	if page < 1 {
		page = 1
	}
	start := (page - 1) * 20
	if start > total {
		start = total
	}
	end := start + 20
	if end > total {
		end = total
	}
	return rows[start:end], income, expense, total, nil
}

type Forecast struct {
	Available                                                                        bool
	Months, HistoryDays, RemainingDays, NextDays                                     int
	From, To, AsOf                                                                   string
	HistoryExpense, Balance, Income, Payments, NextPayments, EndBalance, NextExpense Money
	Daily                                                                            string
}

// roundFraction rounds the final rational amount (kopecks), never the daily average.
func roundFraction(n *big.Int, den int64) Money {
	d := big.NewInt(den)
	q, r := new(big.Int), new(big.Int)
	q.QuoRem(n, d, r)
	abs := new(big.Int).Abs(r)
	if abs.Mul(abs, big.NewInt(2)).Cmp(d) >= 0 {
		if n.Sign() < 0 {
			q.Sub(q, big.NewInt(1))
		} else {
			q.Add(q, big.NewInt(1))
		}
	}
	return Money(q.Int64())
}
func (s *State) Forecast(uid string, today time.Time) Forecast {
	today = Day(today)
	m := Month(today)
	next := m.AddDate(0, 1, 0)
	nextEnd := next.AddDate(0, 1, -1)
	f := Forecast{AsOf: today.Format("2006-01-02"), Balance: s.Balance(uid), RemainingDays: Days(today) - today.Day(), NextDays: Days(next)}
	u := s.User(uid)
	if u == nil || u.Started == "" {
		return f
	}
	for i := 3; i >= 1; i-- {
		h := m.AddDate(0, -i, 0)
		if h.Format("2006-01-02") >= u.Started {
			f.Months++
			f.HistoryDays += Days(h)
			if f.From == "" {
				f.From = h.Format("2006-01-02")
			}
			f.To = h.AddDate(0, 1, -1).Format("2006-01-02")
		}
	}
	if f.Months == 0 {
		return f
	}
	f.Available = true
	linked := map[string]bool{}
	for _, o := range s.Occurrences {
		if o.UserID == uid && o.Status == "done" {
			linked[o.OperationID] = true
		}
	}
	for _, o := range s.Operations {
		if o.UserID == uid && o.Kind == "expense" && o.Date >= f.From && o.Date <= f.To && !linked[o.ID] {
			f.HistoryExpense += o.Amount
		}
	}
	for _, o := range s.Occurrences {
		if o.UserID != uid || o.Status != "pending" {
			continue
		}
		if o.Date < next.Format("2006-01-02") {
			if o.Kind == "income" {
				f.Income += o.Amount
			} else {
				f.Payments += o.Amount
			}
		} else if o.Date <= nextEnd.Format("2006-01-02") && o.Kind == "expense" {
			f.NextPayments += o.Amount
		}
	}
	den := int64(f.HistoryDays)
	daily := new(big.Rat).SetFrac(big.NewInt(int64(f.HistoryExpense)), big.NewInt(den*100))
	f.Daily = daily.FloatString(4)
	base := new(big.Int).Mul(big.NewInt(int64(f.Balance+f.Income-f.Payments)), big.NewInt(den))
	cost := new(big.Int).Mul(big.NewInt(int64(f.HistoryExpense)), big.NewInt(int64(f.RemainingDays)))
	f.EndBalance = roundFraction(new(big.Int).Sub(base, cost), den)
	n := new(big.Int).Mul(big.NewInt(int64(f.HistoryExpense)), big.NewInt(int64(f.NextDays)))
	n.Add(n, new(big.Int).Mul(big.NewInt(int64(f.NextPayments)), big.NewInt(den)))
	f.NextExpense = roundFraction(n, den)
	return f
}
