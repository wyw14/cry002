package application

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"time"
)

type BorrowService struct {
	repo  Repository
	clock Clock
	ids   IDGenerator
}

func NewBorrowService(repo Repository, clock Clock, ids IDGenerator) *BorrowService {
	return &BorrowService{repo: repo, clock: clock, ids: ids}
}
func (s *BorrowService) Apply(ctx context.Context, actor domain.User, caseID, purpose string, due time.Time, meta RequestMeta) (domain.BorrowRequest, error) {
	if purpose == "" || due.Before(s.clock.Now()) {
		return domain.BorrowRequest{}, domain.ErrValidation
	}
	c, err := s.repo.CaseByID(ctx, caseID)
	if err != nil {
		return domain.BorrowRequest{}, err
	}
	if c.Status != domain.CaseArchived {
		return domain.BorrowRequest{}, domain.ErrInvalidState
	}
	b := domain.BorrowRequest{ID: s.ids.New(), CaseID: caseID, ApplicantID: actor.ID, Purpose: purpose, Status: domain.BorrowPending, DueAt: due, CreatedAt: s.clock.Now(), UpdatedAt: s.clock.Now()}
	if err := s.repo.CreateBorrow(ctx, b); err != nil {
		return domain.BorrowRequest{}, err
	}
	return b, s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: "borrow.applied", Resource: "borrow", ResourceID: b.ID, RequestID: meta.RequestID, CreatedAt: s.clock.Now()})
}
func (s *BorrowService) Approve(ctx context.Context, actor domain.User, id string, approved bool, meta RequestMeta) error {
	b, err := s.repo.BorrowByID(ctx, id)
	if err != nil {
		return err
	}
	if !b.CanApprove(actor) {
		return domain.ErrForbidden
	}
	b.ApproverID = actor.ID
	if approved {
		b.Status = domain.BorrowApproved
	} else {
		b.Status = domain.BorrowRejected
	}
	b.UpdatedAt = s.clock.Now()
	if err := s.repo.UpdateBorrow(ctx, b); err != nil {
		return err
	}
	return s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: "borrow.reviewed", Resource: "borrow", ResourceID: id, RequestID: meta.RequestID, CreatedAt: s.clock.Now()})
}
func (s *BorrowService) Checkout(ctx context.Context, actor domain.User, id string, meta RequestMeta) error {
	b, err := s.repo.BorrowByID(ctx, id)
	if err != nil {
		return err
	}
	if b.ApplicantID != actor.ID && actor.Role != domain.RoleArchivist && !actor.IsAdministrator() {
		return domain.ErrForbidden
	}
	if b.Status != domain.BorrowApproved {
		return domain.ErrInvalidState
	}
	if err := s.repo.CheckoutCase(ctx, b.CaseID, id, s.clock.Now()); err != nil {
		return err
	}
	return s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: "borrow.checked_out", Resource: "borrow", ResourceID: id, RequestID: meta.RequestID, CreatedAt: s.clock.Now()})
}
func (s *BorrowService) Renew(ctx context.Context, actor domain.User, id string, days int, meta RequestMeta) error {
	b, err := s.repo.BorrowByID(ctx, id)
	if err != nil {
		return err
	}
	if b.ApplicantID != actor.ID {
		return domain.ErrForbidden
	}
	if days < 1 || days > 30 || !b.CanRenew(s.clock.Now()) {
		return domain.ErrInvalidState
	}
	b.DueAt = b.DueAt.Add(time.Duration(days) * 24 * time.Hour)
	b.Renewals++
	b.UpdatedAt = s.clock.Now()
	if err := s.repo.UpdateBorrow(ctx, b); err != nil {
		return err
	}
	return s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: "borrow.renewed", Resource: "borrow", ResourceID: id, RequestID: meta.RequestID, CreatedAt: s.clock.Now()})
}
func (s *BorrowService) Return(ctx context.Context, actor domain.User, id string, meta RequestMeta) error {
	b, err := s.repo.BorrowByID(ctx, id)
	if err != nil {
		return err
	}
	if actor.ID != b.ApplicantID && actor.Role != domain.RoleArchivist && !actor.IsAdministrator() {
		return domain.ErrForbidden
	}
	if b.Status != domain.BorrowCheckedOut && b.Status != domain.BorrowOverdue {
		return domain.ErrInvalidState
	}
	now := s.clock.Now()
	b.Status = domain.BorrowReturned
	b.ReturnedAt = &now
	b.UpdatedAt = now
	if err := s.repo.UpdateBorrow(ctx, b); err != nil {
		return err
	}
	c, err := s.repo.CaseByID(ctx, b.CaseID)
	if err != nil {
		return err
	}
	c.Status = domain.CaseArchived
	c.UpdatedAt = now
	ms, _ := s.repo.Materials(ctx, c.ID)
	if err := s.repo.UpdateCase(ctx, c, ms); err != nil {
		return err
	}
	return s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: "borrow.returned", Resource: "borrow", ResourceID: id, RequestID: meta.RequestID, CreatedAt: now})
}
func (s *BorrowService) MarkOverdue(ctx context.Context, now time.Time) error {
	borrows, err := s.repo.ListBorrows(ctx, "")
	if err != nil {
		return err
	}
	for _, b := range borrows {
		if b.Status == domain.BorrowCheckedOut && now.After(b.DueAt) {
			b.Status = domain.BorrowOverdue
			if err := s.repo.UpdateBorrow(ctx, b); err != nil {
				return err
			}
		}
	}
	return nil
}
