package httptransport

import (
	"github.com/gin-gonic/gin"
	"github.com/wyw14/cry002/internal/application"
	appmw "github.com/wyw14/cry002/internal/middleware"
	"go.uber.org/zap"
	"time"
)

func Router(h *Handler, tokens application.TokenManager, log *zap.Logger, timeout time.Duration, origins []string) *gin.Engine {
	r := gin.New()
	r.Use(appmw.RequestID(), appmw.Recover(log), appmw.Logger(log), appmw.SecurityHeaders(), appmw.CORS(origins), appmw.Timeout(timeout))
	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Readiness)
	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", h.Register)
	v1.POST("/auth/login", h.Login)
	v1.POST("/auth/refresh", h.Refresh)
	v1.POST("/auth/logout", h.Logout)
	s := v1.Group("")
	s.Use(appmw.Auth(tokens))
	s.GET("/cases", h.ListCases)
	s.POST("/cases", h.CreateCase)
	s.GET("/cases/export", h.Export)
	s.GET("/cases/:caseID", h.GetCase)
	s.POST("/cases/:caseID/submit", h.SubmitCase)
	s.POST("/cases/:caseID/review", h.ReviewCase)
	s.PATCH("/cases/:caseID/archived-metadata", h.UpdateArchived)
	s.POST("/cases/:caseID/attachments", h.Upload)
	s.GET("/attachments/:attachmentID", h.Download)
	s.GET("/classifications", h.ListClasses)
	s.PATCH("/classifications/:classID/move", h.MoveClass)
	s.GET("/borrows", h.ListBorrows)
	s.POST("/borrows", h.ApplyBorrow)
	s.POST("/borrows/:borrowID/review", h.ReviewBorrow)
	s.POST("/borrows/:borrowID/checkout", h.Checkout)
	s.POST("/borrows/:borrowID/return", h.Return)
	s.GET("/audits", h.Audits)
	s.GET("/admin/stats", h.Stats)
	return r
}
