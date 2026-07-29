package store

import (
	"sync"
	"time"

	"amol-nv/user_crud/internal/users"
)

type MemoryStore struct {
	users    *MemoryUserStore
	payments *MemoryPaymentStore
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:    NewMemoryUserStore(),
		payments: NewMemoryPaymentStore(),
	}
}

func (m *MemoryStore) User() *MemoryUserStore { return m.users }

func (m *MemoryStore) Payment() *MemoryPaymentStore { return m.payments }

// --- Users ---

type MemoryUserStore struct {
	mu     sync.RWMutex
	seq    int
	byID   map[string]*users.User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{byID: map[string]*users.User{}}
}

func (s *MemoryUserStore) Create(u *users.User) (*users.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	id := usersID(s.seq)
	now := time.Now().UTC()
	copy := *u
	copy.ID = id
	copy.CreatedAt = now
	copy.UpdatedAt = now
	s.byID[id] = &copy
	return &copy, nil
}

func (s *MemoryUserStore) GetByID(id string) (*users.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.byID[id]
	if !ok {
		return nil, users.ErrUserNotFound
	}
	copy := *u
	return &copy, nil
}

func (s *MemoryUserStore) List() ([]*users.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*users.User, 0, len(s.byID))
	for _, u := range s.byID {
		copy := *u
		out = append(out, &copy)
	}
	return out, nil
}

func (s *MemoryUserStore) Update(u *users.User) (*users.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[u.ID]; !ok {
		return nil, users.ErrUserNotFound
	}
	now := time.Now().UTC()
	copy := *u
	copy.UpdatedAt = now
	s.byID[u.ID] = &copy
	return &copy, nil
}

func (s *MemoryUserStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return users.ErrUserNotFound
	}
	delete(s.byID, id)
	return nil
}

func usersID(seq int) string {
	return "u-" + itoa(seq)
}

// --- Payments ---

type MemoryPaymentStore struct {
	mu     sync.RWMutex
	seq    int
	byID   map[string]*Payment
}

func NewMemoryPaymentStore() *MemoryPaymentStore {
	return &MemoryPaymentStore{byID: map[string]*Payment{}}
}

func (s *MemoryPaymentStore) Create(p *Payment) (*Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	id := paymentsID(s.seq)
	now := time.Now().UTC()
	copy := *p
	copy.ID = id
	copy.CreatedAt = now
	copy.UpdatedAt = now
	s.byID[id] = &copy
	return &copy, nil
}

func (s *MemoryPaymentStore) GetByID(id string) (*Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.byID[id]
	if !ok {
		return nil, ErrPaymentNotFound
	}
	copy := *p
	return &copy, nil
}

func (s *MemoryPaymentStore) List() ([]*Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Payment, 0, len(s.byID))
	for _, p := range s.byID {
		copy := *p
		out = append(out, &copy)
	}
	return out, nil
}

func (s *MemoryPaymentStore) Update(p *Payment) (*Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[p.ID]; !ok {
		return nil, ErrPaymentNotFound
	}
	now := time.Now().UTC()
	copy := *p
	copy.UpdatedAt = now
	s.byID[p.ID] = &copy
	return &copy, nil
}

func (s *MemoryPaymentStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return ErrPaymentNotFound
	}
	delete(s.byID, id)
	return nil
}

func paymentsID(seq int) string {
	return "p-" + itoa(seq)
}

// --- shared helpers ---

func itoa(i int) string {
	// small local helper to avoid importing strconv in multiple files
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	buf := make([]byte, 0, 12)
	for i > 0 {
		buf = append(buf, byte('0'+(i%10)))
		i /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for l, r := 0, len(buf)-1; l < r; l, r = l+1, r-1 {
		buf[l], buf[r] = buf[r], buf[l]
	}
	return string(buf)
}

// ErrPaymentNotFound is used by the payment store.
var ErrPaymentNotFound = users.ErrUserNotFound
