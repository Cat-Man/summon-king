package commerce

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/vip", h.vip)
	group.POST("/vip/chest/claim", h.claimChest)
	group.POST("/orders/create", h.createOrder)
	group.POST("/orders/pay-callback", h.payCallback)
	group.POST("/orders/deliver", h.deliver)
	group.POST("/signin/claim", h.signin)
	group.POST("/redeem", h.redeem)
}

func (h *Handler) vip(c *gin.Context) {
	var req struct {
		PlayerID int64 `form:"player_id"`
	}
	if err := c.ShouldBindQuery(&req); err != nil || req.PlayerID == 0 {
		badCommerceReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.GetVIPState(c.Request.Context(), req.PlayerID), middleware.GetTraceID(c)))
}

func (h *Handler) claimChest(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badCommerceReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ClaimVIPDailyChest(c.Request.Context(), req.PlayerID), middleware.GetTraceID(c)))
}

func (h *Handler) createOrder(c *gin.Context) {
	var req struct {
		PlayerID       int64  `json:"player_id"`
		IdempotencyKey string `json:"idempotency_key"`
		ProductID      int64  `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badCommerceReq(c)
		return
	}
	order, err := h.service.CreateOrder(c.Request.Context(), req.PlayerID, req.IdempotencyKey, req.ProductID)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, err.Error(), middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(order, middleware.GetTraceID(c)))
}

func (h *Handler) payCallback(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badCommerceReq(c)
		return
	}
	order, err := h.service.MarkOrderPaid(c.Request.Context(), req.OrderNo)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, err.Error(), middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(order, middleware.GetTraceID(c)))
}

func (h *Handler) deliver(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badCommerceReq(c)
		return
	}
	order, err := h.service.DeliverOrder(c.Request.Context(), req.OrderNo)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4004, err.Error(), middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(order, middleware.GetTraceID(c)))
}

func (h *Handler) signin(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badCommerceReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ClaimSignin(c.Request.Context(), req.PlayerID), middleware.GetTraceID(c)))
}

func (h *Handler) redeem(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badCommerceReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.RedeemCode(c.Request.Context(), req.PlayerID, req.Code), middleware.GetTraceID(c)))
}

func badCommerceReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
}
