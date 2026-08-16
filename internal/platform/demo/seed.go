package demo

import (
	"context"
	"errors"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	"time"
)

func Seed(ctx context.Context, r application.Repository, h application.PasswordHasher, now time.Time) error {
	hash, err := h.Hash("Demo123456!")
	if err != nil {
		return err
	}
	deps := []domain.Department{{ID: "dept-archives", Name: "档案管理部", Code: "ARCH", Enabled: true, CreatedAt: now}, {ID: "dept-legal", Name: "法务部", Code: "LEGAL", Enabled: true, CreatedAt: now}}
	for _, d := range deps {
		if err := r.CreateDepartment(ctx, d); err != nil && !errors.Is(err, domain.ErrConflict) {
			return err
		}
	}
	users := []domain.User{{ID: "user-admin", Email: "admin@archive.local", DisplayName: "平台管理员", Role: domain.RoleAdministrator, DepartmentID: "dept-archives"}, {ID: "user-archivist", Email: "archivist@archive.local", DisplayName: "档案管理员", Role: domain.RoleArchivist, DepartmentID: "dept-archives"}, {ID: "user-auditor", Email: "auditor@archive.local", DisplayName: "审核员", Role: domain.RoleAuditor, DepartmentID: "dept-archives"}, {ID: "user-manager", Email: "manager@archive.local", DisplayName: "部门负责人", Role: domain.RoleDepartment, DepartmentID: "dept-legal"}, {ID: "user-borrower", Email: "borrower@archive.local", DisplayName: "借阅人", Role: domain.RoleBorrower, DepartmentID: "dept-legal"}}
	for _, u := range users {
		u.PasswordHash = hash
		u.Status = domain.UserActive
		u.CreatedAt = now
		u.UpdatedAt = now
		if err := r.CreateUser(ctx, u); err != nil && !errors.Is(err, domain.ErrConflict) {
			return err
		}
	}
	classes := []domain.Classification{{ID: "class-contract", Name: "合同档案", Code: "CONTRACT", SortOrder: 1, Enabled: true, CreatedAt: now, UpdatedAt: now}, {ID: "class-project", Name: "项目档案", Code: "PROJECT", SortOrder: 2, Enabled: true, CreatedAt: now, UpdatedAt: now}}
	for _, c := range classes {
		if err := r.CreateClassification(ctx, c); err != nil && !errors.Is(err, domain.ErrConflict) {
			return err
		}
	}
	cases := []domain.CaseFile{{ID: "case-archived", CaseNumber: "ARCH-2026-0001", Title: "年度采购合同", ClassificationID: "class-contract", SecurityLevel: domain.SecurityInternal, Retention: domain.RetentionThirtyYears, ResponsibleDepartment: "dept-legal", HandlerID: "user-manager", Year: 2026, Keywords: []string{"采购", "合同"}, Summary: "已归档演示案卷", Status: domain.CaseArchived, Version: 1, CreatedBy: "user-manager", CreatedAt: now, UpdatedAt: now}, {ID: "case-draft", CaseNumber: "ARCH-2026-0002", Title: "新项目立项材料", ClassificationID: "class-project", SecurityLevel: domain.SecurityPublic, Retention: domain.RetentionTenYears, ResponsibleDepartment: "dept-legal", HandlerID: "user-manager", Year: 2026, Keywords: []string{"项目"}, Summary: "待提交案卷", Status: domain.CaseDraft, Version: 1, CreatedBy: "user-manager", CreatedAt: now, UpdatedAt: now}}
	for _, c := range cases {
		ms := []domain.Material{{ID: "material-" + c.ID, CaseID: c.ID, Title: "案卷首页", PageFrom: 1, PageTo: 3, DocumentAt: now, Responsible: "法务部", CreatedAt: now}}
		if err := r.CreateCase(ctx, c, ms); err != nil && !errors.Is(err, domain.ErrConflict) {
			return err
		}
	}
	return nil
}
