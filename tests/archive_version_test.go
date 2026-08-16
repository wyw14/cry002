package tests

import (
	"context"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	"github.com/wyw14/cry002/internal/repository/memory"
	"testing"
	"time"
)

func TestArchivedUpdateCreatesVersionAndPreservesSnapshot(t *testing.T) {
	ctx := context.Background()
	repo := memory.New()
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	old := domain.CaseFile{ID: "case", CaseNumber: "A-1", Title: "旧标题", Summary: "旧摘要", Status: domain.CaseArchived, Version: 1, CreatedBy: "owner", CreatedAt: now, UpdatedAt: now}
	ms := []domain.Material{{ID: "m1", CaseID: "case", Title: "原材料", PageFrom: 1, PageTo: 2}}
	if err := repo.CreateCase(ctx, old, ms); err != nil {
		t.Fatal(err)
	}
	patch := old
	patch.Title = "新标题"
	svc := application.NewCaseService(repo, fixedClock{now.Add(time.Hour)}, &seqIDs{})
	if err := svc.UpdateMetadata(ctx, actor("arch", domain.RoleArchivist), "case", "勘误", patch, application.RequestMeta{}); err != nil {
		t.Fatal(err)
	}
	got, _ := repo.CaseByID(ctx, "case")
	if got.Version != 2 || got.Title != "新标题" {
		t.Fatalf("unexpected updated case: %+v", got)
	}
	versions, _ := repo.Versions(ctx, "case")
	if len(versions) != 1 || versions[0].Snapshot.Case.Title != "旧标题" || len(versions[0].Diff) == 0 {
		t.Fatalf("version did not preserve old snapshot: %+v", versions)
	}
}
