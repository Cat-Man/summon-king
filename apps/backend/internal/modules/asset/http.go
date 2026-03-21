package asset

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
	group.GET("/wallet", h.wallet)
	group.GET("/inventory", h.inventory)
	group.POST("/inventory/use", h.useInventory)
	group.POST("/inventory/sell", h.sellInventory)
}

func (h *Handler) wallet(c *gin.Context) {
	playerID, ok := parsePlayerIDQuery(c)
	if !ok {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return
	}

	wallet, err := h.service.GetWallet(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to get wallet", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(wallet, middleware.GetTraceID(c)))
}

func (h *Handler) inventory(c *gin.Context) {
	playerID, ok := parsePlayerIDQuery(c)
	if !ok {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return
	}

	inventory, err := h.service.GetInventory(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to get inventory", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(inventory, middleware.GetTraceID(c)))
}

func (h *Handler) useInventory(c *gin.Context) {
	var req InventoryOperateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request body", middleware.GetTraceID(c)))
		return
	}

	result, err := h.service.UseInventory(c.Request.Context(), req)
	if err != nil {
		h.handleInventoryError(c, err)
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func (h *Handler) sellInventory(c *gin.Context) {
	var req InventoryOperateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request body", middleware.GetTraceID(c)))
		return
	}

	result, err := h.service.SellInventory(c.Request.Context(), req)
	if err != nil {
		h.handleInventoryError(c, err)
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func (h *Handler) handleInventoryError(c *gin.Context, err error) {
	switch err {
	case ErrInvalidRequest:
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "invalid asset request", middleware.GetTraceID(c)))
	case ErrInventoryNotEnough:
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4004, "inventory not enough", middleware.GetTraceID(c)))
	case ErrInventoryItemAbsent:
		c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "inventory item not found", middleware.GetTraceID(c)))
	default:
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5003, "failed to operate inventory", middleware.GetTraceID(c)))
	}
}

func parsePlayerIDQuery(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		return 0, false
	}
	return playerID, true
}
