package store

import (
	"sync"
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Prompt    string    `json:"prompt"`
	Status    string    `json:"status"`
	ResultURL string    `json:"result_url,omitempty"`
	Owner     string    `json:"owner,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskStore struct {
	sync.RWMutex
	tasks    map[string]*Task
	onUpdate func(*Task)
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]*Task),
	}
}

func (s *TaskStore) Save(t *Task) {
	s.Lock()
	defer s.Unlock()
	t.UpdatedAt = time.Now()
	s.tasks[t.ID] = t
	if s.onUpdate != nil {
		s.onUpdate(t)
	}
}

func (s *TaskStore) Get(id string) (*Task, bool) {
	s.RLock()
	defer s.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStore) List() []*Task {
	s.RLock()
	defer s.RUnlock()
	list := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

func (s *TaskStore) ListByOwner(owner string) []*Task {
	s.RLock()
	defer s.RUnlock()
	list := make([]*Task, 0)
	for _, t := range s.tasks {
		if t.Owner == owner {
			list = append(list, t)
		}
	}
	return list
}

func (s *TaskStore) SetOnUpdate(f func(*Task)) {
	s.onUpdate = f
}
