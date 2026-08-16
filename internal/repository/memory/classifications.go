package memory

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"time"
)

func (s *Store) CreateClassification(ctx context.Context, n domain.Classification) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, x := range s.classifications {
		if x.Code == n.Code {
			return domain.ErrConflict
		}
	}
	s.classifications[n.ID] = n
	return nil
}
func (s *Store) UpdateClassification(ctx context.Context, n domain.Classification) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.classifications[n.ID]; !ok {
		return domain.ErrNotFound
	}
	s.classifications[n.ID] = n
	return nil
}
func (s *Store) ClassificationByID(ctx context.Context, id string) (domain.Classification, error) {
	if err := check(ctx); err != nil {
		return domain.Classification{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.classifications[id]
	if !ok {
		return domain.Classification{}, domain.ErrNotFound
	}
	return n, nil
}
func (s *Store) ListClassifications(ctx context.Context) ([]domain.Classification, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Classification, 0, len(s.classifications))
	for _, n := range s.classifications {
		out = append(out, n)
	}
	return out, nil
}
func (s *Store) MoveClassification(ctx context.Context, id, parent string) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.faultMu.RLock()
	fail := s.failClassificationMove
	s.faultMu.RUnlock()
	if fail {
		return domain.ErrConflict
	}
	n, ok := s.classifications[id]
	if !ok {
		return domain.ErrNotFound
	}
	n.ParentID = parent
	n.UpdatedAt = time.Now()
	s.classifications[id] = n
	return nil
}
