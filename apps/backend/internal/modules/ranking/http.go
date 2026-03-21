package ranking

import (
	stdhttp "net/http"

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
	group.GET("/power", h.power)
	group.GET("/arena-daily", h.arenaDaily)
}

func (h *Handler) power(c *gin.Context) {
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.Snapshot("power"), middleware.GetTraceID(c)))
}

func (h *Handler) arenaDaily(c *gin.Context) {
	c.JSON(stdhttp.StatusOK, httpx.Success(h.service.Snapshot("arena-daily"), middleware.GetTraceID(c)))
}
