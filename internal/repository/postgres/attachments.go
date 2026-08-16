package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

func attachmentToRow(a domain.Attachment) attachmentRow {
	return attachmentRow{a.ID, a.CaseID, a.MaterialID, a.OriginalName, a.StorageKey, a.ContentType, a.Size, a.SHA256, a.Version, a.UploadedBy, a.CreatedAt}
}
func rowToAttachment(x attachmentRow) domain.Attachment {
	return domain.Attachment{ID: x.ID, CaseID: x.CaseID, MaterialID: x.MaterialID, OriginalName: x.OriginalName, StorageKey: x.StorageKey, ContentType: x.ContentType, Size: x.Size, SHA256: x.SHA256, Version: x.Version, UploadedBy: x.UploadedBy, CreatedAt: x.CreatedAt}
}
func (r *Repository) CreateAttachment(ctx context.Context, a domain.Attachment) error {
	x := attachmentToRow(a)
	return r.db.WithContext(ctx).Create(&x).Error
}
func (r *Repository) Attachments(ctx context.Context, id string) ([]domain.Attachment, error) {
	var xs []attachmentRow
	if err := r.db.WithContext(ctx).Where("case_id = ?", id).Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Attachment, len(xs))
	for i, x := range xs {
		out[i] = rowToAttachment(x)
	}
	return out, nil
}
func (r *Repository) AttachmentByID(ctx context.Context, id string) (domain.Attachment, error) {
	var x attachmentRow
	err := r.db.WithContext(ctx).First(&x, "id = ?", id).Error
	return rowToAttachment(x), dbError(err)
}
