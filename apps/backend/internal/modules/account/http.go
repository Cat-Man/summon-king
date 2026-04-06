package account

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type guestLoginRequest struct {
	Nickname string `json:"nickname"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/guest-login", h.guestLogin)
}

func (h *Handler) guestLogin(c *gin.Context) {
	var req guestLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request", middleware.GetTraceID(c)))
		return
	}

	resp, err := h.service.GuestLogin(c.Request.Context(), req.Nickname)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, err.Error(), middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(resp, middleware.GetTraceID(c)))
}
