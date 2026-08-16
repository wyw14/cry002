package postgres

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
)

func userToRow(u domain.User) userRow {
	return userRow{u.ID, u.Email, u.PasswordHash, u.DisplayName, u.AvatarURL, u.DepartmentID, string(u.Role), string(u.Status), u.CreatedAt, u.UpdatedAt}
}
func rowToUser(x userRow) domain.User {
	return domain.User{ID: x.ID, Email: x.Email, PasswordHash: x.PasswordHash, DisplayName: x.DisplayName, AvatarURL: x.AvatarURL, DepartmentID: x.DepartmentID, Role: domain.Role(x.Role), Status: domain.UserStatus(x.Status), CreatedAt: x.CreatedAt, UpdatedAt: x.UpdatedAt}
}
func (r *Repository) CreateUser(ctx context.Context, u domain.User) error {
	return r.db.WithContext(ctx).Create(&[]userRow{userToRow(u)}).Error
}
func (r *Repository) UpdateUser(ctx context.Context, u domain.User) error {
	return r.db.WithContext(ctx).Save(&[]userRow{userToRow(u)}).Error
}
func (r *Repository) UserByID(ctx context.Context, id string) (domain.User, error) {
	var x userRow
	err := r.db.WithContext(ctx).First(&x, "id = ?", id).Error
	return rowToUser(x), dbError(err)
}
func (r *Repository) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	var x userRow
	err := r.db.WithContext(ctx).First(&x, "email = ?", email).Error
	return rowToUser(x), dbError(err)
}
func (r *Repository) ListUsers(ctx context.Context) ([]domain.User, error) {
	var xs []userRow
	if err := r.db.WithContext(ctx).Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.User, len(xs))
	for i, x := range xs {
		out[i] = rowToUser(x)
	}
	return out, nil
}
func (r *Repository) CreateDepartment(ctx context.Context, d domain.Department) error {
	return r.db.WithContext(ctx).Create(&departmentRow{d.ID, d.Name, d.Code, d.Enabled, d.CreatedAt}).Error
}
func (r *Repository) ListDepartments(ctx context.Context) ([]domain.Department, error) {
	var xs []departmentRow
	if err := r.db.WithContext(ctx).Find(&xs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Department, len(xs))
	for i, x := range xs {
		out[i] = domain.Department{ID: x.ID, Name: x.Name, Code: x.Code, Enabled: x.Enabled, CreatedAt: x.CreatedAt}
	}
	return out, nil
}
