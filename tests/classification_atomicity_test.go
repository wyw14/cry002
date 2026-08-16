package tests

import (
	"context"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	"github.com/wyw14/cry002/internal/repository/memory"
	"testing"
	"time"
)

func TestClassificationMoveFailurePreservesTree(t *testing.T) {
	ctx := context.Background()
	repo := memory.New()
	now := time.Now()
	for _, n := range []domain.Classification{{ID: "root", Code: "R", Name: "根", Enabled: true}, {ID: "child", Code: "C", Name: "子", ParentID: "root", Enabled: true}} {
		n.CreatedAt = now
		n.UpdatedAt = now
		if err := repo.CreateClassification(ctx, n); err != nil {
			t.Fatal(err)
		}
	}
	repo.SetClassificationFailure(true)
	svc := application.NewClassificationService(repo, fixedClock{now}, &seqIDs{})
	if err := svc.Move(ctx, actor("arch", domain.RoleArchivist), "child", ""); err == nil {
		t.Fatal("expected injected move failure")
	}
	got, err := repo.ClassificationByID(ctx, "child")
	if err != nil {
		t.Fatal(err)
	}
	if got.ParentID != "root" {
		t.Fatalf("partial update left parent=%q", got.ParentID)
	}
}
