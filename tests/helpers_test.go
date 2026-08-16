package tests

import (
	"github.com/wyw14/cry002/internal/domain"
	"sync"
	"time"
)

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

type seqIDs struct {
	mu sync.Mutex
	n  int
}

func (s *seqIDs) New() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
	return "id-" + time.Unix(int64(s.n), 0).UTC().Format("150405")
}
func actor(id string, role domain.Role) domain.User {
	return domain.User{ID: id, Role: role, Status: domain.UserActive}
}
