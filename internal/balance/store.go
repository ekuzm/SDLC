package balance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"syscall"
	"time"
)

// Store is a single-process transactional store. Every mutation works on a copy;
// memory is published only after fsync + atomic rename of the complete snapshot.
type Store struct {
	mu    sync.Mutex
	path  string
	state State
	lock  *os.File
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		lock.Close()
		return nil, errors.New("хранилище уже открыто другим процессом")
	}
	s := &Store{path: path, lock: lock, state: State{Version: 1}}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		err = s.write(s.state)
	} else if err == nil {
		err = json.Unmarshal(data, &s.state)
		if err == nil && s.state.Version != 1 {
			err = errors.New("неподдерживаемая версия хранилища")
		}
	}
	if err != nil {
		s.Close()
		return nil, err
	}
	if err = os.Chmod(path, 0600); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	if s.lock != nil {
		return s.lock.Close()
	}
	return nil
}
func clone(in State) State {
	out := in
	out.Receipts = append([]Receipt(nil), in.Receipts...)
	out.Users = append([]User(nil), in.Users...)
	out.Categories = append([]Category(nil), in.Categories...)
	out.Operations = append([]Operation(nil), in.Operations...)
	out.Rules = append([]Rule(nil), in.Rules...)
	out.Occurrences = append([]Occurrence(nil), in.Occurrences...)
	out.Limits = append([]Limit(nil), in.Limits...)
	return out
}
func (s *Store) Snapshot() State { s.mu.Lock(); defer s.mu.Unlock(); return clone(s.state) }
func (s *Store) Update(fn func(*State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := clone(s.state)
	if err := fn(&next); err != nil {
		return err
	}
	if err := s.write(next); err != nil {
		return fmt.Errorf("%w: %v", ErrStorage, err)
	}
	s.state = next
	return nil
}

var ErrStorage = errors.New("Не удалось сохранить данные. Повторите попытку; введённые значения сохранены в форме")

func atomicFile(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".balance-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err = f.Write(data); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if err = os.Rename(tmp, path); err != nil {
		return err
	}
	// Rename is the commit point. Syncing the directory improves crash durability;
	// a post-commit sync failure cannot be reported as an uncommitted transaction.
	if dir, e := os.Open(filepath.Dir(path)); e == nil {
		_ = dir.Sync()
		_ = dir.Close()
	}
	return nil
}
func (s *Store) write(state State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return atomicFile(s.path, data)
}
func (s *Store) Ready() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, e := os.Open(s.path)
	if e != nil {
		return e
	}
	defer f.Close()
	var check State
	if e = json.NewDecoder(f).Decode(&check); e != nil {
		return e
	}
	if check.Version != 1 {
		return errors.New("invalid version")
	}
	probe, e := os.CreateTemp(filepath.Dir(s.path), ".ready-*")
	if e != nil {
		return e
	}
	name := probe.Name()
	e = probe.Close()
	_ = os.Remove(name)
	return e
}
func (s *Store) Backup(now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	dir := s.path + ".backups"
	if e := os.MkdirAll(dir, 0700); e != nil {
		return e
	}
	data, e := json.Marshal(s.state)
	if e != nil {
		return e
	}
	name := filepath.Join(dir, now.Format("2006-01-02")+".json")
	if _, e = os.Stat(name); errors.Is(e, os.ErrNotExist) {
		if e = atomicFile(name, data); e != nil {
			return e
		}
	} else if e != nil {
		return e
	}
	files, e := filepath.Glob(filepath.Join(dir, "*.json"))
	if e != nil {
		return e
	}
	sort.Strings(files)
	for len(files) > 7 {
		if e = os.Remove(files[0]); e != nil {
			return e
		}
		files = files[1:]
	}
	return nil
}
func Restore(path, backup string) error {
	data, e := os.ReadFile(backup)
	if e != nil {
		return e
	}
	var state State
	if e = json.Unmarshal(data, &state); e != nil {
		return e
	}
	if state.Version != 1 {
		return errors.New("unsupported backup version")
	}
	s, e := Open(path)
	if e != nil {
		return e
	}
	defer s.Close()
	return s.Update(func(target *State) error { *target = state; return nil })
}
