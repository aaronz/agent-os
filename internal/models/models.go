package models

import (
	"time"
)

// Organization represents a tenant in the system
type Organization struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	Owner           string    `json:"owner" db:"owner"`
	GovernanceRules string    `json:"governance_rules" db:"governance_rules"` // JSON string
	Status          string    `json:"status" db:"status"`                     // active, disabled
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Agent represents an AI agent in the system
type Agent struct {
	ID                string    `json:"id" db:"id"`
	OrgID             string    `json:"org_id" db:"org_id"`
	Name              string    `json:"name" db:"name"`
	Role              string    `json:"role" db:"role"`                 // planning_agent, developer, reviewer, arbitrator
	Capabilities      string    `json:"capabilities" db:"capabilities"` // JSON array of capability names
	Reputation        int       `json:"reputation" db:"reputation"`
	MemoryRefs        string    `json:"memory_refs" db:"memory_refs"` // JSON array of memory IDs
	Status            string    `json:"status" db:"status"`           // idle, busy, disabled, terminated
	APIKey            string    `json:"api_key,omitempty" db:"api_key"`
	APIKeyHash        string    `json:"api_key_hash" db:"api_key_hash"`
	CreatedBy         string    `json:"created_by" db:"created_by"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	LastActiveAt      time.Time `json:"last_active_at" db:"last_active_at"`
	BusyTaskCount     int       `json:"busy_task_count" db:"busy_task_count"`
	AbandonedCount    int       `json:"abandoned_count" db:"abandoned_count"`
	OverturnedReviews int       `json:"overturned_reviews" db:"overturned_reviews"`
}

// Capability represents an agent's ability
type Capability struct {
	Name                  string `json:"name" db:"name"`
	Description           string `json:"description" db:"description"`
	Toolset               string `json:"toolset" db:"toolset"` // JSON array
	RequiredMinReputation int    `json:"required_min_reputation" db:"required_min_reputation"`
}

// Intent represents a goal in the system
type Intent struct {
	ID                  string     `json:"id" db:"id"`
	OrgID               string     `json:"org_id" db:"org_id"`
	TraceID             string     `json:"trace_id" db:"trace_id"`
	Title               string     `json:"title" db:"title"`
	Description         string     `json:"description" db:"description"`
	Constraints         string     `json:"constraints" db:"constraints"`           // JSON array
	SuccessCriteria     string     `json:"success_criteria" db:"success_criteria"` // JSON array
	Priority            string     `json:"priority" db:"priority"`                 // high, medium, low
	CreatedBy           string     `json:"created_by" db:"created_by"`
	Status              string     `json:"status" db:"status"`     // draft, open, planning, executing, paused, completed, failed, cancelled
	PlanRef             string     `json:"plan_ref" db:"plan_ref"` // Task Graph ID
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	ExpectedCompletedAt time.Time  `json:"expected_completed_at" db:"expected_completed_at"`
	ActualCompletedAt   *time.Time `json:"actual_completed_at" db:"actual_completed_at"`
}

// TaskGraph represents a DAG of tasks
type TaskGraph struct {
	ID              string    `json:"id" db:"id"`
	IntentID        string    `json:"intent_id" db:"intent_id"`
	OrgID           string    `json:"org_id" db:"org_id"`
	TraceID         string    `json:"trace_id" db:"trace_id"`
	Tasks           string    `json:"tasks" db:"tasks"`
	Dependencies    string    `json:"dependencies" db:"dependencies"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	Version         int       `json:"version" db:"version"`
	PreviousVersion int       `json:"previous_version" db:"previous_version"`
	Status          string    `json:"status" db:"status"`
	ChangeReason    string    `json:"change_reason" db:"change_reason"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Task represents a single executable unit
type Task struct {
	ID                   string     `json:"id" db:"id"`
	GraphID              string     `json:"graph_id" db:"graph_id"`
	IntentID             string     `json:"intent_id" db:"intent_id"`
	OrgID                string     `json:"org_id" db:"org_id"`
	TraceID              string     `json:"trace_id" db:"trace_id"`
	Title                string     `json:"title" db:"title"`
	Description          string     `json:"description" db:"description"`
	RequiredCapabilities string     `json:"required_capabilities" db:"required_capabilities"` // JSON array
	AcceptanceCriteria   string     `json:"acceptance_criteria" db:"acceptance_criteria"`     // JSON array
	Dependencies         string     `json:"dependencies" db:"dependencies"`                   // JSON array of task IDs
	Priority             string     `json:"priority" db:"priority"`
	EstimatedDurationMin int        `json:"estimated_duration_min" db:"estimated_duration_min"`
	MaxExecutionTimeMin  int        `json:"max_execution_time_min" db:"max_execution_time_min"`
	AssignedAgentID      string     `json:"assigned_agent_id" db:"assigned_agent_id"`
	Status               string     `json:"status" db:"status"` // pending, open, bidding, assigned, executing, reviewing, completed, failed, cancelled
	BidWinnerRule        string     `json:"bid_winner_rule" db:"bid_winner_rule"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeadlineAt           *time.Time `json:"deadline_at" db:"deadline_at"`
	AttemptCount         int        `json:"attempt_count" db:"attempt_count"`
	RejectionCount       int        `json:"rejection_count" db:"rejection_count"`
	BiddingEndTime       *time.Time `json:"bidding_end_time" db:"bidding_end_time"`
}

// Artifact represents a deliverable
type Artifact struct {
	ID           string    `json:"id" db:"id"`
	OrgID        string    `json:"org_id" db:"org_id"`
	TaskID       string    `json:"task_id" db:"task_id"`
	IntentID     string    `json:"intent_id" db:"intent_id"`
	TraceID      string    `json:"trace_id" db:"trace_id"`
	Type         string    `json:"type" db:"type"` // code, document, plan, model, dataset, report, config, other
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	ContentRef   string    `json:"content_ref" db:"content_ref"`
	ContentHash  string    `json:"content_hash" db:"content_hash"`
	Dependencies string    `json:"dependencies" db:"dependencies"` // JSON array of artifact IDs
	CreatedBy    string    `json:"created_by" db:"created_by"`
	Version      int       `json:"version" db:"version"`
	Status       string    `json:"status" db:"status"` // pending_review, approved, rejected, deprecated
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Review represents an artifact review
type Review struct {
	ID              string    `json:"id" db:"id"`
	OrgID           string    `json:"org_id" db:"org_id"`
	ArtifactID      string    `json:"artifact_id" db:"artifact_id"`
	TaskID          string    `json:"task_id" db:"task_id"`
	IntentID        string    `json:"intent_id" db:"intent_id"`
	TraceID         string    `json:"trace_id" db:"trace_id"`
	ReviewerAgentID string    `json:"reviewer_agent_id" db:"reviewer_agent_id"`
	Score           int       `json:"score" db:"score"` // 0-100
	IsApproved      bool      `json:"is_approved" db:"is_approved"`
	Comments        string    `json:"comments" db:"comments"`
	RejectionReason string    `json:"rejection_reason" db:"rejection_reason"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Memory represents experience memory
type Memory struct {
	ID              string    `json:"id" db:"id"`
	OrgID           string    `json:"org_id" db:"org_id"`
	Type            string    `json:"type" db:"type"` // knowledge, project, task, failure, best_practice, review
	Title           string    `json:"title" db:"title"`
	Content         string    `json:"content" db:"content"`
	EmbeddingID     string    `json:"embedding_id" db:"embedding_id"`
	RelatedEntities string    `json:"related_entities" db:"related_entities"` // JSON
	Source          string    `json:"source" db:"source"`
	Validity        string    `json:"validity" db:"validity"` // valid, invalid
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	LastRetrievedAt time.Time `json:"last_retrieved_at" db:"last_retrieved_at"`
	RetrievalCount  int       `json:"retrieval_count" db:"retrieval_count"`
	CitationCount   int       `json:"citation_count" db:"citation_count"`
}

// Bid represents a task bid
type Bid struct {
	ID               string    `json:"id" db:"id"`
	OrgID            string    `json:"org_id" db:"org_id"`
	TaskID           string    `json:"task_id" db:"task_id"`
	AgentID          string    `json:"agent_id" db:"agent_id"`
	EstimatedTimeMin int       `json:"estimated_time_min" db:"estimated_time_min"`
	EstimatedCost    int       `json:"estimated_cost" db:"estimated_cost"`
	Confidence       int       `json:"confidence" db:"confidence"` // 0-100
	Proposal         string    `json:"proposal" db:"proposal"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Status           string    `json:"status" db:"status"` // pending, won, lost, cancelled
}

// Arbitration represents a dispute resolution
type Arbitration struct {
	ID                string     `json:"id" db:"id"`
	OrgID             string     `json:"org_id" db:"org_id"`
	Type              string     `json:"type" db:"type"` // review_dispute, bid_dispute, task_conflict, violation_appeal
	ApplicantID       string     `json:"applicant_id" db:"applicant_id"`
	RespondentID      string     `json:"respondent_id" db:"respondent_id"`
	RelatedEntityIDs  string     `json:"related_entity_ids" db:"related_entity_ids"` // JSON
	Claim             string     `json:"claim" db:"claim"`
	Evidence          string     `json:"evidence" db:"evidence"` // JSON array
	ArbitratorAgentID string     `json:"arbitrator_agent_id" db:"arbitrator_agent_id"`
	Ruling            string     `json:"ruling" db:"ruling"`
	IsApplicantWin    bool       `json:"is_applicant_win" db:"is_applicant_win"`
	PenaltyDecision   string     `json:"penalty_decision" db:"penalty_decision"` // JSON array
	IsFinal           bool       `json:"is_final" db:"is_final"`
	Status            string     `json:"status" db:"status"` // pending, ruled, closed
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	RuledAt           *time.Time `json:"ruled_at" db:"ruled_at"`
}

// Event represents a domain event
type Event struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	OrgID      string    `json:"org_id"`
	TraceID    string    `json:"trace_id"`
	ActorID    string    `json:"actor_id"`
	Data       string    `json:"data"` // JSON
	CreatedAt  time.Time `json:"created_at"`
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID         string    `json:"id" db:"id"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	ActorID    string    `json:"actor_id" db:"actor_id"`
	ActorType  string    `json:"actor_type" db:"actor_type"` // human, agent
	Action     string    `json:"action" db:"action"`
	Details    string    `json:"details" db:"details"` // JSON
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ReputationEvent represents a reputation change
type ReputationEvent struct {
	ID        string    `json:"id" db:"id"`
	AgentID   string    `json:"agent_id" db:"agent_id"`
	Type      string    `json:"type" db:"type"` // reward, penalty, decay
	Points    int       `json:"points" db:"points"`
	Reason    string    `json:"reason" db:"reason"`
	EntityID  string    `json:"entity_id" db:"entity_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Violation represents an agent violation
type Violation struct {
	ID        string    `json:"id" db:"id"`
	OrgID     string    `json:"org_id" db:"org_id"`
	AgentID   string    `json:"agent_id" db:"agent_id"`
	Type      string    `json:"type" db:"type"`       // bidding,履约, delivery, review, arbitration, system, compliance
	Details   string    `json:"details" db:"details"` // JSON
	Penalty   int       `json:"penalty" db:"penalty"`
	Status    string    `json:"status" db:"status"` // pending, processed
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
