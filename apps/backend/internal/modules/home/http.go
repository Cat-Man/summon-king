package home

import (
	stdhttp "net/http"
	"strconv"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/overview", h.overview)
}

func (h *Handler) overview(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	overview, err := h.service.GetOverview(c.Request.Context(), playerID, middleware.GetAuthToken(c))
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to load overview", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(overview, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
