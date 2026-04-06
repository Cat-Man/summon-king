package tower

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

type challengeRequest struct {
	PlayerID int64 `json:"player_id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/pagoda/status", h.pagodaStatus)
	group.GET("/spirit-tower/status", h.spiritStatus)
	group.POST("/pagoda/start", h.startPagoda)
	group.POST("/spirit-tower/start", h.startSpirit)
}

func (h *Handler) pagodaStatus(c *gin.Context) {
	h.status(c, "pagoda")
}

func (h *Handler) spiritStatus(c *gin.Context) {
	h.status(c, "spirit")
}

func (h *Handler) status(c *gin.Context, tower string) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.GetStatus(c.Request.Context(), playerID, tower), middleware.GetTraceID(c)))
}

func (h *Handler) startPagoda(c *gin.Context) {
	h.start(c, "pagoda")
}

func (h *Handler) startSpirit(c *gin.Context) {
	h.start(c, "spirit")
}

func (h *Handler) start(c *gin.Context, tower string) {
	var req challengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id required", middleware.GetTraceID(c)))
		return
	}

	result, err := h.service.StartChallenge(c.Request.Context(), req.PlayerID, tower)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to start challenge", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
