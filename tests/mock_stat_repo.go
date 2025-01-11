package tests

import (
	"time"

	"github.com/navidrome/navidrome/model"
)

type MockStatRepo struct {
	model.StatRepository
}

func CreateMockStatRepo() *MockStatRepo {
	return &MockStatRepo{}
}

func (s *MockStatRepo) RecordPlay(id string, ts time.Time) error {
	return nil
}
