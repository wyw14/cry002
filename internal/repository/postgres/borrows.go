package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

func borrowToRow(b domain.BorrowRequest) borrowRow {
	return borrowRow{b.ID, b.CaseID, b.ApplicantID, b.ApproverID, b.Purpose, string(b.Status), b.DueAt, b.CheckedOutAt, b.ReturnedAt, b.Renewals, b.CreatedAt, b.UpdatedAt}
}
func rowToBorrow(x borrowRow) domain.BorrowRequest {
	return domain.BorrowRequest{ID: x.ID, CaseID: x.CaseID, ApplicantID: x.ApplicantID, ApproverID: x.ApproverID, Purpose: x.Purpose, Status: domain.BorrowStatus(x.Status), DueAt: x.DueAt, CheckedOutAt: x.CheckedOutAt, ReturnedAt: x.ReturnedAt, Renewals: x.Renewals, CreatedAt: x.CreatedAt, UpdatedAt: x.UpdatedAt}
}
func (r *Repository) CreateBorrow(ctx context.Context, b domain.BorrowRequest) error {
	x := borrowToRow(b)
	return r.db.WithContext(ctx).Create(&x).Error
}
func (r *Repository) BorrowByID(ctx context.Context, id string) (domain.BorrowRequest, error) {
	var x borrowRow
	err := r.db.WithContext(ctx).First(&x, "id = ?", id).Error
	return rowToBorrow(x), dbError(err)
}
func (r *Repository) ListBorrows(ctx context.Context, status string) ([]domain.BorrowRequest, error) {
	q := r.db.WithContext(ctx)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var xs []borrowRow
	if err := q.Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.BorrowRequest, len(xs))
	for i, x := range xs {
		out[i] = rowToBorrow(x)
	}
	return out, nil
}
func (r *Repository) UpdateBorrow(ctx context.Context, b domain.BorrowRequest) error {
	x := borrowToRow(b)
	return r.db.WithContext(ctx).Save(&x).Error
}
func (r *Repository) CheckoutCase(ctx context.Context, cid, bid string, now time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var c caseRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&c, "id = ?", cid).Error; err != nil {
			return dbError(err)
		}
		var b borrowRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&b, "id = ?", bid).Error; err != nil {
			return dbError(err)
		}
		if c.Status != string(domain.CaseArchived) || b.CaseID != cid || b.Status != string(domain.BorrowApproved) {
			return domain.ErrAlreadyBorrowed
		}
		c.Status = string(domain.CaseBorrowed)
		c.UpdatedAt = now
		b.Status = string(domain.BorrowCheckedOut)
		b.CheckedOutAt = &now
		b.UpdatedAt = now
		if err := tx.Save(&c).Error; err != nil {
			return translateCheckoutErr(err)
		}
		if err := tx.Save(&b).Error; err != nil {
			return translateCheckoutErr(err)
		}
		return nil
	})
}
func translateCheckoutErr(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "uq_active_checkout") || strings.Contains(msg, "unique constraint") {
		return domain.ErrAlreadyBorrowed
	}
	return err
}
