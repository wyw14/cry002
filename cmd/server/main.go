package main

import (
	"context"
	"errors"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/config"
	"github.com/wyw14/cry002/internal/platform/demo"
	"github.com/wyw14/cry002/internal/platform/storage"
	"github.com/wyw14/cry002/internal/repository/memory"
	pg "github.com/wyw14/cry002/internal/repository/postgres"
	"github.com/wyw14/cry002/internal/service"
	transport "github.com/wyw14/cry002/internal/transport/http"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log, err := zap.NewProduction()
	if cfg.Environment == "development" {
		log, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	root, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var repo application.Repository
	ready := func() error { return nil }
	if cfg.AllowInMemory {
		repo = memory.New()
	} else {
		ctx, stop := context.WithTimeout(root, 10*time.Second)
		db, e := pg.Open(ctx, cfg.DatabaseURL)
		stop()
		if e != nil {
			log.Fatal("connect database", zap.Error(e))
		}
		if e = db.Migrate(root); e != nil {
			log.Fatal("migrate", zap.Error(e))
		}
		repo = pg.NewRepository(db.DB)
		ready = func() error {
			ctx, c := context.WithTimeout(context.Background(), time.Second)
			defer c()
			return db.Ready(ctx)
		}
	}
	clock := service.RealClock{}
	ids := service.UUIDGenerator{}
	passwords := service.BcryptHasher{}
	tokens := service.JWTManager{Secret: []byte(cfg.JWTSecret), AccessTTL: cfg.AccessTTL}
	if err := demo.Seed(root, repo, passwords, clock.Now()); err != nil {
		log.Fatal("seed", zap.Error(err))
	}
	files, err := storage.NewLocal("./data/attachments")
	if err != nil {
		log.Fatal("storage", zap.Error(err))
	}
	h := &transport.Handler{Repo: repo, Auth: application.NewAuthService(repo, clock, ids, passwords, tokens, cfg.RefreshTTL), Cases: application.NewCaseService(repo, clock, ids), Borrows: application.NewBorrowService(repo, clock, ids), Classes: application.NewClassificationService(repo, clock, ids), Attachments: application.NewAttachmentService(repo, files, clock, ids), Admin: application.NewAdminService(repo), Ready: ready}
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: transport.Router(h, tokens, log, cfg.RequestTimeout, cfg.CORSOrigins), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		log.Info("server listening", zap.String("address", cfg.HTTPAddr))
		if e := srv.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			log.Fatal("http server", zap.Error(e))
		}
	}()
	<-root.Done()
	ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown", zap.Error(err))
		os.Exit(1)
	}
}
