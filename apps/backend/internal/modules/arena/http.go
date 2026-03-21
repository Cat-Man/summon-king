package arena

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/opponents", h.opponents)
	group.POST("/refresh", h.refresh)
	group.POST("/battle", h.battle)
	group.POST("/reward/claim", h.claimReward)
	group.GET("/daily-record", h.dailyRecord)
}

func (h *Handler) opponents(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.GetOpponents(c.Request.Context(), playerID), middleware.GetTraceID(c)))
}

func (h *Handler) refresh(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badArenaReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.RefreshOpponents(c.Request.Context(), req.PlayerID), middleware.GetTraceID(c)))
}

func (h *Handler) battle(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
		Won      bool  `json:"won"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badArenaReq(c)
		return
	}
	_ = h.service.RecordBattleResult(c.Request.Context(), req.PlayerID, req.Won)
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.GetDailyRecord(c.Request.Context(), req.PlayerID), middleware.GetTraceID(c)))
}

func (h *Handler) claimReward(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badArenaReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ClaimReward(c.Request.Context(), req.PlayerID), middleware.GetTraceID(c)))
}

func (h *Handler) dailyRecord(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.GetDailyRecord(c.Request.Context(), playerID), middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	var req struct {
		PlayerID int64 `form:"player_id"`
	}
	if err := c.ShouldBindQuery(&req); err != nil || req.PlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return req.PlayerID, true
}

func badArenaReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
}
