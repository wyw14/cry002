package application

import (
	"context"
	"io"
	"time"

	"github.com/wyw14/cry002/internal/domain"
)

type Repository interface {
	CreateUser(context.Context, domain.User) error
	UpdateUser(context.Context, domain.User) error
	UserByID(context.Context, string) (domain.User, error)
	UserByEmail(context.Context, string) (domain.User, error)
	ListUsers(context.Context) ([]domain.User, error)
	CreateDepartment(context.Context, domain.Department) error
	ListDepartments(context.Context) ([]domain.Department, error)

	CreateClassification(context.Context, domain.Classification) error
	UpdateClassification(context.Context, domain.Classification) error
	ClassificationByID(context.Context, string) (domain.Classification, error)
	ListClassifications(context.Context) ([]domain.Classification, error)
	MoveClassification(context.Context, string, string) error

	CreateCase(context.Context, domain.CaseFile, []domain.Material) error
	UpdateCase(context.Context, domain.CaseFile, []domain.Material) error
	CaseByID(context.Context, string) (domain.CaseFile, error)
	ListCases(context.Context, CaseFilter) ([]domain.CaseFile, int, error)
	Materials(context.Context, string) ([]domain.Material, error)
	CreateVersion(context.Context, domain.CaseVersion) error
	Versions(context.Context, string) ([]domain.CaseVersion, error)
	CreateReview(context.Context, domain.Review) error
	Reviews(context.Context, string) ([]domain.Review, error)

	CreateBorrow(context.Context, domain.BorrowRequest) error
	BorrowByID(context.Context, string) (domain.BorrowRequest, error)
	ListBorrows(context.Context, string) ([]domain.BorrowRequest, error)
	UpdateBorrow(context.Context, domain.BorrowRequest) error
	CheckoutCase(context.Context, string, string, time.Time) error

	CreateAttachment(context.Context, domain.Attachment) error
	Attachments(context.Context, string) ([]domain.Attachment, error)
	AttachmentByID(context.Context, string) (domain.Attachment, error)

	AppendAudit(context.Context, domain.AuditEvent) error
	ListAudits(context.Context, string, string) ([]domain.AuditEvent, error)

	StoreRefreshToken(context.Context, domain.RefreshToken) error
	RefreshToken(context.Context, string) (domain.RefreshToken, error)
	ConsumeRefreshToken(context.Context, string, time.Time) (domain.RefreshToken, error)
	RevokeTokenFamily(context.Context, string, time.Time) error
	DoIdempotent(context.Context, string, func() (string, error)) (string, error)
}

type Clock interface{ Now() time.Time }
type IDGenerator interface{ New() string }
type PasswordHasher interface {
	Hash(string) (string, error)
	Compare(string, string) error
}
type TokenManager interface {
	IssueAccess(domain.User) (string, error)
	ParseAccess(string) (string, error)
	NewRefreshToken() (string, string, string, error)
	HashRefresh(string) string
}
type Storage interface {
	Put(context.Context, string, io.Reader, int64) (string, string, error)
	Open(context.Context, string) (io.ReadCloser, error)
	Remove(context.Context, string) error
}
type RequestMeta struct{ RequestID string }
type CaseFilter struct {
	Query            string
	Status           domain.CaseStatus
	ClassificationID string
	Department       string
	Year             int
	Page             int
	PageSize         int
	Sort             string
}
