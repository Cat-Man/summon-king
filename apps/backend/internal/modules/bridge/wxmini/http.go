package wxmini

import (
	"errors"
	"io"
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type loginExchangeRequest struct {
	Code string `json:"code"`
}

type sessionBootstrapRequest struct {
	UnifiedToken string `json:"unified_token"`
	GameBaseURL  string `json:"game_base_url"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/login/exchange", h.loginExchange)
	group.POST("/session/bootstrap", h.sessionBootstrap)
}

func (h *Handler) loginExchange(c *gin.Context) {
	var req loginExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}

	resp, err := h.service.ExchangeLogin(c.Request.Context(), req.Code)
	if err != nil {
		switch {
		case errors.Is(err, ErrCodeRequired):
			c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "code is required", middleware.GetTraceID(c)))
		default:
			c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to exchange wxmini login", middleware.GetTraceID(c)))
		}
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(resp, middleware.GetTraceID(c)))
}

func (h *Handler) sessionBootstrap(c *gin.Context) {
	var req sessionBootstrapRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}

	resp, err := h.service.BootstrapSession(c.Request.Context(), req.UnifiedToken, req.GameBaseURL)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnifiedTokenRequired):
			c.JSON(stdhttp.StatusBadRequest, httpx.Error(4004, "unified_token is required", middleware.GetTraceID(c)))
		case errors.Is(err, ErrSessionNotFound):
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "session not found", middleware.GetTraceID(c)))
		default:
			c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to bootstrap wxmini session", middleware.GetTraceID(c)))
		}
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(resp, middleware.GetTraceID(c)))
}
