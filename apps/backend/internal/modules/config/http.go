package config

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/modules", h.modules)
	group.GET("/detail", h.detail)
	group.POST("/save-draft", h.saveDraft)
	group.POST("/validate", h.validate)
	group.POST("/publish", h.publish)
	group.POST("/rollback", h.rollback)
	group.GET("/publish-logs", h.publishLogs)
}

func (h *Handler) modules(c *gin.Context) {
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ListModules(c.Request.Context()), middleware.GetTraceID(c)))
}

func (h *Handler) detail(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		badConfigReq(c)
		return
	}
	detail, err := h.service.GetDetail(c.Request.Context(), module)
	if err != nil {
		configErr(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(detail, middleware.GetTraceID(c)))
}

func (h *Handler) saveDraft(c *gin.Context) {
	var req struct {
		Module  string `json:"module"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badConfigReq(c)
		return
	}
	detail, err := h.service.SaveDraft(c.Request.Context(), req.Module, req.Content)
	if err != nil {
		configErr(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(detail, middleware.GetTraceID(c)))
}

func (h *Handler) validate(c *gin.Context) {
	var req struct {
		Module string `json:"module"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badConfigReq(c)
		return
	}
	result, err := h.service.Validate(c.Request.Context(), req.Module)
	if err != nil {
		configErr(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func (h *Handler) publish(c *gin.Context) {
	var req struct {
		Module     string `json:"module"`
		OperatorID int64  `json:"operator_id"`
		Summary    string `json:"summary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badConfigReq(c)
		return
	}
	log, err := h.service.Publish(c.Request.Context(), req.Module, req.OperatorID, req.Summary)
	if err != nil {
		configErr(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(log, middleware.GetTraceID(c)))
}

func (h *Handler) rollback(c *gin.Context) {
	var req struct {
		Module     string `json:"module"`
		Version    int    `json:"version"`
		OperatorID int64  `json:"operator_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badConfigReq(c)
		return
	}
	log, err := h.service.Rollback(c.Request.Context(), req.Module, req.Version, req.OperatorID)
	if err != nil {
		configErr(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(log, middleware.GetTraceID(c)))
}

func (h *Handler) publishLogs(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		badConfigReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.PublishLogs(c.Request.Context(), module), middleware.GetTraceID(c)))
}

func badConfigReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4101, "invalid request body", middleware.GetTraceID(c)))
}

func configErr(c *gin.Context, err error) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4102, err.Error(), middleware.GetTraceID(c)))
}
