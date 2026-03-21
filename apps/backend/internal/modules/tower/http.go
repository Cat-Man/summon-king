package tower

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/pagoda/challenge", h.pagoda)
	group.POST("/spirit/challenge", h.spirit)
}

func (h *Handler) pagoda(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
		Victory  bool  `json:"victory"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ChallengePagoda(c.Request.Context(), req.PlayerID, req.Victory), middleware.GetTraceID(c)))
}

func (h *Handler) spirit(c *gin.Context) {
	var req struct {
		PlayerID int64 `json:"player_id"`
		Victory  bool  `json:"victory"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.ChallengeSpiritTower(c.Request.Context(), req.PlayerID, req.Victory), middleware.GetTraceID(c)))
}

func badReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
}
