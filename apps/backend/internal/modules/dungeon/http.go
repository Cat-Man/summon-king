package dungeon

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
	group.GET("/overview", h.worldMap)
	group.POST("/teleport", h.teleport)
	group.POST("/enter", h.enterDungeon)
	group.POST("/roll", h.rollDice)
	group.GET("/status", h.status)
	group.POST("/cultivation/start", h.startCultivation)
	group.POST("/cultivation/claim", h.claimCultivation)
}

func (h *Handler) worldMap(c *gin.Context) {
	world, err := h.service.GetWorldMap(c.Request.Context())
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to load world map", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(world, middleware.GetTraceID(c)))
}

func (h *Handler) teleport(c *gin.Context) {
	cityIDStr := c.PostForm("city_id")
	cityID, err := strconv.ParseInt(cityIDStr, 10, 64)
	if err != nil || cityID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "city_id is required", middleware.GetTraceID(c)))
		return
	}
	city, err := h.service.Teleport(c.Request.Context(), cityID)
	if err != nil {
		c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "city not found", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(city, middleware.GetTraceID(c)))
}

func (h *Handler) enterDungeon(c *gin.Context) {
	req, ok := h.parseDungeonRequest(c)
	if !ok {
		return
	}
	run, err := h.service.EnterDungeon(c.Request.Context(), req.playerID, req.dungeonID)
	if err != nil {
		if err == ErrDungeonLocked {
			c.JSON(stdhttp.StatusForbidden, httpx.Error(4031, "dungeon locked", middleware.GetTraceID(c)))
			return
		}
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to enter dungeon", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(run, middleware.GetTraceID(c)))
}

func (h *Handler) rollDice(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	run, err := h.service.RollDice(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5003, "failed to roll dice", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(run, middleware.GetTraceID(c)))
}

func (h *Handler) status(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	run, err := h.service.GetRun(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "dungeon run not found", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(run, middleware.GetTraceID(c)))
}

func (h *Handler) startCultivation(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	status, err := h.service.StartCultivation(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5004, "failed to start cultivation", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(status, middleware.GetTraceID(c)))
}

func (h *Handler) claimCultivation(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	status, err := h.service.ClaimCultivation(c.Request.Context(), playerID)
	if err != nil {
		if err == ErrCultivationNotFound {
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "cultivation not found", middleware.GetTraceID(c)))
			return
		}
		if err == ErrCultivationNotReady {
			c.JSON(stdhttp.StatusConflict, httpx.Error(4091, "cultivation not ready", middleware.GetTraceID(c)))
			return
		}
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5005, "failed to claim cultivation", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(status, middleware.GetTraceID(c)))
}

type dungeonRequest struct {
	playerID  int64
	dungeonID int64
}

func (h *Handler) parseDungeonRequest(c *gin.Context) (dungeonRequest, bool) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return dungeonRequest{}, false
	}
	dungeonID, err := strconv.ParseInt(c.PostForm("dungeon_id"), 10, 64)
	if err != nil || dungeonID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "dungeon_id is required", middleware.GetTraceID(c)))
		return dungeonRequest{}, false
	}
	return dungeonRequest{playerID: playerID, dungeonID: dungeonID}, true
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerIDStr := c.Query("player_id")
	if playerIDStr == "" {
		playerIDStr = c.PostForm("player_id")
	}
	if playerIDStr == "" {
		playerIDStr = c.GetHeader("X-Player-ID")
	}
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
