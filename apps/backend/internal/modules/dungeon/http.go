package dungeon

import (
	stdhttp "net/http"
	"strconv"
	"strings"

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
	maps := group.Group("/maps")
	maps.GET("/world", h.world)
	maps.POST("/teleport", h.teleport)

	dungeons := group.Group("/dungeons")
	dungeons.POST("/enter", h.enter)
	dungeons.POST("/roll", h.roll)
	dungeons.POST("/boss/claim", h.claimBoss)

	cultivation := group.Group("/cultivation")
	cultivation.POST("/start", h.startCultivation)
	cultivation.POST("/claim", h.claimCultivation)
}

func (h *Handler) world(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	world, err := h.service.GetWorldMap(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to get world map", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(world, middleware.GetTraceID(c)))
}

func (h *Handler) teleport(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		CityID   string `json:"city_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	city, err := h.service.Teleport(c.Request.Context(), req.PlayerID, req.CityID)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "failed to teleport", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(city, middleware.GetTraceID(c)))
}

func (h *Handler) enter(c *gin.Context) {
	var req struct {
		PlayerID  int64 `json:"player_id"`
		DungeonID int64 `json:"dungeon_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	run, err := h.service.EnterDungeon(c.Request.Context(), req.PlayerID, req.DungeonID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to enter dungeon", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(run, middleware.GetTraceID(c)))
}

func (h *Handler) roll(c *gin.Context) {
	var req struct {
		RunID string `json:"run_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	result, err := h.service.RollDice(c.Request.Context(), req.RunID)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "failed to roll dice", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func (h *Handler) claimBoss(c *gin.Context) {
	var req struct {
		RunID string `json:"run_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	run, err := h.service.ClaimBossReward(c.Request.Context(), req.RunID)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4004, "failed to claim boss reward", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(run, middleware.GetTraceID(c)))
}

func (h *Handler) startCultivation(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		MapID    string `json:"map_id"`
		Hours    int    `json:"hours"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	record, err := h.service.StartCultivation(c.Request.Context(), req.PlayerID, req.MapID, req.Hours)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5003, "failed to start cultivation", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(record, middleware.GetTraceID(c)))
}

func (h *Handler) claimCultivation(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		RecordID string `json:"record_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	record, err := h.service.ClaimCultivation(c.Request.Context(), req.PlayerID, req.RecordID)
	if err != nil {
		msg := "failed to claim cultivation"
		if err == ErrCultivationUnfinished {
			msg = "cultivation not finished"
		}
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4005, msg, middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(record, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(strings.TrimSpace(c.Query("player_id")), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
