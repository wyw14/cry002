package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

func (r *Repository) CreateVersion(ctx context.Context, v domain.CaseVersion) error {
	x := versionRow{ID: v.ID, CaseID: v.CaseID, Version: v.Version, Reason: v.Reason, ChangedBy: v.ChangedBy, SnapshotJSON: jsonText(v.Snapshot), DiffJSON: jsonText(v.Diff), CreatedAt: v.CreatedAt}
	return r.db.WithContext(ctx).Create(&x).Error
}
func (r *Repository) Versions(ctx context.Context, id string) ([]domain.CaseVersion, error) {
	var xs []versionRow
	if err := r.db.WithContext(ctx).Where("case_id = ?", id).Order("version DESC").Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CaseVersion, len(xs))
	for i, x := range xs {
		var s domain.CaseSnapshot
		var d []string
		fromJSON(x.SnapshotJSON, &s)
		fromJSON(x.DiffJSON, &d)
		out[i] = domain.CaseVersion{ID: x.ID, CaseID: x.CaseID, Version: x.Version, Reason: x.Reason, ChangedBy: x.ChangedBy, Snapshot: s, Diff: d, CreatedAt: x.CreatedAt}
	}
	return out, nil
}
func (r *Repository) CreateReview(ctx context.Context, v domain.Review) error {
	x := reviewRow{v.ID, v.CaseID, v.ReviewerID, v.Approved, v.Opinion, v.CreatedAt}
	return r.db.WithContext(ctx).Create(&x).Error
}
func (r *Repository) Reviews(ctx context.Context, id string) ([]domain.Review, error) {
	var xs []reviewRow
	if err := r.db.WithContext(ctx).Where("case_id = ?", id).Order("created_at").Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Review, len(xs))
	for i, x := range xs {
		out[i] = domain.Review{ID: x.ID, CaseID: x.CaseID, ReviewerID: x.ReviewerID, Approved: x.Approved, Opinion: x.Opinion, CreatedAt: x.CreatedAt}
	}
	return out, nil
}
