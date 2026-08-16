package postgres

import "time"

type userRow struct {
	ID           string `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string
	DisplayName  string
	AvatarURL    string
	DepartmentID string `gorm:"index"`
	Role         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type departmentRow struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Code      string `gorm:"uniqueIndex"`
	Enabled   bool
	CreatedAt time.Time
}
type classificationRow struct {
	ID        string `gorm:"primaryKey"`
	ParentID  string `gorm:"index"`
	Name      string
	Code      string `gorm:"uniqueIndex"`
	SortOrder int
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
type caseRow struct {
	ID                    string `gorm:"primaryKey"`
	CaseNumber            string `gorm:"uniqueIndex;not null"`
	Title                 string
	ClassificationID      string `gorm:"index"`
	SecurityLevel         string
	Retention             string
	ResponsibleDepartment string `gorm:"index"`
	HandlerID             string `gorm:"index"`
	Year                  int    `gorm:"index"`
	KeywordsJSON          string `gorm:"type:text"`
	Summary               string `gorm:"type:text"`
	Status                string `gorm:"index"`
	Version               int
	CreatedBy             string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DestroyedReason       string
}
type materialRow struct {
	ID           string `gorm:"primaryKey"`
	CaseID       string `gorm:"index;not null"`
	Title        string
	PageFrom     int
	PageTo       int
	DocumentAt   time.Time
	Responsible  string
	AttachmentID string
	CreatedAt    time.Time
}
type versionRow struct {
	ID           string `gorm:"primaryKey"`
	CaseID       string `gorm:"index;not null"`
	Version      int
	Reason       string
	ChangedBy    string
	SnapshotJSON string `gorm:"type:text"`
	DiffJSON     string `gorm:"type:text"`
	CreatedAt    time.Time
}
type reviewRow struct {
	ID         string `gorm:"primaryKey"`
	CaseID     string `gorm:"index;not null"`
	ReviewerID string
	Approved   bool
	Opinion    string
	CreatedAt  time.Time
}
type borrowRow struct {
	ID           string `gorm:"primaryKey"`
	CaseID       string `gorm:"index;not null"`
	ApplicantID  string `gorm:"index"`
	ApproverID   string
	Purpose      string
	Status       string `gorm:"index"`
	DueAt        time.Time
	CheckedOutAt *time.Time
	ReturnedAt   *time.Time
	Renewals     int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type attachmentRow struct {
	ID           string `gorm:"primaryKey"`
	CaseID       string `gorm:"index;not null"`
	MaterialID   string
	OriginalName string
	StorageKey   string `gorm:"uniqueIndex;not null"`
	ContentType  string
	Size         int64
	SHA256       string
	Version      int
	UploadedBy   string
	CreatedAt    time.Time
}
type auditRow struct {
	ID           string `gorm:"primaryKey"`
	ActorID      string
	Action       string `gorm:"index"`
	Resource     string `gorm:"index:idx_resource"`
	ResourceID   string `gorm:"index:idx_resource"`
	Reason       string
	RequestID    string    `gorm:"index"`
	MetadataJSON string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"index"`
}
type tokenRow struct {
	Hash      string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	FamilyID  string `gorm:"index"`
	ExpiresAt time.Time
	UsedAt    *time.Time
	RevokedAt *time.Time
}
