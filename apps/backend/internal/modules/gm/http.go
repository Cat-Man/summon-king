package gm

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/grant-currency", h.grantCurrency)
	group.POST("/grant-item", h.grantItem)
	group.POST("/grant-pet", h.grantPet)
	group.POST("/send-mail", h.sendMail)
	group.POST("/reset-daily", h.resetDaily)
	group.POST("/unblock-player", h.unblock)
	group.POST("/mute-player", h.mute)
	group.POST("/ban-player", h.ban)
	group.POST("/fix-alliance", h.fixAlliance)
	group.GET("/audit-logs", h.auditLogs)
}

func (h *Handler) grantCurrency(c *gin.Context) {
	var req GrantCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.GrantCurrency(c.Request.Context(), req))
}

func (h *Handler) grantItem(c *gin.Context) {
	var req GrantItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.GrantItem(c.Request.Context(), req))
}

func (h *Handler) grantPet(c *gin.Context) {
	var req GrantPetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.GrantPet(c.Request.Context(), req))
}

func (h *Handler) sendMail(c *gin.Context) {
	var req MailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.SendMail(c.Request.Context(), req))
}

func (h *Handler) resetDaily(c *gin.Context) {
	var req ResetDailyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.ResetDaily(c.Request.Context(), req))
}

func (h *Handler) unblock(c *gin.Context) {
	var req PlayerFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.UnblockPlayer(c.Request.Context(), req))
}

func (h *Handler) mute(c *gin.Context) {
	var req PlayerFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.MutePlayer(c.Request.Context(), req))
}

func (h *Handler) ban(c *gin.Context) {
	var req PlayerFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.BanPlayer(c.Request.Context(), req))
}

func (h *Handler) fixAlliance(c *gin.Context) {
	var req FixAllianceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badGMReq(c)
		return
	}
	h.exec(c, h.service.FixAlliance(c.Request.Context(), req))
}

func (h *Handler) auditLogs(c *gin.Context) {
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ListAuditLogs(c.Request.Context()), middleware.GetTraceID(c)))
}

func (h *Handler) exec(c *gin.Context, err error) {
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4202, err.Error(), middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"ok": true}, middleware.GetTraceID(c)))
}

func badGMReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4201, "invalid request body", middleware.GetTraceID(c)))
}
