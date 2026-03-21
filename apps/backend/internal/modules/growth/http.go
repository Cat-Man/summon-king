package growth

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/wallet", h.wallet)
	group.POST("/bones/upgrade", h.upgradeBone)
	group.POST("/spirits/wash", h.washSpirit)
	group.POST("/souls/hunt", h.huntSoul)
	group.POST("/cauldron/refine", h.refineCauldron)
	group.POST("/ascension/start", h.startAscension)
	group.POST("/manor/plant", h.plantCrop)
	group.POST("/manor/harvest", h.harvestCrop)
}

func (h *Handler) wallet(c *gin.Context) {
	playerID, ok := parsePlayerIDBodyOrQuery(c)
	if !ok {
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.GetWallet(c.Request.Context(), playerID), middleware.GetTraceID(c)))
}

func (h *Handler) upgradeBone(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		BoneID   string `json:"bone_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, _ := h.service.UpgradeBone(c.Request.Context(), req.PlayerID, req.BoneID)
	c.JSON(stdhttp.StatusOK, httpx.Success(res, middleware.GetTraceID(c)))
}

func (h *Handler) washSpirit(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
		SpiritID int64 `json:"spirit_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	if err := h.service.WashSpirit(c.Request.Context(), req.PlayerID, req.SpiritID, WashOption{}); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "failed to wash spirit", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"washed": true}, middleware.GetTraceID(c)))
}

func (h *Handler) huntSoul(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		SoulID   string `json:"soul_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, _ := h.service.HuntSoul(c.Request.Context(), req.PlayerID, req.SoulID)
	c.JSON(stdhttp.StatusOK, httpx.Success(res, middleware.GetTraceID(c)))
}

func (h *Handler) refineCauldron(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, _ := h.service.RefineCauldron(c.Request.Context(), req.PlayerID)
	c.JSON(stdhttp.StatusOK, httpx.Success(res, middleware.GetTraceID(c)))
}

func (h *Handler) startAscension(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
		PetID    int64 `json:"pet_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, _ := h.service.StartAscension(c.Request.Context(), req.PlayerID, req.PetID)
	c.JSON(stdhttp.StatusOK, httpx.Success(res, middleware.GetTraceID(c)))
}

func (h *Handler) plantCrop(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		SeedID   string `json:"seed_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, _ := h.service.PlantCrop(c.Request.Context(), req.PlayerID, req.SeedID)
	c.JSON(stdhttp.StatusOK, httpx.Success(res, middleware.GetTraceID(c)))
}

func (h *Handler) harvestCrop(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.HarvestCrop(c.Request.Context(), req.PlayerID)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "crop not ready", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(res, middleware.GetTraceID(c)))
}

func badReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
}

func parsePlayerIDBodyOrQuery(c *gin.Context) (int64, bool) {
	if playerID := c.Query("player_id"); playerID != "" {
		var req struct {
			PlayerID int64 `form:"player_id"`
		}
		if err := c.ShouldBindQuery(&req); err == nil && req.PlayerID != 0 {
			return req.PlayerID, true
		}
	}
	var req struct {
		PlayerID int64 `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.PlayerID != 0 {
		return req.PlayerID, true
	}
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4004, "player_id is required", middleware.GetTraceID(c)))
	return 0, false
}
