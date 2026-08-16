package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

func (r *Repository) AppendAudit(ctx context.Context, a domain.AuditEvent) error {
	x := auditRow{ID: a.ID, ActorID: a.ActorID, Action: a.Action, Resource: a.Resource, ResourceID: a.ResourceID, Reason: a.Reason, RequestID: a.RequestID, MetadataJSON: jsonText(a.Metadata), CreatedAt: a.CreatedAt}
	return r.db.WithContext(ctx).Create(&x).Error
}
func (r *Repository) ListAudits(ctx context.Context, res, id string) ([]domain.AuditEvent, error) {
	q := r.db.WithContext(ctx)
	if res != "" {
		q = q.Where("resource = ?", res)
	}
	if id != "" {
		q = q.Where("resource_id = ?", id)
	}
	var xs []auditRow
	if err := q.Order("created_at DESC").Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.AuditEvent, len(xs))
	for i, x := range xs {
		var m map[string]any
		fromJSON(x.MetadataJSON, &m)
		out[i] = domain.AuditEvent{ID: x.ID, ActorID: x.ActorID, Action: x.Action, Resource: x.Resource, ResourceID: x.ResourceID, Reason: x.Reason, RequestID: x.RequestID, Metadata: m, CreatedAt: x.CreatedAt}
	}
	return out, nil
}
