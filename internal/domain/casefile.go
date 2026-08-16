package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"time"
)

type CaseStatus string

const (
	CaseDraft     CaseStatus = "draft"
	CasePending   CaseStatus = "pending_review"
	CaseArchived  CaseStatus = "archived"
	CaseBorrowed  CaseStatus = "borrowed"
	CaseSealed    CaseStatus = "sealed"
	CaseDestroyed CaseStatus = "destroyed"
)

type SecurityLevel string

const (
	SecurityPublic       SecurityLevel = "public"
	SecurityInternal     SecurityLevel = "internal"
	SecurityConfidential SecurityLevel = "confidential"
	SecuritySecret       SecurityLevel = "secret"
)

type RetentionPeriod string

const (
	RetentionTenYears    RetentionPeriod = "10_years"
	RetentionThirtyYears RetentionPeriod = "30_years"
	RetentionPermanent   RetentionPeriod = "permanent"
)

type CaseFile struct {
	ID                    string          `json:"id"`
	CaseNumber            string          `json:"case_number"`
	Title                 string          `json:"title"`
	ClassificationID      string          `json:"classification_id"`
	SecurityLevel         SecurityLevel   `json:"security_level"`
	Retention             RetentionPeriod `json:"retention"`
	ResponsibleDepartment string          `json:"responsible_department"`
	HandlerID             string          `json:"handler_id"`
	Year                  int             `json:"year"`
	Keywords              []string        `json:"keywords"`
	Summary               string          `json:"summary"`
	Status                CaseStatus      `json:"status"`
	Version               int             `json:"version"`
	CreatedBy             string          `json:"created_by"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	DestroyedReason       string          `json:"destroyed_reason,omitempty"`
}
type Material struct {
	ID           string    `json:"id"`
	CaseID       string    `json:"case_id"`
	Title        string    `json:"title"`
	PageFrom     int       `json:"page_from"`
	PageTo       int       `json:"page_to"`
	DocumentAt   time.Time `json:"document_at"`
	Responsible  string    `json:"responsible"`
	AttachmentID string    `json:"attachment_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
type CaseSnapshot struct {
	Case      CaseFile   `json:"case"`
	Materials []Material `json:"materials"`
}
type CaseVersion struct {
	ID        string       `json:"id"`
	CaseID    string       `json:"case_id"`
	Version   int          `json:"version"`
	Reason    string       `json:"reason"`
	ChangedBy string       `json:"changed_by"`
	Snapshot  CaseSnapshot `json:"snapshot"`
	Diff      []string     `json:"diff"`
	CreatedAt time.Time    `json:"created_at"`
}
type Review struct {
	ID         string    `json:"id"`
	CaseID     string    `json:"case_id"`
	ReviewerID string    `json:"reviewer_id"`
	Approved   bool      `json:"approved"`
	Opinion    string    `json:"opinion"`
	CreatedAt  time.Time `json:"created_at"`
}

func (c CaseFile) ValidateForSubmit(materials []Material) error {
	if strings.TrimSpace(c.CaseNumber) == "" || strings.TrimSpace(c.Title) == "" || c.ClassificationID == "" || c.ResponsibleDepartment == "" || c.HandlerID == "" || c.Year < 1900 || len(materials) == 0 {
		return ErrValidation
	}
	for _, m := range materials {
		if m.PageFrom <= 0 || m.PageTo < m.PageFrom || strings.TrimSpace(m.Title) == "" {
			return ErrValidation
		}
	}
	return nil
}
func (c CaseFile) CanTransition(to CaseStatus, actor User) bool {
	switch c.Status {
	case CaseDraft:
		return to == CasePending && (actor.ID == c.CreatedBy || actor.Role == RoleArchivist)
	case CasePending:
		return (to == CaseDraft || to == CaseArchived) && (actor.Role == RoleArchivist || actor.Role == RoleAuditor || actor.IsAdministrator())
	case CaseArchived:
		return (to == CaseBorrowed && (actor.Role == RoleArchivist || actor.Role == RoleBorrower || actor.IsAdministrator())) || (to == CaseSealed && actor.IsAdministrator())
	case CaseBorrowed:
		return to == CaseArchived && (actor.Role == RoleArchivist || actor.IsAdministrator())
	case CaseSealed:
		return to == CaseDestroyed && actor.IsAdministrator()
	case CaseDestroyed:
		return to == CaseSealed && actor.IsAdministrator()
	default:
		return false
	}
}
func (s CaseSnapshot) Digest() string {
	parts := []string{s.Case.CaseNumber, s.Case.Title, s.Case.Summary, string(s.Case.Status)}
	for _, m := range s.Materials {
		parts = append(parts, m.ID, m.Title)
	}
	sort.Strings(parts)
	h := sha256.Sum256([]byte(strings.Join(parts, "\x1f")))
	return hex.EncodeToString(h[:])
}
