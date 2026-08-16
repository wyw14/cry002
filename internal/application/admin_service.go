package application

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

type AdminStats struct {
	TotalCases     int `json:"total_cases"`
	ArchivedCases  int `json:"archived_cases"`
	PendingBorrows int `json:"pending_borrows"`
	AuditEvents    int `json:"audit_events"`
}
type AdminService struct{ repo Repository }

func NewAdminService(repo Repository) *AdminService { return &AdminService{repo: repo} }
func (s *AdminService) Stats(ctx context.Context, actor domain.User) (AdminStats, error) {
	if !actor.IsAdministrator() && actor.Role != domain.RoleArchivist {
		return AdminStats{}, domain.ErrForbidden
	}
	cases, total, err := s.repo.ListCases(ctx, CaseFilter{Page: 1, PageSize: 100})
	if err != nil {
		return AdminStats{}, err
	}
	arch := 0
	for _, c := range cases {
		if c.Status == domain.CaseArchived {
			arch++
		}
	}
	borrows, err := s.repo.ListBorrows(ctx, "")
	if err != nil {
		return AdminStats{}, err
	}
	pending := 0
	for _, b := range borrows {
		if b.Status == domain.BorrowPending {
			pending++
		}
	}
	audits, err := s.repo.ListAudits(ctx, "", "")
	if err != nil {
		return AdminStats{}, err
	}
	return AdminStats{TotalCases: total, ArchivedCases: arch, PendingBorrows: pending, AuditEvents: len(audits)}, nil
}
