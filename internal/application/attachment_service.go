package application

import (
	"context"
	"fmt"
	"github.com/wyw14/cry002/internal/domain"
	"io"
	"path/filepath"
	"strings"
	"unicode"
)

type AttachmentService struct {
	repo    Repository
	storage Storage
	clock   Clock
	ids     IDGenerator
}

func NewAttachmentService(r Repository, s Storage, c Clock, i IDGenerator) *AttachmentService {
	return &AttachmentService{r, s, c, i}
}
func (s *AttachmentService) Upload(ctx context.Context, actor domain.User, caseID, materialID, name, contentType string, size int64, r io.Reader, meta RequestMeta) (domain.Attachment, error) {
	if size <= 0 || size > 50*1024*1024 {
		return domain.Attachment{}, domain.ErrValidation
	}
	if !allowedType(contentType) {
		return domain.Attachment{}, domain.ErrValidation
	}
	base, ok := safeAttachmentName(name)
	if !ok {
		return domain.Attachment{}, domain.ErrValidation
	}
	c, err := s.repo.CaseByID(ctx, caseID)
	if err != nil {
		return domain.Attachment{}, err
	}
	if !canAccessCase(actor, c) {
		return domain.Attachment{}, domain.ErrForbidden
	}
	key := fmt.Sprintf("%s/%s/%s", caseID, s.ids.New(), base)
	storageKey, sum, err := s.storage.Put(ctx, key, r, size)
	if err != nil {
		return domain.Attachment{}, err
	}
	a := domain.Attachment{ID: s.ids.New(), CaseID: caseID, MaterialID: materialID, OriginalName: base, StorageKey: storageKey, ContentType: contentType, Size: size, SHA256: sum, Version: 1, UploadedBy: actor.ID, CreatedAt: s.clock.Now()}
	if err := s.repo.CreateAttachment(ctx, a); err != nil {
		_ = s.storage.Remove(ctx, storageKey)
		return domain.Attachment{}, err
	}
	_ = s.repo.AppendAudit(ctx, domain.AuditEvent{ID: s.ids.New(), ActorID: actor.ID, Action: "attachment.uploaded", Resource: "case", ResourceID: caseID, RequestID: meta.RequestID, CreatedAt: s.clock.Now()})
	return a, nil
}
func (s *AttachmentService) Open(ctx context.Context, actor domain.User, id string) (domain.Attachment, io.ReadCloser, error) {
	a, err := s.repo.AttachmentByID(ctx, id)
	if err != nil {
		return domain.Attachment{}, nil, err
	}
	c, err := s.repo.CaseByID(ctx, a.CaseID)
	if err != nil {
		return domain.Attachment{}, nil, err
	}
	if !canAccessCase(actor, c) && !s.hasActiveBorrow(ctx, actor, a.CaseID) {
		return domain.Attachment{}, nil, domain.ErrForbidden
	}
	r, err := s.storage.Open(ctx, a.StorageKey)
	return a, r, err
}
func (s *AttachmentService) hasActiveBorrow(ctx context.Context, actor domain.User, caseID string) bool {
	if actor.ID == "" {
		return false
	}
	borrows, err := s.repo.ListBorrows(ctx, string(domain.BorrowCheckedOut))
	if err != nil {
		return false
	}
	for _, b := range borrows {
		if b.CaseID == caseID && b.ApplicantID == actor.ID {
			return true
		}
	}
	overdue, err := s.repo.ListBorrows(ctx, string(domain.BorrowOverdue))
	if err != nil {
		return false
	}
	for _, b := range overdue {
		if b.CaseID == caseID && b.ApplicantID == actor.ID {
			return true
		}
	}
	return false
}
func canAccessCase(a domain.User, c domain.CaseFile) bool {
	return a.IsAdministrator() || a.Role == domain.RoleArchivist || a.Role == domain.RoleAuditor || a.ID == c.CreatedBy || a.ID == c.HandlerID || (a.DepartmentID != "" && a.DepartmentID == c.ResponsibleDepartment)
}
func allowedType(t string) bool {
	switch strings.ToLower(t) {
	case "application/pdf", "image/jpeg", "image/png", "text/plain":
		return true
	default:
		return false
	}
}
func safeAttachmentName(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	for _, r := range name {
		if r == unicode.ReplacementChar {
			return "", false
		}
	}
	base := filepath.Base(filepath.FromSlash(name))
	if base == "." || base == ".." || base == string(filepath.Separator) || strings.TrimSpace(base) == "" {
		return "", false
	}
	for _, r := range base {
		if unicode.IsControl(r) {
			return "", false
		}
	}
	return base, true
}
