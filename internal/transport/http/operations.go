package httptransport

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
	appmw "github.com/wyw14/cry002/internal/middleware"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) ApplyBorrow(c *gin.Context) {
	a, err := h.actor(c)
	var in struct {
		CaseID  string    `json:"case_id" binding:"required"`
		Purpose string    `json:"purpose" binding:"required"`
		DueAt   time.Time `json:"due_at" binding:"required"`
	}
	if err == nil {
		err = c.ShouldBindJSON(&in)
	}
	var out domain.BorrowRequest
	if err == nil {
		out, err = h.Borrows.Apply(c, a, in.CaseID, in.Purpose, in.DueAt, appmw.Meta(c))
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}
func (h *Handler) ListBorrows(c *gin.Context) {
	items, err := h.Repo.ListBorrows(c, c.Query("status"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *Handler) ReviewBorrow(c *gin.Context) {
	a, err := h.actor(c)
	var in struct {
		Approved bool `json:"approved"`
	}
	if err == nil {
		err = c.ShouldBindJSON(&in)
	}
	if err == nil {
		err = h.Borrows.Approve(c, a, c.Param("borrowID"), in.Approved, appmw.Meta(c))
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) Checkout(c *gin.Context) {
	a, err := h.actor(c)
	if err == nil {
		err = h.Borrows.Checkout(c, a, c.Param("borrowID"), appmw.Meta(c))
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) Return(c *gin.Context) {
	a, err := h.actor(c)
	if err == nil {
		err = h.Borrows.Return(c, a, c.Param("borrowID"), appmw.Meta(c))
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) Upload(c *gin.Context) {
	a, err := h.actor(c)
	if err != nil {
		writeError(c, err)
		return
	}
	f, err := c.FormFile("file")
	if err != nil {
		writeError(c, err)
		return
	}
	src, err := f.Open()
	if err != nil {
		writeError(c, err)
		return
	}
	defer src.Close()
	out, err := h.Attachments.Upload(c, a, c.Param("caseID"), c.PostForm("material_id"), f.Filename, f.Header.Get("Content-Type"), f.Size, src, appmw.Meta(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}
func (h *Handler) Download(c *gin.Context) {
	a, err := h.actor(c)
	if err != nil {
		writeError(c, err)
		return
	}
	meta, r, err := h.Attachments.Open(c, a, c.Param("attachmentID"))
	if err != nil {
		writeError(c, err)
		return
	}
	defer r.Close()
	c.Header("Content-Disposition", "attachment; filename="+strconv.Quote(meta.OriginalName))
	c.DataFromReader(http.StatusOK, meta.Size, meta.ContentType, r, nil)
}
func (h *Handler) Audits(c *gin.Context) {
	items, err := h.Repo.ListAudits(c, c.Query("resource"), c.Query("resource_id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *Handler) Stats(c *gin.Context) {
	a, err := h.actor(c)
	var out application.AdminStats
	if err == nil {
		out, err = h.Admin.Stats(c, a)
	}
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func writeError(c *gin.Context, err error) {
	status, code, msg := http.StatusBadRequest, "invalid_request", "request validation failed"
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status, code, msg = http.StatusNotFound, "not_found", "resource not found"
	case errors.Is(err, domain.ErrUnauthorized):
		status, code, msg = http.StatusUnauthorized, "unauthorized", "authentication failed"
	case errors.Is(err, domain.ErrForbidden):
		status, code, msg = http.StatusForbidden, "forbidden", "operation is not allowed"
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrAlreadyBorrowed):
		status, code, msg = http.StatusConflict, "conflict", "resource state conflicts with request"
	case errors.Is(err, domain.ErrInvalidState):
		status, code, msg = http.StatusUnprocessableEntity, "invalid_state", "state transition is not allowed"
	}
	c.JSON(status, appmw.ErrorBody(c, code, msg, nil))
}
