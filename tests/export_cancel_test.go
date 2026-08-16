package tests

import (
	"bytes"
	"context"
	"errors"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	"github.com/wyw14/cry002/internal/repository/memory"
	"testing"
	"time"
)

func TestCanceledExportLeavesNoAudit(t *testing.T) {
	repo := memory.New()
	repo.SetDelay(100 * time.Millisecond)
	svc := application.NewCaseService(repo, fixedClock{time.Now()}, &seqIDs{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	err := svc.ExportCSV(ctx, actor("u", domain.RoleArchivist), application.CaseFilter{}, &bytes.Buffer{}, application.RequestMeta{RequestID: "cancel"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err=%v", err)
	}
	audits, _ := repo.ListAudits(context.Background(), "case", "")
	if len(audits) != 0 {
		t.Fatalf("canceled export wrote %d audits", len(audits))
	}
}
