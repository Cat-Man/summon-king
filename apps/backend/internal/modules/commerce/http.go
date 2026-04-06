package commerce

import (
	"errors"
	"io"
	stdhttp "net/http"
	"strconv"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type claimRequest struct {
	PlayerID int64 `json:"player_id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterSigninRoutes(group *gin.RouterGroup) {
	group.GET("/index", h.signinIndex)
	group.POST("/claim", h.signinClaim)
}

func (h *Handler) RegisterVIPRoutes(group *gin.RouterGroup) {
	group.GET("/index", h.vipIndex)
	group.POST("/claim-daily", h.vipClaimDaily)
}

func (h *Handler) signinIndex(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	index, err := h.service.GetSigninIndex(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to load signin index", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(index, middleware.GetTraceID(c)))
}

func (h *Handler) signinClaim(c *gin.Context) {
	var req claimRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return
	}

	result, err := h.service.ClaimSignin(c.Request.Context(), req.PlayerID)
	if err != nil {
		if errors.Is(err, ErrAlreadyClaimedToday) {
			c.JSON(stdhttp.StatusConflict, httpx.Error(4091, "signin already claimed today", middleware.GetTraceID(c)))
			return
		}
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to claim signin", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func (h *Handler) vipIndex(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	index, err := h.service.GetVIPIndex(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5003, "failed to load vip index", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(index, middleware.GetTraceID(c)))
}

func (h *Handler) vipClaimDaily(c *gin.Context) {
	var req claimRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return
	}

	result, err := h.service.ClaimVIPDaily(c.Request.Context(), req.PlayerID)
	if err != nil {
		if errors.Is(err, ErrAlreadyClaimedToday) {
			c.JSON(stdhttp.StatusConflict, httpx.Error(4092, "vip daily already claimed today", middleware.GetTraceID(c)))
			return
		}
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5004, "failed to claim vip daily", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
