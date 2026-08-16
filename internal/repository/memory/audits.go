package memory

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

func (s *Store) AppendAudit(ctx context.Context, a domain.AuditEvent) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audits = append(s.audits, a)
	return nil
}
func (s *Store) ListAudits(ctx context.Context, resource, id string) ([]domain.AuditEvent, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.AuditEvent{}
	for _, a := range s.audits {
		if (resource == "" || a.Resource == resource) && (id == "" || a.ResourceID == id) {
			out = append(out, a)
		}
	}
	return out, nil
}
