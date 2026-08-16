package tests

import (
	"context"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	"github.com/wyw14/cry002/internal/repository/memory"
	"sync"
	"testing"
	"time"
)

func TestConcurrentCheckoutSingleWinner(t *testing.T) {
	ctx := context.Background()
	repo := memory.New()
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	c := domain.CaseFile{ID: "case", CaseNumber: "A-1", Status: domain.CaseArchived, CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateCase(ctx, c, nil); err != nil {
		t.Fatal(err)
	}
	for _, b := range []domain.BorrowRequest{{ID: "b1", CaseID: "case", ApplicantID: "u1", Status: domain.BorrowApproved}, {ID: "b2", CaseID: "case", ApplicantID: "u2", Status: domain.BorrowApproved}} {
		if err := repo.CreateBorrow(ctx, b); err != nil {
			t.Fatal(err)
		}
	}
	svc := application.NewBorrowService(repo, fixedClock{now}, &seqIDs{})
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for _, x := range []struct{ id, user string }{{"b1", "u1"}, {"b2", "u2"}} {
		wg.Add(1)
		go func(id, user string) {
			defer wg.Done()
			<-start
			errs <- svc.Checkout(ctx, actor(user, domain.RoleBorrower), id, application.RequestMeta{})
		}(x.id, x.user)
	}
	close(start)
	wg.Wait()
	close(errs)
	ok := 0
	for err := range errs {
		if err == nil {
			ok++
		}
	}
	if ok != 1 {
		t.Fatalf("successful checkouts=%d, want 1", ok)
	}
	got, _ := repo.CaseByID(ctx, "case")
	if got.Status != domain.CaseBorrowed {
		t.Fatalf("status=%s", got.Status)
	}
}
