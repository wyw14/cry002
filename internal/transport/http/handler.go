package httptransport

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	appmw "github.com/wyw14/cry002/internal/middleware"
	"net/http"
	"strconv"
)

type Handler struct {
	Repo        application.Repository
	Auth        *application.AuthService
	Cases       *application.CaseService
	Borrows     *application.BorrowService
	Classes     *application.ClassificationService
	Attachments *application.AttachmentService
	Admin       *application.AdminService
	Ready       func() error
}

func (h *Handler) actor(c *gin.Context) (domain.User, error) {
	return h.Repo.UserByID(c.Request.Context(), appmw.ActorID(c))
}
func (h *Handler) Health(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func (h *Handler) Readiness(c *gin.Context) {
	if h.Ready != nil {
		if err := h.Ready(); err != nil {
			c.JSON(http.StatusServiceUnavailable, appmw.ErrorBody(c, "not_ready", "database unavailable", nil))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
func (h *Handler) Register(c *gin.Context) {
	var in struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		DisplayName string `json:"display_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		writeError(c, err)
		return
	}
	u, err := h.Auth.Register(c, in.Email, in.Password, in.DisplayName, domain.RoleBorrower)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, u)
}
func (h *Handler) Login(c *gin.Context) {
	var in struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		writeError(c, err)
		return
	}
	v, err := h.Auth.Login(c, in.Email, in.Password)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}
func (h *Handler) Refresh(c *gin.Context) {
	var in struct {
		Token string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		writeError(c, err)
		return
	}
	v, err := h.Auth.Refresh(c, in.Token)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}
func (h *Handler) Logout(c *gin.Context) {
	var in struct {
		Token string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		writeError(c, err)
		return
	}
	if err := h.Auth.Logout(c, in.Token); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) ListCases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	year, _ := strconv.Atoi(c.Query("year"))
	items, total, err := h.Cases.Search(c, application.CaseFilter{Query: c.Query("q"), Status: domain.CaseStatus(c.Query("status")), ClassificationID: c.Query("classification_id"), Department: c.Query("department_id"), Year: year, Page: page, PageSize: size, Sort: c.Query("sort")})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": page, "page_size": size})
}
func (h *Handler) GetCase(c *gin.Context) {
	x, err := h.Repo.CaseByID(c, c.Param("caseID"))
	if err != nil {
		writeError(c, err)
		return
	}
	ms, _ := h.Repo.Materials(c, x.ID)
	vs, _ := h.Repo.Versions(c, x.ID)
	c.JSON(http.StatusOK, gin.H{"case": x, "materials": ms, "versions": vs})
}
func (h *Handler) CreateCase(c *gin.Context) {
	actor, err := h.actor(c)
	if err != nil {
		writeError(c, err)
		return
	}
	var in struct {
		Case      domain.CaseFile   `json:"case"`
		Materials []domain.Material `json:"materials"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		writeError(c, err)
		return
	}
	out, err := h.Cases.Create(c, actor, in.Case, in.Materials, appmw.Meta(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}
func (h *Handler) SubmitCase(c *gin.Context) {
	a, err := h.actor(c)
	if err == nil {
		err = h.Cases.Submit(c, a, c.Param("caseID"), appmw.Meta(c))
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) ReviewCase(c *gin.Context) {
	a, err := h.actor(c)
	var in struct {
		Approved bool   `json:"approved"`
		Opinion  string `json:"opinion"`
	}
	if err == nil {
		err = c.ShouldBindJSON(&in)
	}
	if err == nil {
		err = h.Cases.Review(c, a, c.Param("caseID"), in.Approved, in.Opinion, appmw.Meta(c))
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) UpdateArchived(c *gin.Context) {
	a, err := h.actor(c)
	var in struct {
		Reason string          `json:"reason"`
		Case   domain.CaseFile `json:"case"`
	}
	if err == nil {
		err = c.ShouldBindJSON(&in)
	}
	if err == nil {
		err = h.Cases.UpdateMetadata(c, a, c.Param("caseID"), in.Reason, in.Case, appmw.Meta(c))
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) Export(c *gin.Context) {
	a, err := h.actor(c)
	if err != nil {
		writeError(c, err)
		return
	}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=cases.csv")
	err = h.Cases.ExportCSV(context.WithoutCancel(c.Request.Context()), a, application.CaseFilter{Query: c.Query("q"), Page: 1, PageSize: 100}, c.Writer, appmw.Meta(c))
	if err != nil {
		writeError(c, err)
	}
}
func (h *Handler) ListClasses(c *gin.Context) {
	items, err := h.Classes.Tree(c)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *Handler) MoveClass(c *gin.Context) {
	a, err := h.actor(c)
	var in struct {
		ParentID string `json:"parent_id"`
	}
	if err == nil {
		err = c.ShouldBindJSON(&in)
	}
	if err == nil {
		err = h.Classes.Move(c, a, c.Param("classID"), in.ParentID)
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
