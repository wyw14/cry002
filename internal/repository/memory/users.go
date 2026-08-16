package memory

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

func (s *Store) CreateUser(ctx context.Context, u domain.User) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.emails[u.Email]; ok {
		return domain.ErrConflict
	}
	s.users[u.ID] = u
	s.emails[u.Email] = u.ID
	return nil
}
func (s *Store) UpdateUser(ctx context.Context, u domain.User) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[u.ID]; !ok {
		return domain.ErrNotFound
	}
	s.users[u.ID] = u
	s.emails[u.Email] = u.ID
	return nil
}
func (s *Store) UserByID(ctx context.Context, id string) (domain.User, error) {
	if err := check(ctx); err != nil {
		return domain.User{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}
func (s *Store) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	if err := check(ctx); err != nil {
		return domain.User{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.emails[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return s.users[id], nil
}
func (s *Store) ListUsers(ctx context.Context) ([]domain.User, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out, nil
}
func (s *Store) CreateDepartment(ctx context.Context, d domain.Department) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.departments[d.ID] = d
	return nil
}
func (s *Store) ListDepartments(ctx context.Context) ([]domain.Department, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Department, 0, len(s.departments))
	for _, d := range s.departments {
		out = append(out, d)
	}
	return out, nil
}
