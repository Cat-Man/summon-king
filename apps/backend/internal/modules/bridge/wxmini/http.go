package wxmini

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/login", h.login)
	group.POST("/pay", h.pay)
}

func (h *Handler) login(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ExchangeLogin(req.Code), middleware.GetTraceID(c)))
}

func (h *Handler) pay(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.BuildPayParams(req.OrderNo), middleware.GetTraceID(c)))
}

func badReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
}
