package battle

import (
	stdhttp "net/http"
	"strings"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/detail", h.detail)
	group.GET("/replay", h.replay)
}

func (h *Handler) detail(c *gin.Context) {
	battleNo := strings.TrimSpace(c.Query("battle_no"))
	if battleNo == "" {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "battle_no is required", middleware.GetTraceID(c)))
		return
	}

	result, err := h.store.GetDetail(c.Request.Context(), battleNo)
	if err != nil {
		c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "battle not found", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(result, middleware.GetTraceID(c)))
}

func (h *Handler) replay(c *gin.Context) {
	battleNo := strings.TrimSpace(c.Query("battle_no"))
	if battleNo == "" {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "battle_no is required", middleware.GetTraceID(c)))
		return
	}

	replay, err := h.store.GetReplay(c.Request.Context(), battleNo)
	if err != nil {
		c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "battle not found", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(replay, middleware.GetTraceID(c)))
}
