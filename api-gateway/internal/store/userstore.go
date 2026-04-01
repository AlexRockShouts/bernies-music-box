package store

import (
	"encoding/json"
	"os"
	"sync"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type UserStore struct {
	sync.RWMutex
	users map[string]*User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*User),
	}
}

func (s *UserStore) LoadJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var us []User
	if err := json.Unmarshal(data, &us); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	for _, u := range us {
		s.users[u.Username] = &u
	}
	return nil
}

func (s *UserStore) Get(username string) (*User, bool) {
	s.RLock()
	defer s.RUnlock()
	u, ok := s.users[username]
	return u, ok
}

func (s *UserStore) SaveJSON(filename string) error {
	s.RLock()
	defer s.RUnlock()
	var us []User
	for _, u := range s.users {
		us = append(us, *u)
	}
	data, err := json.MarshalIndent(us, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
