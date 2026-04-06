package alliance

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

type createAllianceRequest struct {
	PlayerID int64  `json:"player_id"`
	Name     string `json:"name"`
}

type applyAllianceRequest struct {
	PlayerID   int64 `json:"player_id"`
	AllianceID int64 `json:"alliance_id"`
}

type approveAllianceRequest struct {
	PlayerID          int64 `json:"player_id"`
	ApplicantPlayerID int64 `json:"applicant_player_id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/index", h.index)
	group.GET("/hall", h.hall)
	group.GET("/apply-list", h.applyList)
	group.POST("/create", h.create)
	group.POST("/apply", h.apply)
	group.POST("/apply/approve", h.approve)
}

func (h *Handler) index(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	index, err := h.service.GetIndex(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5001, "failed to load alliance index", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(index, middleware.GetTraceID(c)))
}

func (h *Handler) hall(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	hall, err := h.service.ListHall(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5002, "failed to load alliance hall", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(hall, middleware.GetTraceID(c)))
}

func (h *Handler) applyList(c *gin.Context) {
	playerID, ok := parsePlayerID(c)
	if !ok {
		return
	}

	applications, err := h.service.ListApplications(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, ErrAlliancePermissionDenied) {
			c.JSON(stdhttp.StatusForbidden, httpx.Error(4031, "alliance permission denied", middleware.GetTraceID(c)))
			return
		}
		if errors.Is(err, ErrAllianceNotFound) {
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "alliance not found", middleware.GetTraceID(c)))
			return
		}
		c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5003, "failed to load alliance applications", middleware.GetTraceID(c)))
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(applications, middleware.GetTraceID(c)))
}

func (h *Handler) create(c *gin.Context) {
	var req createAllianceRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return
	}

	index, err := h.service.CreateAlliance(c.Request.Context(), req.PlayerID, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrAllianceNameRequired):
			c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, "alliance name is required", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAllianceNameTaken):
			c.JSON(stdhttp.StatusConflict, httpx.Error(4091, "alliance name already exists", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAlreadyInAlliance):
			c.JSON(stdhttp.StatusConflict, httpx.Error(4092, "player already joined alliance", middleware.GetTraceID(c)))
		default:
			c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5004, "failed to create alliance", middleware.GetTraceID(c)))
		}
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(index, middleware.GetTraceID(c)))
}

func (h *Handler) apply(c *gin.Context) {
	var req applyAllianceRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 || req.AllianceID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4004, "player_id and alliance_id are required", middleware.GetTraceID(c)))
		return
	}

	if err := h.service.Apply(c.Request.Context(), req.PlayerID, req.AllianceID); err != nil {
		switch {
		case errors.Is(err, ErrAllianceNotFound):
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "alliance not found", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAlreadyInAlliance):
			c.JSON(stdhttp.StatusConflict, httpx.Error(4092, "player already joined alliance", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAlreadyApplied):
			c.JSON(stdhttp.StatusConflict, httpx.Error(4093, "alliance application already exists", middleware.GetTraceID(c)))
		default:
			c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5005, "failed to apply alliance", middleware.GetTraceID(c)))
		}
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"status": "applied"}, middleware.GetTraceID(c)))
}

func (h *Handler) approve(c *gin.Context) {
	var req approveAllianceRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, "invalid request", middleware.GetTraceID(c)))
		return
	}
	if req.PlayerID == 0 || req.ApplicantPlayerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4005, "player_id and applicant_player_id are required", middleware.GetTraceID(c)))
		return
	}

	index, err := h.service.ApproveApplication(c.Request.Context(), req.PlayerID, req.ApplicantPlayerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrAllianceNotFound):
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4041, "alliance not found", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAlliancePermissionDenied):
			c.JSON(stdhttp.StatusForbidden, httpx.Error(4031, "alliance permission denied", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAllianceApplicationMissing):
			c.JSON(stdhttp.StatusNotFound, httpx.Error(4042, "alliance application not found", middleware.GetTraceID(c)))
		case errors.Is(err, ErrAllianceFull):
			c.JSON(stdhttp.StatusConflict, httpx.Error(4094, "alliance is full", middleware.GetTraceID(c)))
		default:
			c.JSON(stdhttp.StatusInternalServerError, httpx.Error(5006, "failed to approve alliance application", middleware.GetTraceID(c)))
		}
		return
	}

	c.JSON(stdhttp.StatusOK, httpx.Success(index, middleware.GetTraceID(c)))
}

func parsePlayerID(c *gin.Context) (int64, bool) {
	playerID, err := strconv.ParseInt(c.Query("player_id"), 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "player_id is required", middleware.GetTraceID(c)))
		return 0, false
	}
	return playerID, true
}
