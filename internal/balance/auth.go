package balance

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/mail"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const passwordIterations = 600000

func passwordKey(password, salt string) string {
	key, err := pbkdf2.Key(sha256.New, password, []byte(salt), passwordIterations, 32)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}
func (s *Store) Register(email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	addr, e := mail.ParseAddress(email)
	if e != nil || addr.Address != email || len(email) > 254 {
		return "", errors.New("Введите корректный email")
	}
	if n := utf8.RuneCountInString(password); n < 12 || n > 128 {
		return "", errors.New("Пароль должен содержать 12–128 символов")
	}
	u := User{ID: ID(), Email: email, Salt: ID()}
	u.PasswordHash = passwordKey(password, u.Salt)
	e = s.Update(func(st *State) error {
		for _, v := range st.Users {
			if v.Email == email {
				return errors.New("Этот email уже зарегистрирован")
			}
		}
		st.Users = append(st.Users, u)
		for _, c := range []Category{{Name: "Зарплата", Kind: "income"}, {Name: "Продукты", Kind: "expense"}, {Name: "Транспорт", Kind: "expense"}, {Name: "Жильё", Kind: "expense"}, {Name: "Другое", Kind: "expense"}} {
			if e := st.SaveCategory(u.ID, c); e != nil {
				return e
			}
		}
		return nil
	})
	return u.ID, e
}

var ErrLogin = errors.New("Неверные данные входа или вход временно недоступен")

func (s *Store) Authenticate(email, password string, now time.Time) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	st := s.Snapshot()
	var user User
	for _, u := range st.Users {
		if u.Email == email {
			user = u
			break
		}
	}
	salt := user.Salt
	if salt == "" {
		salt = "unknown-account-timing-padding"
	}
	hash := passwordKey(password, salt)
	if user.ID == "" {
		return "", ErrLogin
	}
	valid := subtle.ConstantTimeCompare([]byte(hash), []byte(user.PasswordHash)) == 1
	accepted := false
	err := s.Update(func(st *State) error {
		u := st.User(user.ID)
		if now.Before(u.BlockedUntil) {
			return nil
		}
		if valid {
			u.Failed = 0
			u.FailWindow = time.Time{}
			u.BlockedUntil = time.Time{}
			accepted = true
			return nil
		}
		if u.FailWindow.IsZero() || now.Sub(u.FailWindow) >= 15*time.Minute {
			u.Failed = 0
			u.FailWindow = now
		}
		u.Failed++
		if u.Failed >= 5 {
			u.BlockedUntil = now.Add(15 * time.Minute)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if !accepted {
		return "", ErrLogin
	}
	return user.ID, nil
}

type session struct {
	UserID, CSRF string
	Last         time.Time
}
type sessions struct {
	mu    sync.Mutex
	items map[string]session
}

func (s *sessions) new(uid string, now time.Time) (string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.items == nil {
		s.items = map[string]session{}
	}
	for k, v := range s.items {
		if now.Sub(v.Last) >= 30*time.Minute {
			delete(s.items, k)
		}
	}
	key, csrf := ID(), ID()
	s.items[key] = session{uid, csrf, now}
	return key, csrf
}
func (s *sessions) get(key string, now time.Time) (session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[key]
	if !ok {
		return v, false
	}
	if now.Sub(v.Last) >= 30*time.Minute {
		delete(s.items, key)
		return session{}, false
	}
	v.Last = now
	s.items[key] = v
	return v, true
}
func (s *sessions) delete(key string) { s.mu.Lock(); defer s.mu.Unlock(); delete(s.items, key) }
