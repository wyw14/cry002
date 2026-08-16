package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func tokenToRow(t domain.RefreshToken) tokenRow {
	return tokenRow{t.Hash, t.UserID, t.FamilyID, t.ExpiresAt, t.UsedAt, t.RevokedAt}
}
func rowToToken(x tokenRow) domain.RefreshToken {
	return domain.RefreshToken{Hash: x.Hash, UserID: x.UserID, FamilyID: x.FamilyID, ExpiresAt: x.ExpiresAt, UsedAt: x.UsedAt, RevokedAt: x.RevokedAt}
}
func (r *Repository) StoreRefreshToken(ctx context.Context, t domain.RefreshToken) error {
	x := tokenToRow(t)
	return r.db.WithContext(ctx).Create(&x).Error
}
func (r *Repository) RefreshToken(ctx context.Context, h string) (domain.RefreshToken, error) {
	var x tokenRow
	err := r.db.WithContext(ctx).First(&x, "hash = ?", h).Error
	return rowToToken(x), dbError(err)
}
func (r *Repository) ConsumeRefreshToken(ctx context.Context, h string, now time.Time) (domain.RefreshToken, error) {
	var out domain.RefreshToken
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var x tokenRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&x, "hash = ?", h).Error; err != nil {
			return dbError(err)
		}
		t := rowToToken(x)
		if err := t.CanRotate(now); err != nil {
			return err
		}
		t.UsedAt = &now
		x = tokenToRow(t)
		out = t
		return tx.Save(&x).Error
	})
	return out, err
}
func (r *Repository) RevokeTokenFamily(ctx context.Context, f string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&tokenRow{}).Where("family_id = ?", f).Update("revoked_at", now).Error
}

var _ interface {
	ConsumeRefreshToken(context.Context, string, time.Time) (domain.RefreshToken, error)
} = (*Repository)(nil)
var _ *gorm.DB
