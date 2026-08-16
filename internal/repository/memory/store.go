package memory

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"sort"
	"strings"
	"sync"
	"time"
)

type idemResult struct {
	done  chan struct{}
	value string
	err   error
}
type Store struct {
	mu                     sync.RWMutex
	users                  map[string]domain.User
	emails                 map[string]string
	departments            map[string]domain.Department
	classifications        map[string]domain.Classification
	cases                  map[string]domain.CaseFile
	materials              map[string][]domain.Material
	versions               map[string][]domain.CaseVersion
	reviews                map[string][]domain.Review
	borrows                map[string]domain.BorrowRequest
	attachments            map[string]domain.Attachment
	audits                 []domain.AuditEvent
	tokens                 map[string]domain.RefreshToken
	idemMu                 sync.Mutex
	idem                   map[string]*idemResult
	faultMu                sync.RWMutex
	delay                  time.Duration
	failClassificationMove bool
}

func New() *Store {
	return &Store{users: map[string]domain.User{}, emails: map[string]string{}, departments: map[string]domain.Department{}, classifications: map[string]domain.Classification{}, cases: map[string]domain.CaseFile{}, materials: map[string][]domain.Material{}, versions: map[string][]domain.CaseVersion{}, reviews: map[string][]domain.Review{}, borrows: map[string]domain.BorrowRequest{}, attachments: map[string]domain.Attachment{}, tokens: map[string]domain.RefreshToken{}, idem: map[string]*idemResult{}}
}
func (s *Store) SetDelay(d time.Duration) { s.faultMu.Lock(); s.delay = d; s.faultMu.Unlock() }
func (s *Store) SetClassificationFailure(v bool) {
	s.faultMu.Lock()
	s.failClassificationMove = v
	s.faultMu.Unlock()
}
func check(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func (s *Store) wait(ctx context.Context) error {
	s.faultMu.RLock()
	d := s.delay
	s.faultMu.RUnlock()
	if d <= 0 {
		return check(ctx)
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
func (s *Store) DoIdempotent(ctx context.Context, key string, fn func() (string, error)) (string, error) {
	s.idemMu.Lock()
	if old, ok := s.idem[key]; ok {
		s.idemMu.Unlock()
		select {
		case <-old.done:
			return old.value, old.err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	cur := &idemResult{done: make(chan struct{})}
	s.idem[key] = cur
	s.idemMu.Unlock()
	cur.value, cur.err = fn()
	close(cur.done)
	return cur.value, cur.err
}

var _ = sort.Slice
var _ = strings.TrimSpace
