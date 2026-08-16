package memory

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"time"
)

func (s *Store) CreateBorrow(ctx context.Context, b domain.BorrowRequest) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.borrows[b.ID] = b
	return nil
}
func (s *Store) BorrowByID(ctx context.Context, id string) (domain.BorrowRequest, error) {
	if err := check(ctx); err != nil {
		return domain.BorrowRequest{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.borrows[id]
	if !ok {
		return domain.BorrowRequest{}, domain.ErrNotFound
	}
	return b, nil
}
func (s *Store) ListBorrows(ctx context.Context, status string) ([]domain.BorrowRequest, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.BorrowRequest{}
	for _, b := range s.borrows {
		if status == "" || string(b.Status) == status {
			out = append(out, b)
		}
	}
	return out, nil
}
func (s *Store) UpdateBorrow(ctx context.Context, b domain.BorrowRequest) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.borrows[b.ID]; !ok {
		return domain.ErrNotFound
	}
	s.borrows[b.ID] = b
	return nil
}
func (s *Store) CheckoutCase(ctx context.Context, caseID, borrowID string, now time.Time) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c, caseOK := s.cases[caseID]
	b, borrowOK := s.borrows[borrowID]
	if !caseOK || !borrowOK {
		return domain.ErrNotFound
	}
	if c.Status != domain.CaseArchived || b.Status != domain.BorrowApproved {
		return domain.ErrAlreadyBorrowed
	}
	c.Status = domain.CaseBorrowed
	c.UpdatedAt = now
	b.Status = domain.BorrowCheckedOut
	b.CheckedOutAt = &now
	b.UpdatedAt = now
	s.cases[caseID] = c
	s.borrows[borrowID] = b
	return nil
}
