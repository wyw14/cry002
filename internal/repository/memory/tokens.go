package memory

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"time"
)

func (s *Store) StoreRefreshToken(ctx context.Context, t domain.RefreshToken) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[t.Hash] = t
	return nil
}
func (s *Store) RefreshToken(ctx context.Context, h string) (domain.RefreshToken, error) {
	if err := check(ctx); err != nil {
		return domain.RefreshToken{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tokens[h]
	if !ok {
		return domain.RefreshToken{}, domain.ErrNotFound
	}
	return t, nil
}
func (s *Store) ConsumeRefreshToken(ctx context.Context, h string, now time.Time) (domain.RefreshToken, error) {
	if err := check(ctx); err != nil {
		return domain.RefreshToken{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tokens[h]
	if !ok {
		return domain.RefreshToken{}, domain.ErrUnauthorized
	}
	if err := t.CanRotate(now); err != nil {
		return domain.RefreshToken{}, err
	}
	t.UsedAt = &now
	s.tokens[h] = t
	return t, nil
}
func (s *Store) RevokeTokenFamily(ctx context.Context, f string, now time.Time) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for h, t := range s.tokens {
		if t.FamilyID == f {
			t.RevokedAt = &now
			s.tokens[h] = t
		}
	}
	return nil
}

var _ interface {
	StoreRefreshToken(context.Context, domain.RefreshToken) error
} = (*Store)(nil)
