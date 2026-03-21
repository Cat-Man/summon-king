package config

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrDraftNotFound   = errors.New("config draft not found")
	ErrInvalidDraft    = errors.New("config draft is invalid")
	ErrVersionNotFound = errors.New("config version not found")
)

type Detail struct {
	Module           string `json:"module"`
	DraftContent     string `json:"draft_content"`
	PublishedContent string `json:"published_content"`
	PublishedVersion int    `json:"published_version"`
}

type ValidationResult struct {
	Module   string   `json:"module"`
	Valid    bool     `json:"valid"`
	Messages []string `json:"messages"`
}

type PublishLog struct {
	Module        string    `json:"module"`
	Version       int       `json:"version"`
	Action        string    `json:"action"`
	Content       string    `json:"content"`
	Summary       string    `json:"summary"`
	OperatorID    int64     `json:"operator_id"`
	PublishedAt   time.Time `json:"published_at"`
	RollbackToVer int       `json:"rollback_to_ver,omitempty"`
}

type moduleRecord struct {
	draft       string
	published   string
	version     int
	versions    map[int]string
	publishLogs []PublishLog
}

type Service struct {
	mu      sync.Mutex
	modules map[string]*moduleRecord
	now     func() time.Time
}

func NewService() *Service {
	return &Service{
		modules: map[string]*moduleRecord{},
		now:     time.Now,
	}
}

func (s *Service) ListModules(_ context.Context) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	seen := map[string]struct{}{
		"pet_catalog": {},
		"skills":      {},
		"dungeons":    {},
		"cultivation": {},
		"tower":       {},
		"growth":      {},
		"alliance":    {},
		"arena":       {},
		"commerce":    {},
	}
	for module := range s.modules {
		seen[module] = struct{}{}
	}
	modules := make([]string, 0, len(seen))
	for module := range seen {
		modules = append(modules, module)
	}
	sort.Strings(modules)
	return modules
}

func (s *Service) GetDetail(_ context.Context, module string) (Detail, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.detailLocked(strings.TrimSpace(module)), nil
}

func (s *Service) SaveDraft(_ context.Context, module, content string) (Detail, error) {
	module = strings.TrimSpace(module)
	content = strings.TrimSpace(content)
	if module == "" || content == "" {
		return Detail{}, ErrInvalidDraft
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	record := s.ensureModuleLocked(module)
	record.draft = content
	return s.detailLocked(module), nil
}

func (s *Service) Validate(_ context.Context, module string) (ValidationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.modules[strings.TrimSpace(module)]
	if !ok || strings.TrimSpace(record.draft) == "" {
		return ValidationResult{}, ErrDraftNotFound
	}
	if !json.Valid([]byte(record.draft)) {
		return ValidationResult{Module: module, Valid: false, Messages: []string{"draft is not valid json"}}, ErrInvalidDraft
	}
	return ValidationResult{Module: module, Valid: true, Messages: []string{"ok"}}, nil
}

func (s *Service) Publish(ctx context.Context, module string, operatorID int64, summary string) (PublishLog, error) {
	if _, err := s.Validate(ctx, module); err != nil {
		return PublishLog{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	record := s.ensureModuleLocked(strings.TrimSpace(module))
	record.version++
	record.published = record.draft
	record.versions[record.version] = record.published
	log := PublishLog{
		Module:      strings.TrimSpace(module),
		Version:     record.version,
		Action:      "publish",
		Content:     record.published,
		Summary:     strings.TrimSpace(summary),
		OperatorID:  operatorID,
		PublishedAt: s.now(),
	}
	record.publishLogs = append(record.publishLogs, log)
	return log, nil
}

func (s *Service) Rollback(_ context.Context, module string, version int, operatorID int64) (PublishLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureModuleLocked(strings.TrimSpace(module))
	content, ok := record.versions[version]
	if !ok {
		return PublishLog{}, ErrVersionNotFound
	}
	record.version++
	record.draft = content
	record.published = content
	record.versions[record.version] = content
	log := PublishLog{
		Module:        strings.TrimSpace(module),
		Version:       record.version,
		Action:        "rollback",
		Content:       content,
		Summary:       "rollback",
		OperatorID:    operatorID,
		PublishedAt:   s.now(),
		RollbackToVer: version,
	}
	record.publishLogs = append(record.publishLogs, log)
	return log, nil
}

func (s *Service) PublishLogs(_ context.Context, module string) []PublishLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	record := s.ensureModuleLocked(strings.TrimSpace(module))
	logs := make([]PublishLog, len(record.publishLogs))
	copy(logs, record.publishLogs)
	return logs
}

func (s *Service) ensureModuleLocked(module string) *moduleRecord {
	if record, ok := s.modules[module]; ok {
		if record.versions == nil {
			record.versions = map[int]string{}
		}
		return record
	}
	record := &moduleRecord{versions: map[int]string{}}
	s.modules[module] = record
	return record
}

func (s *Service) detailLocked(module string) Detail {
	record := s.ensureModuleLocked(module)
	return Detail{
		Module:           module,
		DraftContent:     record.draft,
		PublishedContent: record.published,
		PublishedVersion: record.version,
	}
}
