package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func classToRow(n domain.Classification) classificationRow {
	return classificationRow{n.ID, n.ParentID, n.Name, n.Code, n.SortOrder, n.Enabled, n.CreatedAt, n.UpdatedAt}
}
func rowToClass(x classificationRow) domain.Classification {
	return domain.Classification{ID: x.ID, ParentID: x.ParentID, Name: x.Name, Code: x.Code, SortOrder: x.SortOrder, Enabled: x.Enabled, CreatedAt: x.CreatedAt, UpdatedAt: x.UpdatedAt}
}
func (r *Repository) CreateClassification(ctx context.Context, n domain.Classification) error {
	return r.db.WithContext(ctx).Create(&[]classificationRow{classToRow(n)}).Error
}
func (r *Repository) UpdateClassification(ctx context.Context, n domain.Classification) error {
	return r.db.WithContext(ctx).Save(&[]classificationRow{classToRow(n)}).Error
}
func (r *Repository) ClassificationByID(ctx context.Context, id string) (domain.Classification, error) {
	var x classificationRow
	err := r.db.WithContext(ctx).First(&x, "id = ?", id).Error
	return rowToClass(x), dbError(err)
}
func (r *Repository) ListClassifications(ctx context.Context) ([]domain.Classification, error) {
	var xs []classificationRow
	if err := r.db.WithContext(ctx).Order("sort_order,id").Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Classification, len(xs))
	for i, x := range xs {
		out[i] = rowToClass(x)
	}
	return out, nil
}
func (r *Repository) MoveClassification(ctx context.Context, id, parent string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var n classificationRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&n, "id = ?", id).Error; err != nil {
			return dbError(err)
		}
		n.ParentID = parent
		n.UpdatedAt = time.Now()
		return tx.Save(&n).Error
	})
}
