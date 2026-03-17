package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	requestsTotal      uint64
	requestsByEndpoint map[string]*uint64
	errorsTotal        uint64
	latencySum         uint64
	latencyCount       uint64
	activeConnections  int64
	tasksCreated       uint64
	tasksCompleted     uint64
	tasksFailed        uint64
	intentsCreated     uint64
	intentsCompleted   uint64
	agentsActive       uint64
	bidsPlaced         uint64
	reviewsCompleted   uint64
	mu                 sync.RWMutex
}

var globalMetrics = &Metrics{
	requestsByEndpoint: make(map[string]*uint64),
}

func Global() *Metrics {
	return globalMetrics
}

func (m *Metrics) IncRequest(endpoint string) {
	atomic.AddUint64(&m.requestsTotal, 1)

	m.mu.Lock()
	if _, ok := m.requestsByEndpoint[endpoint]; !ok {
		var count uint64
		m.requestsByEndpoint[endpoint] = &count
	}
	m.mu.Unlock()

	atomic.AddUint64(m.requestsByEndpoint[endpoint], 1)
}

func (m *Metrics) IncError() {
	atomic.AddUint64(&m.errorsTotal, 1)
}

func (m *Metrics) RecordLatency(d time.Duration) {
	atomic.AddUint64(&m.latencySum, uint64(d.Milliseconds()))
	atomic.AddUint64(&m.latencyCount, 1)
}

func (m *Metrics) IncActiveConnections(n int64) {
	atomic.AddInt64(&m.activeConnections, n)
}

func (m *Metrics) IncTasksCreated() {
	atomic.AddUint64(&m.tasksCreated, 1)
}

func (m *Metrics) IncTasksCompleted() {
	atomic.AddUint64(&m.tasksCompleted, 1)
}

func (m *Metrics) IncTasksFailed() {
	atomic.AddUint64(&m.tasksFailed, 1)
}

func (m *Metrics) IncIntentsCreated() {
	atomic.AddUint64(&m.intentsCreated, 1)
}

func (m *Metrics) IncIntentsCompleted() {
	atomic.AddUint64(&m.intentsCompleted, 1)
}

func (m *Metrics) SetAgentsActive(n uint64) {
	atomic.StoreUint64(&m.agentsActive, n)
}

func (m *Metrics) IncBidsPlaced() {
	atomic.AddUint64(&m.bidsPlaced, 1)
}

func (m *Metrics) IncReviewsCompleted() {
	atomic.AddUint64(&m.reviewsCompleted, 1)
}

type Snapshot struct {
	RequestsTotal     uint64  `json:"requests_total"`
	ErrorsTotal       uint64  `json:"errors_total"`
	AvgLatencyMs      float64 `json:"avg_latency_ms"`
	ActiveConnections int64   `json:"active_connections"`
	TasksCreated      uint64  `json:"tasks_created"`
	TasksCompleted    uint64  `json:"tasks_completed"`
	TasksFailed       uint64  `json:"tasks_failed"`
	IntentsCreated    uint64  `json:"intents_created"`
	IntentsCompleted  uint64  `json:"intents_completed"`
	AgentsActive      uint64  `json:"agents_active"`
	BidsPlaced        uint64  `json:"bids_placed"`
	ReviewsCompleted  uint64  `json:"reviews_completed"`
	ErrorRate         float64 `json:"error_rate"`
	TaskSuccessRate   float64 `json:"task_success_rate"`
}

func (m *Metrics) Snapshot() Snapshot {
	requests := atomic.LoadUint64(&m.requestsTotal)
	errors := atomic.LoadUint64(&m.errorsTotal)
	latencySum := atomic.LoadUint64(&m.latencySum)
	latencyCount := atomic.LoadUint64(&m.latencyCount)
	activeConns := atomic.LoadInt64(&m.activeConnections)
	tasksCreated := atomic.LoadUint64(&m.tasksCreated)
	tasksCompleted := atomic.LoadUint64(&m.tasksCompleted)
	tasksFailed := atomic.LoadUint64(&m.tasksFailed)

	var avgLatency float64
	if latencyCount > 0 {
		avgLatency = float64(latencySum) / float64(latencyCount)
	}

	var errorRate float64
	if requests > 0 {
		errorRate = float64(errors) / float64(requests) * 100
	}

	var taskSuccessRate float64
	totalTasks := tasksCompleted + tasksFailed
	if totalTasks > 0 {
		taskSuccessRate = float64(tasksCompleted) / float64(totalTasks) * 100
	}

	return Snapshot{
		RequestsTotal:     requests,
		ErrorsTotal:       errors,
		AvgLatencyMs:      avgLatency,
		ActiveConnections: activeConns,
		TasksCreated:      tasksCreated,
		TasksCompleted:    tasksCompleted,
		TasksFailed:       tasksFailed,
		IntentsCreated:    atomic.LoadUint64(&m.intentsCreated),
		IntentsCompleted:  atomic.LoadUint64(&m.intentsCompleted),
		AgentsActive:      atomic.LoadUint64(&m.agentsActive),
		BidsPlaced:        atomic.LoadUint64(&m.bidsPlaced),
		ReviewsCompleted:  atomic.LoadUint64(&m.reviewsCompleted),
		ErrorRate:         errorRate,
		TaskSuccessRate:   taskSuccessRate,
	}
}

func (m *Metrics) GetEndpointCounts() map[string]uint64 {
	result := make(map[string]uint64)
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.requestsByEndpoint {
		result[k] = atomic.LoadUint64(v)
	}
	return result
}
