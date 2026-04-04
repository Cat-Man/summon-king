package growth

import (
	"io"
	stdhttp "net/http"
	"strconv"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	spirit   *SpiritService
	bone     *BoneService
	soul     *SoulService
	cauldron *CauldronService
	manor    *ManorService
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		spirit:   NewSpiritService(repo),
		bone:     NewBoneService(repo),
		soul:     NewSoulService(repo),
		cauldron: NewCauldronService(repo),
		manor:    NewManorService(repo),
	}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/wallet", h.wallet)
	group.POST("/spirit/wash", h.washSpirit)
	group.GET("/bone", h.boneState)
	group.POST("/bone/upgrade", h.boneUpgrade)
	group.GET("/soul", h.soulState)
	group.POST("/soul/upgrade", h.soulUpgrade)
	group.POST("/cauldron/stir", h.cauldronStir)
	group.GET("/manor", h.manorState)
	group.POST("/manor/harvest", h.manorHarvest)
	group.POST("/manor/plant", h.manorPlant)
}

func (h *Handler) wallet(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data := h.spirit.GetWallet(c.Request.Context(), playerID)
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func (h *Handler) washSpirit(c *gin.Context) {
	var req WashOption
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request", middleware.GetTraceID(c)))
		return
	}
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	err := h.spirit.WashSpirit(c.Request.Context(), playerID, req)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "wash failed", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"washed": true}, middleware.GetTraceID(c)))
}

func (h *Handler) boneState(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data := h.bone.GetBoneState(c.Request.Context(), playerID)
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func (h *Handler) boneUpgrade(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data, err := h.bone.Upgrade(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5003, "bone upgrade failed", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func (h *Handler) soulState(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data := h.soul.GetSoulState(c.Request.Context(), playerID)
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func (h *Handler) soulUpgrade(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data, err := h.soul.Upgrade(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5005, "soul upgrade failed", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func (h *Handler) cauldronStir(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	msg, err := h.cauldron.Stir(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "cauldron failed", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"message": msg}, middleware.GetTraceID(c)))
}

func (h *Handler) manorState(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data := h.manor.ListPlots(c.Request.Context(), playerID)
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func (h *Handler) manorHarvest(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data, err := h.manor.Harvest(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5004, "manor harvest failed", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func (h *Handler) manorPlant(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}
	data, err := h.manor.Plant(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5006, "manor plant failed", middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
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
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}

	return playerID, true
}
