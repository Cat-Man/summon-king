package arena

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
	Won      bool  `json:"won"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/status", h.status)
	group.POST("/challenge", h.challenge)
}

func (h *Handler) status(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	record, err := h.service.GetDailyRecord(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "record not found", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(record, middleware.GetTraceID(c)))
}

func (h *Handler) challenge(c *gin.Context) {
	var req challengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return
	}

	if err := h.service.RecordBattleResult(c.Request.Context(), req.PlayerID, req.Won); err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to record result", middleware.GetTraceID(c)))
		return
	}

	record, _ := h.service.GetDailyRecord(c.Request.Context(), req.PlayerID)
	c.JSON(stdhttp.StatusOK, httpx.Success(record, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
