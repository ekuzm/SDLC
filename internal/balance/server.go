package balance

import (
	"bytes"
	"crypto/subtle"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed web/*
var web embed.FS

type Config struct {
	Secure bool
	Now    func() time.Time
}
type Server struct {
	store     *Store
	config    Config
	sessions  sessions
	templates *template.Template
	handler   http.Handler
}
type categoryView struct {
	Category
	Spent, Limit, Left Money
	Percent            int
	Status             string
}
type occurrenceView struct {
	Occurrence
	StatusLabel string
	Overdue     bool
}
type Page struct {
	Title, View, Email, CSRF, Token, Error, Notice, Today, Month, From, To, CompareFrom, CompareTo, Partial, Change string
	Values                                                                                                          url.Values
	User                                                                                                            User
	Categories                                                                                                      []Category
	CategoryViews                                                                                                   []categoryView
	Operations                                                                                                      []Operation
	LinkOperations                                                                                                  []Operation
	Rules                                                                                                           []Rule
	Occurrences                                                                                                     []occurrenceView
	Forecast                                                                                                        Forecast
	Balance, Income, Expense, CompareExpense                                                                        Money
	Total, Page, Pages                                                                                              int
	Links                                                                                                           []string
	EditOperation                                                                                                   *Operation
	EditCategory                                                                                                    *Category
	EditRule                                                                                                        *Rule
	Delete                                                                                                          *Operation
}

func NewServer(store *Store, cfg Config) (*Server, error) {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	s := &Server{store: store, config: cfg}
	funcs := template.FuncMap{"money": func(m Money) string { return m.String() }, "date": func(d string) string {
		t, e := Date(d)
		if e != nil {
			return d
		}
		return t.Format("02.01.2006")
	}, "kind": func(k string) string {
		if k == "income" {
			return "Доход"
		}
		return "Расход"
	}, "value": func(v url.Values, k, def string) string {
		if v.Has(k) {
			return v.Get(k)
		}
		return def
	}, "cat": func(c []Category, id string) string {
		for _, x := range c {
			if x.ID == id {
				return x.Name
			}
		}
		return "—"
	}, "inc": func(i int) int { return i + 1 }, "sub": func(a, b int) int { return a - b }, "neg": func(m Money) bool { return m < 0 }, "percent": func(part, total Money) string {
		if total == 0 {
			return "—"
		}
		return fmt.Sprintf("%.1f %%", float64(part)*100/float64(total))
	}}

	funcs["checked"] = func(v url.Values, key string, def bool) bool {
		if v.Has("csrf") {
			return v.Get(key) == "yes"
		}
		return def
	}
	funcs["decimal"] = func(s string) string { return strings.ReplaceAll(s, ".", ",") }
	funcs["oKind"] = func(o *Operation) string {
		if o != nil {
			return o.Kind
		}
		return "expense"
	}
	funcs["oDate"] = func(o *Operation) string {
		if o != nil {
			return o.Date
		}
		return ""
	}
	funcs["oAmount"] = func(o *Operation) string {
		if o != nil {
			return o.Amount.String()
		}
		return ""
	}
	funcs["oCategory"] = func(o *Operation) string {
		if o != nil {
			return o.CategoryID
		}
		return ""
	}
	funcs["oComment"] = func(o *Operation) string {
		if o != nil {
			return o.Comment
		}
		return ""
	}
	funcs["cName"] = func(c *Category) string {
		if c != nil {
			return c.Name
		}
		return ""
	}
	funcs["cKind"] = func(c *Category) string {
		if c != nil {
			return c.Kind
		}
		return "expense"
	}
	funcs["rName"] = func(r *Rule) string {
		if r != nil {
			return r.Name
		}
		return ""
	}
	funcs["rKind"] = func(r *Rule) string {
		if r != nil {
			return r.Kind
		}
		return "expense"
	}
	funcs["rCategory"] = func(r *Rule) string {
		if r != nil {
			return r.CategoryID
		}
		return ""
	}
	funcs["rAmount"] = func(r *Rule) string {
		if r != nil {
			return r.Amount.String()
		}
		return ""
	}
	funcs["rDay"] = func(r *Rule) string {
		if r != nil {
			return strconv.Itoa(r.Day)
		}
		return "1"
	}
	funcs["rStart"] = func(r *Rule) string {
		if r != nil {
			return r.Start
		}
		return ""
	}
	funcs["rEnd"] = func(r *Rule) string {
		if r != nil {
			return r.End
		}
		return ""
	}
	var e error
	s.templates, e = template.New("pages").Funcs(funcs).ParseFS(web, "web/*.html")
	if e != nil {
		return nil, e
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /static/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		data, _ := web.ReadFile("web/style.css")
		_, _ = w.Write(data)
	})
	mux.HandleFunc("GET /static/app.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		data, _ := web.ReadFile("web/app.js")
		_, _ = w.Write(data)
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/overview", http.StatusSeeOther) })
	for _, path := range []string{"login", "register"} {
		mux.HandleFunc("GET /"+path, s.authPage)
		mux.HandleFunc("POST /"+path, s.authPost)
	}
	for _, path := range []string{"overview", "operations", "categories", "payments", "analytics", "settings"} {
		mux.HandleFunc("GET /"+path, s.page)
	}
	for _, path := range []string{"budget", "category", "limit", "operation", "operation/delete", "rule", "payment", "logout"} {
		mux.HandleFunc("POST /"+path, s.mutate)
	}
	s.handler = http.NewCrossOriginProtection().Handler(mux)
	return s, nil
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := ID()
	w.Header().Set("X-Request-ID", id)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "same-origin")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; base-uri 'none'; frame-ancestors 'none'; form-action 'self'")
	if s.config.Secure {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")
	}
	defer func() {
		if recover() != nil {
			slog.Error("request panic", "request_id", id)
			http.Error(w, "Не удалось обработать запрос. Код: "+id, 500)
		}
	}()
	s.handler.ServeHTTP(w, r)
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if e := s.store.Ready(); e != nil {
		slog.Error("storage unavailable", "request_id", w.Header().Get("X-Request-ID"))
		http.Error(w, "Хранилище недоступно", 503)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
}
func (s *Server) cookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: true, Secure: s.config.Secure, SameSite: http.SameSiteStrictMode, MaxAge: maxAge})
}
func (s *Server) current(r *http.Request) (session, bool) {
	c, e := r.Cookie("balance_session")
	if e != nil {
		return session{}, false
	}
	return s.sessions.get(c.Value, s.config.Now())
}
func (s *Server) render(w http.ResponseWriter, p Page, status int) {
	var b bytes.Buffer
	if e := s.templates.ExecuteTemplate(&b, "layout", p); e != nil {
		slog.Error("render failed", "error", e)
		http.Error(w, "Не удалось открыть страницу", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(b.Bytes())
}
func (s *Server) authPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.current(r); ok {
		http.Redirect(w, r, "/overview", 303)
		return
	}
	token := ID()
	s.cookie(w, "balance_csrf", token, 1800)
	view := strings.TrimPrefix(r.URL.Path, "/")
	s.render(w, Page{View: view, Title: "Вход в Баланс", CSRF: token, Values: url.Values{}}, 200)
}
func readForm(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
	if e := r.ParseForm(); e != nil {
		return errors.New("Форма слишком большая или повреждена")
	}
	return nil
}
func tokenEqual(a, b string) bool {
	return len(a) >= 20 && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
func (s *Server) authPost(w http.ResponseWriter, r *http.Request) {
	if e := readForm(w, r); e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	c, e := r.Cookie("balance_csrf")
	if e != nil || !tokenEqual(c.Value, r.PostForm.Get("csrf")) {
		http.Error(w, "Форма устарела. Обновите страницу.", 403)
		return
	}
	view := strings.TrimPrefix(r.URL.Path, "/")
	var uid string
	if view == "register" {
		uid, e = s.store.Register(r.PostForm.Get("email"), r.PostForm.Get("password"))
	} else {
		uid, e = s.store.Authenticate(r.PostForm.Get("email"), r.PostForm.Get("password"), s.config.Now())
	}
	if e != nil {
		r.PostForm.Del("password")
		status := 422
		if errors.Is(e, ErrStorage) {
			status = 503
			slog.Error("auth save failed", "request_id", w.Header().Get("X-Request-ID"))
		}
		s.render(w, Page{View: view, Title: "Вход в Баланс", CSRF: c.Value, Error: e.Error(), Values: r.PostForm}, status)
		return
	}
	key, _ := s.sessions.new(uid, s.config.Now())
	s.cookie(w, "balance_session", key, 0)
	s.cookie(w, "balance_csrf", "", -1)
	http.Redirect(w, r, "/overview", 303)
}
func (s *Server) build(r *http.Request, v session, values url.Values) (Page, error) {
	now := s.config.Now()
	today := Day(now)
	st := s.store.Snapshot()
	// Derive future occurrences for reads; mutations materialize them in the same transaction.
	st.Materialize(v.UserID, today)
	u := st.User(v.UserID)
	if u == nil {
		return Page{}, errors.New("Учётная запись не найдена")
	}
	view := strings.TrimPrefix(r.URL.Path, "/")
	titles := map[string]string{"overview": "Мой бюджет", "operations": "Операции", "categories": "Категории и лимиты", "payments": "Плановые платежи", "analytics": "Аналитика и прогнозы", "settings": "Настройка бюджета"}
	if titles[view] == "" {
		view = "overview"
	}
	p := Page{Title: titles[view], View: view, User: *u, Email: u.Email, CSRF: v.CSRF, Token: ID(), Today: today.Format("2006-01-02"), Month: today.Format("2006-01"), Balance: st.Balance(u.ID), Forecast: st.Forecast(u.ID, today), Values: values, Page: 1}
	if r.URL.Query().Get("saved") == "1" {
		p.Notice = "Изменения сохранены"
	}
	for _, c := range st.Categories {
		if c.UserID == u.ID {
			p.Categories = append(p.Categories, c)
		}
	}
	sort.Slice(p.Categories, func(i, j int) bool { return p.Categories[i].Name < p.Categories[j].Name })
	q := r.URL.Query()
	if values.Get("month") != "" {
		p.Month = values.Get("month")
	}
	if q.Get("month") != "" {
		if _, e := Date(q.Get("month") + "-01"); e != nil {
			return p, e
		}
		p.Month = q.Get("month")
	}
	start, _ := Date(p.Month + "-01")
	end := start.AddDate(0, 1, -1)
	p.Income, p.Expense = st.Totals(u.ID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	for _, c := range p.Categories {
		cv := categoryView{Category: c}
		for _, o := range st.Operations {
			if o.UserID == u.ID && o.CategoryID == c.ID && o.Date >= start.Format("2006-01-02") && o.Date <= end.Format("2006-01-02") {
				cv.Spent += o.Amount
			}
		}
		for _, l := range st.Limits {
			if l.UserID == u.ID && l.CategoryID == c.ID && l.Month == p.Month {
				cv.Limit = l.Amount
			}
		}
		cv.Status = "Без лимита"
		if cv.Limit > 0 {
			cv.Left = cv.Limit - cv.Spent
			cv.Percent = int(cv.Spent * 100 / cv.Limit)
			cv.Status = "В пределах лимита"
			if cv.Left < 0 {
				cv.Status = "Лимит превышен"
			} else if cv.Left == 0 {
				cv.Status = "Лимит достигнут"
			}
		}
		p.CategoryViews = append(p.CategoryViews, cv)
	}
	for _, o := range st.Occurrences {
		if o.UserID == u.ID {
			label := map[string]string{"pending": "Ожидается", "done": "Исполнено", "skipped": "Пропущено"}[o.Status]
			overdue := o.Status == "pending" && o.Date < p.Today
			if overdue {
				label = "Просрочен"
			}
			p.Occurrences = append(p.Occurrences, occurrenceView{o, label, overdue})
		}
	}
	sort.Slice(p.Occurrences, func(i, j int) bool {
		if p.Occurrences[i].Status == "pending" && p.Occurrences[j].Status != "pending" {
			return true
		}
		if p.Occurrences[j].Status == "pending" && p.Occurrences[i].Status != "pending" {
			return false
		}
		return p.Occurrences[i].Date < p.Occurrences[j].Date
	})
	for _, r := range st.Rules {
		if r.UserID == u.ID {
			p.Rules = append(p.Rules, r)
		}
	}
	if view == "operations" {
		page, _ := strconv.Atoi(q.Get("page"))
		if page < 1 {
			page = 1
		}
		if page > 1000000 {
			page = 1000000
		}
		p.Page = page
		f := Filter{q.Get("from"), q.Get("to"), q.Get("kind"), q.Get("category"), page}
		var e error
		p.Operations, p.Income, p.Expense, p.Total, e = st.Filter(u.ID, f)
		if e != nil {
			return p, e
		}
		p.Pages = (p.Total + 19) / 20
		for i := 1; i <= p.Pages; i++ {
			link := url.Values{}
			for k, vs := range q {
				link[k] = append([]string(nil), vs...)
			}
			link.Set("page", strconv.Itoa(i))
			p.Links = append(p.Links, "/operations?"+link.Encode())
		}
		if id := q.Get("edit"); id != "" {
			p.EditOperation = st.Operation(u.ID, id)
			if p.EditOperation == nil {
				return p, errors.New("Операция не найдена")
			}
		}
		if id := q.Get("delete"); id != "" {
			p.Delete = st.Operation(u.ID, id)
			if p.Delete == nil {
				return p, errors.New("Операция не найдена")
			}
		}
	}
	if view == "categories" && q.Get("edit") != "" {
		p.EditCategory = st.Category(u.ID, q.Get("edit"))
		if p.EditCategory == nil {
			return p, errors.New("Категория не найдена")
		}
	}
	if view == "payments" {
		linked := map[string]bool{}
		for _, o := range st.Occurrences {
			if o.UserID == u.ID && o.OperationID != "" {
				linked[o.OperationID] = true
			}
		}
		for _, o := range st.Operations {
			if o.UserID == u.ID && !linked[o.ID] {
				p.LinkOperations = append(p.LinkOperations, o)
			}
		}
	}
	if view == "payments" && q.Get("edit") != "" {
		p.EditRule = st.Rule(u.ID, q.Get("edit"))
		if p.EditRule == nil {
			return p, errors.New("Правило не найдено")
		}
	}
	if view == "analytics" {
		p.From = start.Format("2006-01-02")
		p.To = end.Format("2006-01-02")
		if p.Month == today.Format("2006-01") {
			p.To = p.Today
		}
		prev := start.AddDate(0, -1, 0)
		p.CompareFrom = prev.Format("2006-01-02")
		p.CompareTo = prev.AddDate(0, 1, -1).Format("2006-01-02")
		if p.Month == today.Format("2006-01") {
			day := today.Day()
			if day > Days(prev) {
				day = Days(prev)
			}
			p.CompareTo = prev.AddDate(0, 0, day-1).Format("2006-01-02")
		}
		for key, dst := range map[string]*string{"from": &p.From, "to": &p.To, "compare_from": &p.CompareFrom, "compare_to": &p.CompareTo} {
			if q.Get(key) != "" {
				*dst = q.Get(key)
			}
			if _, e := Date(*dst); e != nil {
				return p, e
			}
		}
		if p.From > p.To || p.CompareFrom > p.CompareTo {
			return p, errors.New("Начало периода позже конца")
		}
		_, p.Expense = st.Totals(u.ID, p.From, p.To)
		_, p.CompareExpense = st.Totals(u.ID, p.CompareFrom, p.CompareTo)
		p.Change = "— (база сравнения равна нулю)"
		if p.CompareExpense != 0 {
			p.Change = fmt.Sprintf("%+.2f %%", float64(p.Expense-p.CompareExpense)*100/float64(p.CompareExpense))
		}
		// A range shorter than a complete calendar month is explicitly labelled.
		fd, _ := Date(p.From)
		td, _ := Date(p.To)
		cd, _ := Date(p.CompareFrom)
		ct, _ := Date(p.CompareTo)
		if fd.Day() != 1 || td.Day() != Days(td) || cd.Day() != 1 || ct.Day() != Days(ct) {
			p.Partial = "Выбраны неполные календарные месяцы"
		}
		for i := range p.CategoryViews {
			cv := &p.CategoryViews[i]
			cv.Spent = 0
			for _, o := range st.Operations {
				if o.UserID == u.ID && o.CategoryID == cv.ID && o.Date >= p.From && o.Date <= p.To {
					cv.Spent += o.Amount
				}
			}
		}
	}
	return p, nil
}
func (s *Server) page(w http.ResponseWriter, r *http.Request) {
	v, ok := s.current(r)
	if !ok {
		http.Redirect(w, r, "/login", 303)
		return
	}
	p, e := s.build(r, v, r.URL.Query())
	if e != nil {
		p.Error = e.Error()
		s.render(w, p, 422)
		return
	}
	s.render(w, p, 200)
}
func (s *Server) mutate(w http.ResponseWriter, r *http.Request) {
	v, ok := s.current(r)
	if !ok {
		http.Redirect(w, r, "/login", 303)
		return
	}
	if e := readForm(w, r); e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	if !tokenEqual(v.CSRF, r.PostForm.Get("csrf")) {
		http.Error(w, "Форма устарела. Обновите страницу.", 403)
		return
	}
	if r.URL.Path == "/logout" {
		if c, e := r.Cookie("balance_session"); e == nil {
			s.sessions.delete(c.Value)
		}
		s.cookie(w, "balance_session", "", -1)
		http.Redirect(w, r, "/login", 303)
		return
	}
	f := r.PostForm
	target := "overview"
	now := s.config.Now()
	switch r.URL.Path {
	case "/budget":
		target = "settings"
	case "/category", "/limit":
		target = "categories"
	case "/operation", "/operation/delete":
		target = "operations"
	case "/rule", "/payment":
		target = "payments"
	}
	err := s.store.Update(func(st *State) error {
		st.Materialize(v.UserID, now)
		switch r.URL.Path {
		case "/budget":
			amount, e := ParseMoney(f.Get("initial"))
			if e != nil {
				return e
			}
			return st.SetBudget(v.UserID, f.Get("start"), amount, f.Get("confirm") == "yes", now)
		case "/category":
			return st.SaveCategory(v.UserID, Category{ID: f.Get("id"), Name: f.Get("name"), Kind: f.Get("kind"), Archived: f.Get("archived") == "yes"})
		case "/limit":
			amount := Money(0)
			if f.Get("remove") != "yes" {
				var e error
				amount, e = ParseMoney(f.Get("amount"))
				if e != nil {
					return e
				}
			}
			return st.SetLimit(v.UserID, f.Get("category"), f.Get("month"), amount, f.Get("remove") == "yes")
		case "/operation":
			amount, e := ParseMoney(f.Get("amount"))
			if e != nil {
				return e
			}
			_, e = st.SaveOperation(v.UserID, Operation{ID: f.Get("id"), CategoryID: f.Get("category"), Kind: f.Get("kind"), Date: f.Get("date"), Amount: amount, Comment: f.Get("comment"), Token: f.Get("token")}, now)
			return e
		case "/operation/delete":
			return st.DeleteOperation(v.UserID, f.Get("id"), f.Get("confirm") == "yes")
		case "/rule":
			amount, e := ParseMoney(f.Get("amount"))
			if e != nil {
				return e
			}
			day, e := strconv.Atoi(f.Get("day"))
			if e != nil {
				return errors.New("Укажите день 1–31")
			}
			return st.SaveRule(v.UserID, Rule{ID: f.Get("id"), Name: f.Get("name"), CategoryID: f.Get("category"), Kind: f.Get("kind"), Amount: amount, Day: day, Start: f.Get("start"), End: f.Get("end"), Active: f.Get("active") == "yes"}, now)
		case "/payment":
			amount := Money(0)
			if f.Get("action") == "complete" {
				var e error
				amount, e = ParseMoney(f.Get("amount"))
				if e != nil {
					return e
				}
			}
			return st.Resolve(v.UserID, f.Get("id"), f.Get("action"), f.Get("operation"), f.Get("date"), amount, now)
		}
		return errors.New("Неизвестное действие")
	})
	if err != nil {
		copyReq := r.Clone(r.Context())
		u := *r.URL
		u.Path = "/" + target
		u.RawQuery = ""
		q := url.Values{}
		if f.Get("id") != "" && r.URL.Path != "/payment" {
			q.Set("edit", f.Get("id"))
		}
		u.RawQuery = q.Encode()
		copyReq.URL = &u
		p, _ := s.build(copyReq, v, f)
		p.Error = err.Error()
		p.Token = f.Get("token")
		if p.Token == "" {
			p.Token = ID()
		}
		status := 422
		if errors.Is(err, ErrStorage) {
			status = 503
			slog.Error("save failed", "request_id", w.Header().Get("X-Request-ID"))
		}
		s.render(w, p, status)
		return
	}
	http.Redirect(w, r, "/"+target+"?saved=1", 303)
}
