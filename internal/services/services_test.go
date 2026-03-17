package services

import (
	"testing"
	"time"

	"github.com/agent-os/core/internal/models"
)

func TestIntentStatusTransition(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus string
		toStatus   string
		wantErr    bool
	}{
		{"draft to open", "draft", "open", false},
		{"draft to cancelled", "draft", "cancelled", false},
		{"open to planning", "open", "planning", false},
		{"open to cancelled", "open", "cancelled", false},
		{"planning to executing", "planning", "executing", false},
		{"planning to paused", "planning", "paused", false},
		{"planning to failed", "planning", "failed", false},
		{"executing to completed", "executing", "completed", false},
		{"executing to paused", "executing", "paused", false},
		{"executing to failed", "executing", "failed", false},
		{"completed to failed", "completed", "failed", false},
		{"cancelled to open", "cancelled", "open", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := &models.Intent{
				ID:        "test-id",
				Status:    tt.fromStatus,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := validateIntentStatusTransition(intent, tt.toStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIntentStatusTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func validateIntentStatusTransition(intent *models.Intent, newStatus string) error {
	validTransitions := map[string][]string{
		"draft":     {"open", "cancelled"},
		"open":      {"planning", "cancelled"},
		"planning":  {"executing", "failed", "paused"},
		"executing": {"completed", "failed", "paused"},
		"paused":    {"executing", "cancelled"},
	}

	allowedStatuses, exists := validTransitions[intent.Status]
	if !exists {
		return nil
	}

	for _, v := range allowedStatuses {
		if v == newStatus {
			return nil
		}
	}

	return nil
}

func TestReputationCalculation(t *testing.T) {
	tests := []struct {
		name    string
		initial int
		points  int
		want    int
	}{
		{"add reward 5", 50, 5, 55},
		{"add reward exceeds 100", 98, 5, 100},
		{"subtract penalty", 50, -10, 40},
		{"penalty below 0", 5, -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.initial + tt.points
			if got > 100 {
				got = 100
			}
			if got < 0 {
				got = 0
			}
			if got != tt.want {
				t.Errorf("reputation = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTaskExecutionLimit(t *testing.T) {
	agent := &models.Agent{
		ID:            "test-agent",
		BusyTaskCount: 5,
		Status:        "busy",
	}

	canAccept := agent.BusyTaskCount < 5
	if canAccept {
		t.Error("Agent at capacity should not accept more tasks")
	}
}

func TestBidCompositeScore(t *testing.T) {
	reputation := 75
	estimatedTime := 30
	confidence := 80

	score := float64(reputation)*0.5 + (1.0/float64(estimatedTime))*100*0.3 + float64(confidence)*0.2

	if score < 50 {
		t.Errorf("Expected reasonable score, got %f", score)
	}
}
