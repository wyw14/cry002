package domain

import "time"

type BorrowStatus string

const (
	BorrowPending    BorrowStatus = "pending"
	BorrowApproved   BorrowStatus = "approved"
	BorrowCheckedOut BorrowStatus = "checked_out"
	BorrowReturned   BorrowStatus = "returned"
	BorrowRejected   BorrowStatus = "rejected"
	BorrowOverdue    BorrowStatus = "overdue"
)

type BorrowRequest struct {
	ID           string       `json:"id"`
	CaseID       string       `json:"case_id"`
	ApplicantID  string       `json:"applicant_id"`
	ApproverID   string       `json:"approver_id,omitempty"`
	Purpose      string       `json:"purpose"`
	Status       BorrowStatus `json:"status"`
	DueAt        time.Time    `json:"due_at"`
	CheckedOutAt *time.Time   `json:"checked_out_at,omitempty"`
	ReturnedAt   *time.Time   `json:"returned_at,omitempty"`
	Renewals     int          `json:"renewals"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (b BorrowRequest) CanApprove(actor User) bool {
	return b.Status == BorrowPending && (actor.Role == RoleArchivist || actor.IsAdministrator())
}
func (b BorrowRequest) CanRenew(now time.Time) bool {
	return (b.Status == BorrowCheckedOut || b.Status == BorrowOverdue) && b.Renewals < 2 && b.ReturnedAt == nil && now.Before(b.DueAt.Add(30*24*time.Hour))
}
