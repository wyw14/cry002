package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	"gorm.io/gorm"
	"strings"
)

func caseToRow(c domain.CaseFile) caseRow {
	return caseRow{ID: c.ID, CaseNumber: c.CaseNumber, Title: c.Title, ClassificationID: c.ClassificationID, SecurityLevel: string(c.SecurityLevel), Retention: string(c.Retention), ResponsibleDepartment: c.ResponsibleDepartment, HandlerID: c.HandlerID, Year: c.Year, KeywordsJSON: jsonText(c.Keywords), Summary: c.Summary, Status: string(c.Status), Version: c.Version, CreatedBy: c.CreatedBy, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt, DestroyedReason: c.DestroyedReason}
}
func rowToCase(x caseRow) domain.CaseFile {
	var k []string
	fromJSON(x.KeywordsJSON, &k)
	return domain.CaseFile{ID: x.ID, CaseNumber: x.CaseNumber, Title: x.Title, ClassificationID: x.ClassificationID, SecurityLevel: domain.SecurityLevel(x.SecurityLevel), Retention: domain.RetentionPeriod(x.Retention), ResponsibleDepartment: x.ResponsibleDepartment, HandlerID: x.HandlerID, Year: x.Year, Keywords: k, Summary: x.Summary, Status: domain.CaseStatus(x.Status), Version: x.Version, CreatedBy: x.CreatedBy, CreatedAt: x.CreatedAt, UpdatedAt: x.UpdatedAt, DestroyedReason: x.DestroyedReason}
}
func materialToRow(m domain.Material) materialRow {
	return materialRow{m.ID, m.CaseID, m.Title, m.PageFrom, m.PageTo, m.DocumentAt, m.Responsible, m.AttachmentID, m.CreatedAt}
}
func rowToMaterial(x materialRow) domain.Material {
	return domain.Material{ID: x.ID, CaseID: x.CaseID, Title: x.Title, PageFrom: x.PageFrom, PageTo: x.PageTo, DocumentAt: x.DocumentAt, Responsible: x.Responsible, AttachmentID: x.AttachmentID, CreatedAt: x.CreatedAt}
}
func (r *Repository) CreateCase(ctx context.Context, c domain.CaseFile, ms []domain.Material) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		x := caseToRow(c)
		if err := tx.Create(&x).Error; err != nil {
			return err
		}
		rows := make([]materialRow, len(ms))
		for i, m := range ms {
			rows[i] = materialToRow(m)
		}
		if len(rows) > 0 {
			return tx.Create(&rows).Error
		}
		return nil
	})
}
func (r *Repository) UpdateCase(ctx context.Context, c domain.CaseFile, ms []domain.Material) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		x := caseToRow(c)
		if err := tx.Save(&x).Error; err != nil {
			return err
		}
		if err := tx.Where("case_id = ?", c.ID).Delete(&materialRow{}).Error; err != nil {
			return err
		}
		rows := make([]materialRow, len(ms))
		for i, m := range ms {
			rows[i] = materialToRow(m)
		}
		if len(rows) > 0 {
			return tx.Create(&rows).Error
		}
		return nil
	})
}
func (r *Repository) CaseByID(ctx context.Context, id string) (domain.CaseFile, error) {
	var x caseRow
	err := r.db.WithContext(ctx).First(&x, "id = ?", id).Error
	return rowToCase(x), dbError(err)
}
func (r *Repository) ListCases(ctx context.Context, f application.CaseFilter) ([]domain.CaseFile, int, error) {
	q := r.db.WithContext(context.WithoutCancel(ctx)).Model(&caseRow{})
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.ClassificationID != "" {
		q = q.Where("classification_id = ?", f.ClassificationID)
	}
	if f.Department != "" {
		q = q.Where("responsible_department = ?", f.Department)
	}
	if f.Year != 0 {
		q = q.Where("year = ?", f.Year)
	}
	if strings.TrimSpace(f.Query) != "" {
		like := "%" + strings.TrimSpace(f.Query) + "%"
		q = q.Where("case_number ILIKE ? OR title ILIKE ? OR summary ILIKE ?", like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, size := f.Page, f.PageSize
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	var xs []caseRow
	if err := q.Order("updated_at DESC").Offset((page - 1) * size).Limit(size).Find(&xs).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.CaseFile, len(xs))
	for i, x := range xs {
		out[i] = rowToCase(x)
	}
	return out, int(total), nil
}
func (r *Repository) Materials(ctx context.Context, id string) ([]domain.Material, error) {
	var xs []materialRow
	if err := r.db.WithContext(ctx).Where("case_id = ?", id).Order("page_from,id").Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Material, len(xs))
	for i, x := range xs {
		out[i] = rowToMaterial(x)
	}
	return out, nil
}
