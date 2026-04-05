package pet

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

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/team", h.team)
	group.POST("/team/main", h.setMainPet)
}

func (h *Handler) team(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	data, err := h.service.GetCollection(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to load pet collection", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

type setMainPetRequest struct {
	PlayerID int64 `json:"player_id"`
	PetID    int64 `json:"pet_id"`
}

func (h *Handler) setMainPet(c *gin.Context) {
	var req setMainPetRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 || req.PetID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "player_id and pet_id are required", middleware.GetTraceID(c)))
		return
	}

	data, err := h.service.SetMainPet(c.Request.Context(), req.PlayerID, req.PetID)
	if err != nil {
		if errors.Is(err, ErrPetNotFound) {
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "pet not found", middleware.GetTraceID(c)))
			return
		}
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to switch main pet", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
