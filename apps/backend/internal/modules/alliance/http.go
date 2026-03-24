package alliance

import (
	stdhttp "net/http"
	"strings"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/create", h.create)
	group.POST("/apply", h.apply)
	group.POST("/approve", h.approve)
	group.POST("/notice", h.notice)
	group.POST("/donate", h.donate)
	group.POST("/buildings/upgrade", h.upgradeBuilding)
	group.POST("/fire-training/start", h.startFireTraining)
	group.POST("/fire-training/claim", h.claimFireTraining)
	group.POST("/war/sign", h.signWar)
	group.POST("/war/checkin", h.checkinWar)
	group.POST("/war/match", h.matchWar)
	group.POST("/war/settle", h.settleWar)
	group.POST("/war/redeem", h.redeemWar)
}

func (h *Handler) create(c *gin.Context) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	if err := h.service.CreateAlliance(c.Request.Context(), req.PlayerID, req.Name); err != nil {
		msg := "failed to create alliance"
		if err == ErrLevelTooLow {
			msg = "player level too low"
		}
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4002, msg, middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"created": true}, middleware.GetTraceID(c)))
}

func (h *Handler) apply(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		PlayerID   int64  `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	err := h.service.ApplyJoin(c.Request.Context(), req.AllianceID, req.PlayerID)
	handleAllianceErr(c, err, gin.H{"applied": err == nil})
}
func (h *Handler) approve(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		PlayerID   int64  `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	err := h.service.ApproveApplication(c.Request.Context(), req.AllianceID, req.PlayerID)
	handleAllianceErr(c, err, gin.H{"approved": err == nil})
}
func (h *Handler) notice(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		Notice     string `json:"notice"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	err := h.service.UpdateNotice(c.Request.Context(), req.AllianceID, req.Notice)
	handleAllianceErr(c, err, gin.H{"updated": err == nil})
}
func (h *Handler) donate(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		Coins      int64  `json:"coins"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.Donate(c.Request.Context(), req.AllianceID, req.Coins)
	handleAllianceErr(c, err, res)
}
func (h *Handler) upgradeBuilding(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		Building   string `json:"building"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.UpgradeBuilding(c.Request.Context(), req.AllianceID, req.Building)
	handleAllianceErr(c, err, res)
}
func (h *Handler) startFireTraining(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		PlayerID   int64  `json:"player_id"`
		RoomID     string `json:"room_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.StartFireTraining(c.Request.Context(), req.AllianceID, req.PlayerID, req.RoomID)
	handleAllianceErr(c, err, res)
}
func (h *Handler) claimFireTraining(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		PlayerID   int64  `json:"player_id"`
		RecordID   string `json:"record_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.ClaimFireTraining(c.Request.Context(), req.AllianceID, req.PlayerID, req.RecordID)
	handleAllianceErr(c, err, res)
}
func (h *Handler) signWar(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.SignAllianceWar(c.Request.Context(), req.AllianceID)
	handleAllianceErr(c, err, res)
}
func (h *Handler) checkinWar(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		PlayerID   int64  `json:"player_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.CheckInAllianceWar(c.Request.Context(), req.AllianceID, req.PlayerID)
	handleAllianceErr(c, err, res)
}
func (h *Handler) matchWar(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		OpponentID string `json:"opponent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.MatchAllianceWar(c.Request.Context(), req.AllianceID, req.OpponentID)
	handleAllianceErr(c, err, res)
}
func (h *Handler) settleWar(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		Won        bool   `json:"won"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.SettleAllianceWar(c.Request.Context(), req.AllianceID, req.Won)
	handleAllianceErr(c, err, res)
}
func (h *Handler) redeemWar(c *gin.Context) {
	var req struct {
		AllianceID string `json:"alliance_id"`
		Points     int64  `json:"points"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badReq(c)
		return
	}
	res, err := h.service.RedeemWarPoints(c.Request.Context(), req.AllianceID, req.Points)
	handleAllianceErr(c, err, res)
}

func badReq(c *gin.Context) {
	c.JSON(stdhttp.StatusBadRequest, httpx.Error(4001, "invalid request body", middleware.GetTraceID(c)))
}

func handleAllianceErr(c *gin.Context, err error, data interface{}) {
	if err != nil {
		msg := strings.TrimSpace(err.Error())
		c.JSON(stdhttp.StatusBadRequest, httpx.Error(4003, msg, middleware.GetTraceID(c)))
		return
	}
	c.JSON(stdhttp.StatusOK, httpx.Success(data, middleware.GetTraceID(c)))
}
