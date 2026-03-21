package account

import (
	"io"
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type loginRequest struct {
	Channel string `json:"channel"`
}

type loginResponse struct {
	PlayerID int64  `json:"player_id"`
	Token    string `json:"token"`
	Channel  string `json:"channel"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/login", h.login)
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
		return
	}

	player, err := h.service.CreateGuestPlayer(c.Request.Context(), req.Channel)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to create guest player", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(loginResponse{
		PlayerID: player.PlayerID,
		Token:    player.Token,
		Channel:  player.Channel,
	}, middleware.GetTraceID(c)))
}
