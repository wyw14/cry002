package application

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/wyw14/cry002/internal/domain"
)

type CaseService struct {
	repo  Repository
	clock Clock
	ids   IDGenerator
}

func NewCaseService(repo Repository, clock Clock, ids IDGenerator) *CaseService {
	return &CaseService{repo: repo, clock: clock, ids: ids}
}
func (s *CaseService) Create(ctx context.Context, actor domain.User, c domain.CaseFile, materials []domain.Material, meta RequestMeta) (domain.CaseFile, error) {
	if actor.Status != domain.UserActive {
		return domain.CaseFile{}, domain.ErrForbidden
	}
	c.ID = s.ids.New()
	c.CreatedBy = actor.ID
	c.Status = domain.CaseDraft
	c.Version = 1
	c.CreatedAt = s.clock.Now()
	c.UpdatedAt = c.CreatedAt
	for i := range materials {
		materials[i].ID = s.ids.New()
		materials[i].CaseID = c.ID
		materials[i].CreatedAt = c.CreatedAt
	}
	if err := s.repo.CreateCase(ctx, c, materials); err != nil {
		return domain.CaseFile{}, err
	}
	_ = s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: "case.created", Resource: "case", ResourceID: c.ID, RequestID: meta.RequestID, CreatedAt: s.clock.Now()})
	return c, nil
}
func (s *CaseService) Submit(ctx context.Context, actor domain.User, id string, meta RequestMeta) error {
	c, err := s.repo.CaseByID(ctx, id)
	if err != nil {
		return err
	}
	ms, _ := s.repo.Materials(ctx, id)
	if c.CreatedBy != actor.ID && actor.Role != domain.RoleArchivist {
		return domain.ErrForbidden
	}
	if err := c.ValidateForSubmit(ms); err != nil {
		return err
	}
	if !c.CanTransition(domain.CasePending, actor) {
		return domain.ErrInvalidState
	}
	c.Status = domain.CasePending
	c.UpdatedAt = s.clock.Now()
	if err := s.repo.UpdateCase(ctx, c, ms); err != nil {
		return err
	}
	return s.audit(ctx, actor, "case.submitted", id, meta)
}
func (s *CaseService) Review(ctx context.Context, actor domain.User, id string, approved bool, opinion string, meta RequestMeta) error {
	c, err := s.repo.CaseByID(ctx, id)
	if err != nil {
		return err
	}
	if actor.Role != domain.RoleArchivist && actor.Role != domain.RoleAuditor && !actor.IsAdministrator() {
		return domain.ErrForbidden
	}
	to := domain.CaseArchived
	if !approved {
		to = domain.CaseDraft
	}
	if !c.CanTransition(to, actor) {
		return domain.ErrInvalidState
	}
	c.Status = to
	c.UpdatedAt = s.clock.Now()
	ms, _ := s.repo.Materials(ctx, id)
	if err := s.repo.UpdateCase(ctx, c, ms); err != nil {
		return err
	}
	if err := s.repo.CreateReview(ctx, domain.Review{ID: s.ids.New(), CaseID: id, ReviewerID: actor.ID, Approved: approved, Opinion: opinion, CreatedAt: s.clock.Now()}); err != nil {
		return err
	}
	return s.audit(ctx, actor, "case.reviewed", id, meta)
}
func (s *CaseService) UpdateMetadata(ctx context.Context, actor domain.User, id, reason string, patch domain.CaseFile, meta RequestMeta) error {
	c, err := s.repo.CaseByID(ctx, id)
	if err != nil {
		return err
	}
	if c.Status != domain.CaseArchived && c.Status != domain.CaseSealed {
		return domain.ErrInvalidState
	}
	if actor.Role != domain.RoleArchivist && !actor.IsAdministrator() {
		return domain.ErrForbidden
	}
	ms, _ := s.repo.Materials(ctx, id)
	before := domain.CaseSnapshot{Case: c, Materials: ms}
	patch.ID = c.ID
	patch.Status = c.Status
	patch.Version = c.Version + 1
	patch.CreatedBy = c.CreatedBy
	patch.CreatedAt = c.CreatedAt
	patch.UpdatedAt = s.clock.Now()
	if err := s.repo.UpdateCase(ctx, patch, ms); err != nil {
		return err
	}
	diff := []string{}
	if c.Title != patch.Title {
		diff = append(diff, "title")
	}
	if c.Summary != patch.Summary {
		diff = append(diff, "summary")
	}
	if c.SecurityLevel != patch.SecurityLevel {
		diff = append(diff, "security_level")
	}
	if err := s.repo.CreateVersion(ctx, domain.CaseVersion{ID: s.ids.New(), CaseID: id, Version: patch.Version, Reason: reason, ChangedBy: actor.ID, Snapshot: before, Diff: diff, CreatedAt: s.clock.Now()}); err != nil {
		return err
	}
	return s.audit(ctx, actor, "case.versioned", id, meta)
}
func (s *CaseService) Transition(ctx context.Context, actor domain.User, id string, to domain.CaseStatus, reason string, meta RequestMeta) error {
	c, err := s.repo.CaseByID(ctx, id)
	if err != nil {
		return err
	}
	if !c.CanTransition(to, actor) {
		return domain.ErrInvalidState
	}
	if to == domain.CaseDestroyed && strings.TrimSpace(reason) == "" {
		return domain.ErrValidation
	}
	c.Status = to
	c.DestroyedReason = reason
	c.UpdatedAt = s.clock.Now()
	ms, _ := s.repo.Materials(ctx, id)
	if err := s.repo.UpdateCase(ctx, c, ms); err != nil {
		return err
	}
	return s.audit(ctx, actor, "case.status_changed", id, meta)
}
func (s *CaseService) Search(ctx context.Context, filter CaseFilter) ([]domain.CaseFile, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return s.repo.ListCases(ctx, filter)
}
func (s *CaseService) ExportCSV(ctx context.Context, actor domain.User, filter CaseFilter, w io.Writer, meta RequestMeta) error {
	cases, _, err := s.Search(ctx, filter)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"case_number", "title", "status", "year", "version"}); err != nil {
		return err
	}
	for _, c := range cases {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := cw.Write([]string{c.CaseNumber, c.Title, string(c.Status), strconv.Itoa(c.Year), strconv.Itoa(c.Version)}); err != nil {
			return err
		}
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}
	return s.audit(ctx, actor, "case.exported", "search", meta)
}
func (s *CaseService) audit(ctx context.Context, actor domain.User, action, id string, meta RequestMeta) error {
	return s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: action, Resource: "case", ResourceID: id, RequestID: meta.RequestID, CreatedAt: s.clock.Now()})
}

var _ = time.Second
