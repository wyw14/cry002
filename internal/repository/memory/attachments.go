package memory

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

func (s *Store) CreateAttachment(ctx context.Context, a domain.Attachment) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attachments[a.ID] = a
	return nil
}
func (s *Store) Attachments(ctx context.Context, caseID string) ([]domain.Attachment, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Attachment{}
	for _, a := range s.attachments {
		if a.CaseID == caseID {
			out = append(out, a)
		}
	}
	return out, nil
}
func (s *Store) AttachmentByID(ctx context.Context, id string) (domain.Attachment, error) {
	if err := check(ctx); err != nil {
		return domain.Attachment{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.attachments[id]
	if !ok {
		return domain.Attachment{}, domain.ErrNotFound
	}
	return a, nil
}
