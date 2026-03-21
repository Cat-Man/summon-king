package gm

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrGMReasonRequired = errors.New("reason and ticket are required")

type GrantCurrencyRequest struct {
	OperatorID int64  `json:"operator_id"`
	PlayerID   int64  `json:"player_id"`
	Currency   string `json:"currency"`
	Amount     int64  `json:"amount"`
	Reason     string `json:"reason"`
	TicketNo   string `json:"ticket_no"`
}

type GrantItemRequest struct {
	OperatorID int64  `json:"operator_id"`
	PlayerID   int64  `json:"player_id"`
	ItemID     string `json:"item_id"`
	Count      int64  `json:"count"`
	Reason     string `json:"reason"`
	TicketNo   string `json:"ticket_no"`
}

type GrantPetRequest struct {
	OperatorID int64  `json:"operator_id"`
	PlayerID   int64  `json:"player_id"`
	PetID      string `json:"pet_id"`
	Reason     string `json:"reason"`
	TicketNo   string `json:"ticket_no"`
}

type MailRequest struct {
	OperatorID int64  `json:"operator_id"`
	PlayerID   int64  `json:"player_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Reason     string `json:"reason"`
	TicketNo   string `json:"ticket_no"`
}

type ResetDailyRequest struct {
	OperatorID int64  `json:"operator_id"`
	PlayerID   int64  `json:"player_id"`
	ResetType  string `json:"reset_type"`
	Reason     string `json:"reason"`
	TicketNo   string `json:"ticket_no"`
}

type PlayerFlagRequest struct {
	OperatorID int64  `json:"operator_id"`
	PlayerID   int64  `json:"player_id"`
	Reason     string `json:"reason"`
	TicketNo   string `json:"ticket_no"`
}

type FixAllianceRequest struct {
	OperatorID int64  `json:"operator_id"`
	PlayerID   int64  `json:"player_id"`
	AllianceID string `json:"alliance_id"`
	Reason     string `json:"reason"`
	TicketNo   string `json:"ticket_no"`
}

type AuditLog struct {
	OperatorID int64     `json:"operator_id"`
	PlayerID   int64     `json:"player_id"`
	Action     string    `json:"action"`
	Reason     string    `json:"reason"`
	TicketNo   string    `json:"ticket_no"`
	Payload    string    `json:"payload"`
	CreatedAt  time.Time `json:"created_at"`
}

type PlayerOpsState struct {
	PlayerID      int64            `json:"player_id"`
	Wallet        map[string]int64 `json:"wallet"`
	Items         map[string]int64 `json:"items"`
	Pets          []string         `json:"pets"`
	Mails         []string         `json:"mails"`
	ResetRecords  []string         `json:"reset_records"`
	AllianceFixes []string         `json:"alliance_fixes"`
	Muted         bool             `json:"muted"`
	Banned        bool             `json:"banned"`
}

type Service struct {
	mu      sync.Mutex
	players map[int64]*PlayerOpsState
	logs    []AuditLog
	now     func() time.Time
}

func NewService() *Service {
	return &Service{players: map[int64]*PlayerOpsState{}, now: time.Now}
}

func (s *Service) GrantCurrency(_ context.Context, req GrantCurrencyRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.Wallet[req.Currency] += req.Amount
	s.logLocked(req.OperatorID, req.PlayerID, "grant_currency", req.Reason, req.TicketNo, req.Currency)
	return nil
}

func (s *Service) GrantItem(_ context.Context, req GrantItemRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.Items[req.ItemID] += req.Count
	s.logLocked(req.OperatorID, req.PlayerID, "grant_item", req.Reason, req.TicketNo, req.ItemID)
	return nil
}

func (s *Service) GrantPet(_ context.Context, req GrantPetRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.Pets = append(state.Pets, req.PetID)
	s.logLocked(req.OperatorID, req.PlayerID, "grant_pet", req.Reason, req.TicketNo, req.PetID)
	return nil
}

func (s *Service) SendMail(_ context.Context, req MailRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.Mails = append(state.Mails, req.Title)
	s.logLocked(req.OperatorID, req.PlayerID, "send_mail", req.Reason, req.TicketNo, req.Title)
	return nil
}

func (s *Service) ResetDaily(_ context.Context, req ResetDailyRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.ResetRecords = append(state.ResetRecords, req.ResetType)
	s.logLocked(req.OperatorID, req.PlayerID, "reset_daily", req.Reason, req.TicketNo, req.ResetType)
	return nil
}

func (s *Service) UnblockPlayer(_ context.Context, req PlayerFlagRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.Banned = false
	s.logLocked(req.OperatorID, req.PlayerID, "unblock_player", req.Reason, req.TicketNo, "")
	return nil
}

func (s *Service) MutePlayer(_ context.Context, req PlayerFlagRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.Muted = true
	s.logLocked(req.OperatorID, req.PlayerID, "mute_player", req.Reason, req.TicketNo, "")
	return nil
}

func (s *Service) BanPlayer(_ context.Context, req PlayerFlagRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.Banned = true
	s.logLocked(req.OperatorID, req.PlayerID, "ban_player", req.Reason, req.TicketNo, "")
	return nil
}

func (s *Service) FixAlliance(_ context.Context, req FixAllianceRequest) error {
	if err := validateOperatorReason(req.Reason, req.TicketNo); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(req.PlayerID)
	state.AllianceFixes = append(state.AllianceFixes, req.AllianceID)
	s.logLocked(req.OperatorID, req.PlayerID, "fix_alliance", req.Reason, req.TicketNo, req.AllianceID)
	return nil
}

func (s *Service) GetPlayerOpsState(_ context.Context, playerID int64) PlayerOpsState {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.ensurePlayerLocked(playerID)
	copyState := *state
	copyState.Wallet = copyIntMap(state.Wallet)
	copyState.Items = copyIntMap(state.Items)
	copyState.Pets = append([]string{}, state.Pets...)
	copyState.Mails = append([]string{}, state.Mails...)
	copyState.ResetRecords = append([]string{}, state.ResetRecords...)
	copyState.AllianceFixes = append([]string{}, state.AllianceFixes...)
	return copyState
}

func (s *Service) ListAuditLogs(_ context.Context) []AuditLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	logs := make([]AuditLog, len(s.logs))
	copy(logs, s.logs)
	return logs
}

func (s *Service) ensurePlayerLocked(playerID int64) *PlayerOpsState {
	if state, ok := s.players[playerID]; ok {
		return state
	}
	state := &PlayerOpsState{
		PlayerID: playerID,
		Wallet:   map[string]int64{},
		Items:    map[string]int64{},
	}
	s.players[playerID] = state
	return state
}

func (s *Service) logLocked(operatorID, playerID int64, action, reason, ticketNo, payload string) {
	s.logs = append(s.logs, AuditLog{
		OperatorID: operatorID,
		PlayerID:   playerID,
		Action:     action,
		Reason:     reason,
		TicketNo:   ticketNo,
		Payload:    payload,
		CreatedAt:  s.now(),
	})
}

func validateOperatorReason(reason, ticketNo string) error {
	if strings.TrimSpace(reason) == "" || strings.TrimSpace(ticketNo) == "" {
		return ErrGMReasonRequired
	}
	return nil
}

func copyIntMap(src map[string]int64) map[string]int64 {
	dst := make(map[string]int64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
