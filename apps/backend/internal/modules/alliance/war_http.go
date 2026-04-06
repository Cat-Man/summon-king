package alliance

import (
	"errors"
	"io"
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type WarHandler struct {
	service *WarService
}

type registerWarTargetRequest struct {
	PlayerID int64  `json:"player_id"`
	Target   string `json:"target"`
}

func NewWarHandler(service *WarService) *WarHandler {
	return &WarHandler{service: service}
}

func (h *WarHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/index", h.index)
	group.POST("/register-target", h.registerTarget)
}

func (h *WarHandler) index(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	index, err := h.service.GetIndex(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5007, "failed to load alliance war index", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(index, middleware.GetTraceID(c)))
}

func (h *WarHandler) registerTarget(c *gin.Context) {
	var req registerWarTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 || req.Target == "" {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4006, "player_id and target are required", middleware.GetTraceID(c)))
		return
	}

	index, err := h.service.RegisterTarget(c.Request.Context(), req.PlayerID, req.Target)
	if err != nil {
		switch {
		case errors.Is(err, ErrAllianceNotFound):
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "alliance not found", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAlliancePermissionDenied):
			c.JSON(stdhttp.StatusForbidden, httpx.Error(4031, "alliance permission denied", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAllianceWarTargetRequired), errors.Is(err, ErrAllianceWarTargetInvalid):
			c.JSON(stdhttp.StatusBadRequest, httpx.Error(4007, "alliance war target is invalid", middleware.GetTraceID(c)))
		default:
			c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5008, "failed to register alliance war target", middleware.GetTraceID(c)))
		}
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(index, middleware.GetTraceID(c)))
}
