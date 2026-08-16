package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/wyw14/cry002/internal/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	idemMu sync.Mutex
	idem   map[string]*idemResult
}
type idemResult struct {
	done  chan struct{}
	value string
	err   error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db, idem: map[string]*idemResult{}}
}
func dbError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	return err
}
func jsonText(v any) string    { b, _ := json.Marshal(v); return string(b) }
func fromJSON(s string, v any) { _ = json.Unmarshal([]byte(s), v) }
func (r *Repository) DoIdempotent(ctx context.Context, key string, fn func() (string, error)) (string, error) {
	r.idemMu.Lock()
	if old, ok := r.idem[key]; ok {
		r.idemMu.Unlock()
		select {
		case <-old.done:
			return old.value, old.err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	cur := &idemResult{done: make(chan struct{})}
	r.idem[key] = cur
	r.idemMu.Unlock()
	cur.value, cur.err = fn()
	close(cur.done)
	return cur.value, cur.err
}
