package ranking

import (
	"sort"
	"sync"
)

type Entry struct {
	PlayerID int64 `json:"player_id"`
	Score    int64 `json:"score"`
}

type Service struct {
	mu     sync.Mutex
	boards map[string]map[int64]int64
}

func NewService() *Service { return &Service{boards: make(map[string]map[int64]int64)} }

func (s *Service) RecordScore(board string, playerID int64, score int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.boards[board]; !ok {
		s.boards[board] = make(map[int64]int64)
	}
	s.boards[board][playerID] = score
}

func (s *Service) Snapshot(board string) []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries := make([]Entry, 0, len(s.boards[board]))
	for pid, score := range s.boards[board] {
		entries = append(entries, Entry{PlayerID: pid, Score: score})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Score > entries[j].Score })
	return entries
}
