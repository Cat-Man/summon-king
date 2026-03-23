package pet

import (
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

type saveTeamRequest struct {
	PlayerID int64   `json:"player_id"`
	PetIDs   []int64 `json:"pet_ids"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/catalog", h.catalog)
	group.GET("/list", h.list)
	group.GET("/detail", h.detail)
	group.GET("/team", h.team)
	group.POST("/team/save", h.saveTeam)
}

func (h *Handler) catalog(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	items, err := h.service.GetCatalog(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to fetch pet catalog", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(items, middleware.GetTraceID(c)))
}

func (h *Handler) list(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	items, err := h.service.GetPlayerPets(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5003, "failed to fetch player pets", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(items, middleware.GetTraceID(c)))
}

func (h *Handler) detail(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	petID, err := strconv.ParseInt(c.Query("pet_id"), 10, 64)
	if err != nil || petID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4004, "pet_id is required", middleware.GetTraceID(c)))
		return
	}

	item, err := h.service.GetDetail(c.Request.Context(), playerID, petID)
	if err == ErrPetNotFound {
		c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "pet not found", middleware.GetTraceID(c)))
		return
	}
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5004, "failed to fetch pet detail", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(item, middleware.GetTraceID(c)))
}

func (h *Handler) saveTeam(c *gin.Context) {
	var req saveTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return
	}

	if err := h.service.SaveTeam(c.Request.Context(), req.PlayerID, req.PetIDs); err == ErrDuplicatePetInTeam {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4005, "duplicate pet in team", middleware.GetTraceID(c)))
		return
	} else if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5005, "failed to save pet team", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"saved": true}, middleware.GetTraceID(c)))
}

func (h *Handler) team(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	team, err := h.service.GetTeam(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5006, "failed to fetch pet team", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(team, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
