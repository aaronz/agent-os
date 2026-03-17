package repository

import (
	"database/sql"
	"fmt"

	"github.com/agent-os/core/internal/config"
	"github.com/agent-os/core/internal/models"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{db}, nil
}

func (db *DB) InitSchema() error {
	schema := `
	-- Organizations
	CREATE TABLE IF NOT EXISTS organizations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) UNIQUE NOT NULL,
		description TEXT,
		owner VARCHAR(255) NOT NULL,
		governance_rules JSONB DEFAULT '[]',
		status VARCHAR(50) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Agents
	CREATE TABLE IF NOT EXISTS agents (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		name VARCHAR(255) NOT NULL,
		role VARCHAR(100),
		capabilities JSONB DEFAULT '[]',
		reputation INT DEFAULT 50,
		memory_refs JSONB DEFAULT '[]',
		status VARCHAR(50) DEFAULT 'idle',
		api_key_hash VARCHAR(255),
		created_by VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		last_active_at TIMESTAMP DEFAULT NOW(),
		busy_task_count INT DEFAULT 0,
		abandoned_count INT DEFAULT 0,
		overturned_reviews INT DEFAULT 0
	);

	-- Intents
	CREATE TABLE IF NOT EXISTS intents (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		trace_id UUID DEFAULT gen_random_uuid(),
		title VARCHAR(500) NOT NULL,
		description TEXT NOT NULL,
		constraints JSONB DEFAULT '[]',
		success_criteria JSONB NOT NULL,
		priority VARCHAR(20) DEFAULT 'medium',
		created_by VARCHAR(255) NOT NULL,
		status VARCHAR(50) DEFAULT 'draft',
		plan_ref UUID,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		expected_completed_at TIMESTAMP,
		actual_completed_at TIMESTAMP
	);

	-- Task Graphs
	CREATE TABLE IF NOT EXISTS task_graphs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		intent_id UUID REFERENCES intents(id),
		org_id UUID REFERENCES organizations(id),
		trace_id UUID,
		tasks JSONB NOT NULL,
		dependencies JSONB DEFAULT '[]',
		created_by VARCHAR(255) NOT NULL,
		version INT DEFAULT 1,
		status VARCHAR(50) DEFAULT 'draft',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Tasks
	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		graph_id UUID REFERENCES task_graphs(id),
		intent_id UUID REFERENCES intents(id),
		org_id UUID REFERENCES organizations(id),
		trace_id UUID,
		title VARCHAR(500) NOT NULL,
		description TEXT,
		required_capabilities JSONB NOT NULL,
		acceptance_criteria JSONB DEFAULT '[]',
		dependencies JSONB DEFAULT '[]',
		priority VARCHAR(20) DEFAULT 'medium',
		estimated_duration_min INT,
		max_execution_time_min INT DEFAULT 60,
		assigned_agent_id UUID REFERENCES agents(id),
		status VARCHAR(50) DEFAULT 'pending',
		bid_winner_rule VARCHAR(100) DEFAULT 'composite_score',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		deadline_at TIMESTAMP,
		attempt_count INT DEFAULT 0,
		rejection_count INT DEFAULT 0,
		bidding_end_time TIMESTAMP
	);

	-- Artifacts
	CREATE TABLE IF NOT EXISTS artifacts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		task_id UUID REFERENCES tasks(id),
		intent_id UUID REFERENCES intents(id),
		trace_id UUID,
		type VARCHAR(50) NOT NULL,
		title VARCHAR(500) NOT NULL,
		description TEXT,
		content_ref TEXT NOT NULL,
		content_hash VARCHAR(255),
		dependencies JSONB DEFAULT '[]',
		created_by VARCHAR(255) NOT NULL,
		version INT DEFAULT 1,
		status VARCHAR(50) DEFAULT 'pending_review',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Reviews
	CREATE TABLE IF NOT EXISTS reviews (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		artifact_id UUID REFERENCES artifacts(id),
		task_id UUID REFERENCES tasks(id),
		intent_id UUID REFERENCES intents(id),
		trace_id UUID,
		reviewer_agent_id UUID REFERENCES agents(id),
		score INT NOT NULL,
		is_approved BOOLEAN NOT NULL,
		comments TEXT,
		rejection_reason VARCHAR(100),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Memories
	CREATE TABLE IF NOT EXISTS memories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		type VARCHAR(50) NOT NULL,
		title VARCHAR(500) NOT NULL,
		content TEXT NOT NULL,
		embedding_id VARCHAR(255),
		related_entities JSONB DEFAULT '{}',
		source VARCHAR(100),
		validity VARCHAR(20) DEFAULT 'valid',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		last_retrieved_at TIMESTAMP DEFAULT NOW(),
		retrieval_count INT DEFAULT 0,
		citation_count INT DEFAULT 0
	);

	-- Bids
	CREATE TABLE IF NOT EXISTS bids (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		task_id UUID REFERENCES tasks(id),
		agent_id UUID REFERENCES agents(id),
		estimated_time_min INT,
		estimated_cost INT DEFAULT 0,
		confidence INT,
		proposal TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		status VARCHAR(50) DEFAULT 'pending'
	);

	-- Arbitrations
	CREATE TABLE IF NOT EXISTS arbitrations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		type VARCHAR(50) NOT NULL,
		applicant_id UUID REFERENCES agents(id),
		respondent_id UUID REFERENCES agents(id),
		related_entity_ids JSONB DEFAULT '{}',
		claim TEXT NOT NULL,
		evidence JSONB DEFAULT '[]',
		arbitrator_agent_id UUID REFERENCES agents(id),
		ruling TEXT,
		is_applicant_win BOOLEAN,
		penalty_decision JSONB DEFAULT '[]',
		is_final BOOLEAN DEFAULT true,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT NOW(),
		ruled_at TIMESTAMP
	);

	-- Reputation Events
	CREATE TABLE IF NOT EXISTS reputation_events (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		agent_id UUID REFERENCES agents(id),
		type VARCHAR(50) NOT NULL,
		points INT NOT NULL,
		reason VARCHAR(255),
		entity_id UUID,
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Violations
	CREATE TABLE IF NOT EXISTS violations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		org_id UUID REFERENCES organizations(id),
		agent_id UUID REFERENCES agents(id),
		type VARCHAR(50) NOT NULL,
		details JSONB DEFAULT '{}',
		penalty INT DEFAULT 0,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Audit Logs
	CREATE TABLE IF NOT EXISTS audit_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		entity_type VARCHAR(50) NOT NULL,
		entity_id UUID NOT NULL,
		actor_id VARCHAR(255) NOT NULL,
		actor_type VARCHAR(20) NOT NULL,
		action VARCHAR(100) NOT NULL,
		details JSONB DEFAULT '{}',
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_agents_org_id ON agents(org_id);
	CREATE INDEX IF NOT EXISTS idx_agents_status ON agents(status);
	CREATE INDEX IF NOT EXISTS idx_intents_org_id ON intents(org_id);
	CREATE INDEX IF NOT EXISTS idx_intents_status ON intents(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_org_id ON tasks(org_id);
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_intent_id ON tasks(intent_id);
	CREATE INDEX IF NOT EXISTS idx_artifacts_task_id ON artifacts(task_id);
	CREATE INDEX IF NOT EXISTS idx_reviews_artifact_id ON reviews(artifact_id);
	CREATE INDEX IF NOT EXISTS idx_memories_org_id ON memories(org_id);
	CREATE INDEX IF NOT EXISTS idx_memories_type ON memories(type);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
	`

	_, err := db.Exec(schema)
	return err
}

// Organization Repository
type OrganizationRepository struct {
	db *DB
}

func NewOrganizationRepository(db *DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(org *models.Organization) error {
	query := `
		INSERT INTO organizations (id, name, description, owner, governance_rules, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, org.ID, org.Name, org.Description, org.Owner, org.GovernanceRules, org.Status, org.CreatedAt, org.UpdatedAt)
	return err
}

func (r *OrganizationRepository) Get(id string) (*models.Organization, error) {
	query := `SELECT id, name, description, owner, governance_rules, status, created_at, updated_at FROM organizations WHERE id = $1`
	org := &models.Organization{}
	err := r.db.QueryRow(query, id).Scan(&org.ID, &org.Name, &org.Description, &org.Owner, &org.GovernanceRules, &org.Status, &org.CreatedAt, &org.UpdatedAt)
	return org, err
}

func (r *OrganizationRepository) List() ([]*models.Organization, error) {
	query := `SELECT id, name, description, owner, governance_rules, status, created_at, updated_at FROM organizations WHERE status = 'active'`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*models.Organization
	for rows.Next() {
		org := &models.Organization{}
		if err := rows.Scan(&org.ID, &org.Name, &org.Description, &org.Owner, &org.GovernanceRules, &org.Status, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (r *OrganizationRepository) Update(org *models.Organization) error {
	query := `UPDATE organizations SET name=$1, description=$2, governance_rules=$3, status=$4, updated_at=$5 WHERE id=$6`
	_, err := r.db.Exec(query, org.Name, org.Description, org.GovernanceRules, org.Status, org.UpdatedAt, org.ID)
	return err
}

// Agent Repository
type AgentRepository struct {
	db *DB
}

func NewAgentRepository(db *DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) Create(agent *models.Agent) error {
	query := `
		INSERT INTO agents (id, org_id, name, role, capabilities, reputation, memory_refs, status, api_key_hash, created_by, created_at, updated_at, last_active_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(query, agent.ID, agent.OrgID, agent.Name, agent.Role, agent.Capabilities, agent.Reputation, agent.MemoryRefs, agent.Status, agent.APIKeyHash, agent.CreatedBy, agent.CreatedAt, agent.UpdatedAt, agent.LastActiveAt)
	return err
}

func (r *AgentRepository) Get(id string) (*models.Agent, error) {
	query := `SELECT id, org_id, name, role, capabilities, reputation, memory_refs, status, api_key_hash, created_by, created_at, updated_at, last_active_at, busy_task_count, abandoned_count, overturned_reviews FROM agents WHERE id = $1`
	agent := &models.Agent{}
	err := r.db.QueryRow(query, id).Scan(&agent.ID, &agent.OrgID, &agent.Name, &agent.Role, &agent.Capabilities, &agent.Reputation, &agent.MemoryRefs, &agent.Status, &agent.APIKeyHash, &agent.CreatedBy, &agent.CreatedAt, &agent.UpdatedAt, &agent.LastActiveAt, &agent.BusyTaskCount, &agent.AbandonedCount, &agent.OverturnedReviews)
	return agent, err
}

func (r *AgentRepository) List(orgID string, role, status string, minReputation int) ([]*models.Agent, error) {
	query := `SELECT id, org_id, name, role, capabilities, reputation, memory_refs, status, api_key_hash, created_by, created_at, updated_at, last_active_at, busy_task_count, abandoned_count, overturned_reviews FROM agents WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if orgID != "" {
		query += fmt.Sprintf(" AND org_id = $%d", argIdx)
		args = append(args, orgID)
		argIdx++
	}
	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, role)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if minReputation > 0 {
		query += fmt.Sprintf(" AND reputation >= $%d", argIdx)
		args = append(args, minReputation)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		agent := &models.Agent{}
		if err := rows.Scan(&agent.ID, &agent.OrgID, &agent.Name, &agent.Role, &agent.Capabilities, &agent.Reputation, &agent.MemoryRefs, &agent.Status, &agent.APIKeyHash, &agent.CreatedBy, &agent.CreatedAt, &agent.UpdatedAt, &agent.LastActiveAt, &agent.BusyTaskCount, &agent.AbandonedCount, &agent.OverturnedReviews); err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}
	return agents, nil
}

func (r *AgentRepository) Update(agent *models.Agent) error {
	query := `UPDATE agents SET name=$1, role=$2, capabilities=$3, reputation=$4, memory_refs=$5, status=$6, updated_at=$7, last_active_at=$8, busy_task_count=$9, abandoned_count=$10, overturned_reviews=$11 WHERE id=$12`
	_, err := r.db.Exec(query, agent.Name, agent.Role, agent.Capabilities, agent.Reputation, agent.MemoryRefs, agent.Status, agent.UpdatedAt, agent.LastActiveAt, agent.BusyTaskCount, agent.AbandonedCount, agent.OverturnedReviews, agent.ID)
	return err
}

func (r *AgentRepository) GetByAPIKeyHash(hash string) (*models.Agent, error) {
	query := `SELECT id, org_id, name, role, capabilities, reputation, memory_refs, status, api_key_hash, created_by, created_at, updated_at, last_active_at, busy_task_count, abandoned_count, overturned_reviews FROM agents WHERE api_key_hash = $1`
	agent := &models.Agent{}
	err := r.db.QueryRow(query, hash).Scan(&agent.ID, &agent.OrgID, &agent.Name, &agent.Role, &agent.Capabilities, &agent.Reputation, &agent.MemoryRefs, &agent.Status, &agent.APIKeyHash, &agent.CreatedBy, &agent.CreatedAt, &agent.UpdatedAt, &agent.LastActiveAt, &agent.BusyTaskCount, &agent.AbandonedCount, &agent.OverturnedReviews)
	return agent, err
}

func (r *AgentRepository) UpdateReputation(id string, reputation int) error {
	query := `UPDATE agents SET reputation=$1, updated_at=NOW() WHERE id=$2`
	_, err := r.db.Exec(query, reputation, id)
	return err
}

// Intent Repository
type IntentRepository struct {
	db *DB
}

func NewIntentRepository(db *DB) *IntentRepository {
	return &IntentRepository{db: db}
}

func (r *IntentRepository) Create(intent *models.Intent) error {
	query := `
		INSERT INTO intents (id, org_id, trace_id, title, description, constraints, success_criteria, priority, created_by, status, plan_ref, created_at, updated_at, expected_completed_at, actual_completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.db.Exec(query, intent.ID, intent.OrgID, intent.TraceID, intent.Title, intent.Description, intent.Constraints, intent.SuccessCriteria, intent.Priority, intent.CreatedBy, intent.Status, intent.PlanRef, intent.CreatedAt, intent.UpdatedAt, intent.ExpectedCompletedAt, intent.ActualCompletedAt)
	return err
}

func (r *IntentRepository) Get(id string) (*models.Intent, error) {
	query := `SELECT id, org_id, trace_id, title, description, constraints, success_criteria, priority, created_by, status, plan_ref, created_at, updated_at, expected_completed_at, actual_completed_at FROM intents WHERE id = $1`
	intent := &models.Intent{}
	err := r.db.QueryRow(query, id).Scan(&intent.ID, &intent.OrgID, &intent.TraceID, &intent.Title, &intent.Description, &intent.Constraints, &intent.SuccessCriteria, &intent.Priority, &intent.CreatedBy, &intent.Status, &intent.PlanRef, &intent.CreatedAt, &intent.UpdatedAt, &intent.ExpectedCompletedAt, &intent.ActualCompletedAt)
	return intent, err
}

func (r *IntentRepository) List(orgID string, status string, priority string) ([]*models.Intent, error) {
	query := `SELECT id, org_id, trace_id, title, description, constraints, success_criteria, priority, created_by, status, plan_ref, created_at, updated_at, expected_completed_at, actual_completed_at FROM intents WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if orgID != "" {
		query += fmt.Sprintf(" AND org_id = $%d", argIdx)
		args = append(args, orgID)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if priority != "" {
		query += fmt.Sprintf(" AND priority = $%d", argIdx)
		args = append(args, priority)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var intents []*models.Intent
	for rows.Next() {
		intent := &models.Intent{}
		if err := rows.Scan(&intent.ID, &intent.OrgID, &intent.TraceID, &intent.Title, &intent.Description, &intent.Constraints, &intent.SuccessCriteria, &intent.Priority, &intent.CreatedBy, &intent.Status, &intent.PlanRef, &intent.CreatedAt, &intent.UpdatedAt, &intent.ExpectedCompletedAt, &intent.ActualCompletedAt); err != nil {
			return nil, err
		}
		intents = append(intents, intent)
	}
	return intents, nil
}

func (r *IntentRepository) Update(intent *models.Intent) error {
	query := `UPDATE intents SET title=$1, description=$2, constraints=$3, success_criteria=$4, priority=$5, status=$6, plan_ref=$7, updated_at=$8, expected_completed_at=$9, actual_completed_at=$10 WHERE id=$11`
	_, err := r.db.Exec(query, intent.Title, intent.Description, intent.Constraints, intent.SuccessCriteria, intent.Priority, intent.Status, intent.PlanRef, intent.UpdatedAt, intent.ExpectedCompletedAt, intent.ActualCompletedAt, intent.ID)
	return err
}

// Task Repository
type TaskRepository struct {
	db *DB
}

func NewTaskRepository(db *DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (id, graph_id, intent_id, org_id, trace_id, title, description, required_capabilities, acceptance_criteria, dependencies, priority, estimated_duration_min, max_execution_time_min, assigned_agent_id, status, bid_winner_rule, created_at, updated_at, deadline_at, attempt_count, rejection_count, bidding_end_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`
	_, err := r.db.Exec(query, task.ID, task.GraphID, task.IntentID, task.OrgID, task.TraceID, task.Title, task.Description, task.RequiredCapabilities, task.AcceptanceCriteria, task.Dependencies, task.Priority, task.EstimatedDurationMin, task.MaxExecutionTimeMin, task.AssignedAgentID, task.Status, task.BidWinnerRule, task.CreatedAt, task.UpdatedAt, task.DeadlineAt, task.AttemptCount, task.RejectionCount, task.BiddingEndTime)
	return err
}

func (r *TaskRepository) Get(id string) (*models.Task, error) {
	query := `SELECT id, graph_id, intent_id, org_id, trace_id, title, description, required_capabilities, acceptance_criteria, dependencies, priority, estimated_duration_min, max_execution_time_min, assigned_agent_id, status, bid_winner_rule, created_at, updated_at, deadline_at, attempt_count, rejection_count, bidding_end_time FROM tasks WHERE id = $1`
	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.GraphID, &task.IntentID, &task.OrgID, &task.TraceID, &task.Title, &task.Description, &task.RequiredCapabilities, &task.AcceptanceCriteria, &task.Dependencies, &task.Priority, &task.EstimatedDurationMin, &task.MaxExecutionTimeMin, &task.AssignedAgentID, &task.Status, &task.BidWinnerRule, &task.CreatedAt, &task.UpdatedAt, &task.DeadlineAt, &task.AttemptCount, &task.RejectionCount, &task.BiddingEndTime)
	return task, err
}

func (r *TaskRepository) List(orgID, status, capabilities string) ([]*models.Task, error) {
	query := `SELECT id, graph_id, intent_id, org_id, trace_id, title, description, required_capabilities, acceptance_criteria, dependencies, priority, estimated_duration_min, max_execution_time_min, assigned_agent_id, status, bid_winner_rule, created_at, updated_at, deadline_at, attempt_count, rejection_count, bidding_end_time FROM tasks WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if orgID != "" {
		query += fmt.Sprintf(" AND org_id = $%d", argIdx)
		args = append(args, orgID)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	// Capabilities filtering would require JSONB query - simplified for now

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		if err := rows.Scan(&task.ID, &task.GraphID, &task.IntentID, &task.OrgID, &task.TraceID, &task.Title, &task.Description, &task.RequiredCapabilities, &task.AcceptanceCriteria, &task.Dependencies, &task.Priority, &task.EstimatedDurationMin, &task.MaxExecutionTimeMin, &task.AssignedAgentID, &task.Status, &task.BidWinnerRule, &task.CreatedAt, &task.UpdatedAt, &task.DeadlineAt, &task.AttemptCount, &task.RejectionCount, &task.BiddingEndTime); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, required_capabilities=$3, acceptance_criteria=$4, dependencies=$5, priority=$6, estimated_duration_min=$7, max_execution_time_min=$8, assigned_agent_id=$9, status=$10, bid_winner_rule=$11, updated_at=$12, deadline_at=$13, attempt_count=$14, rejection_count=$15, bidding_end_time=$16 WHERE id=$17`
	_, err := r.db.Exec(query, task.Title, task.Description, task.RequiredCapabilities, task.AcceptanceCriteria, task.Dependencies, task.Priority, task.EstimatedDurationMin, task.MaxExecutionTimeMin, task.AssignedAgentID, task.Status, task.BidWinnerRule, task.UpdatedAt, task.DeadlineAt, task.AttemptCount, task.RejectionCount, task.BiddingEndTime, task.ID)
	return err
}

// Artifact Repository
type ArtifactRepository struct {
	db *DB
}

func NewArtifactRepository(db *DB) *ArtifactRepository {
	return &ArtifactRepository{db: db}
}

func (r *ArtifactRepository) Create(artifact *models.Artifact) error {
	query := `
		INSERT INTO artifacts (id, org_id, task_id, intent_id, trace_id, type, title, description, content_ref, content_hash, dependencies, created_by, version, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.Exec(query, artifact.ID, artifact.OrgID, artifact.TaskID, artifact.IntentID, artifact.TraceID, artifact.Type, artifact.Title, artifact.Description, artifact.ContentRef, artifact.ContentHash, artifact.Dependencies, artifact.CreatedBy, artifact.Version, artifact.Status, artifact.CreatedAt, artifact.UpdatedAt)
	return err
}

func (r *ArtifactRepository) Get(id string) (*models.Artifact, error) {
	query := `SELECT id, org_id, task_id, intent_id, trace_id, type, title, description, content_ref, content_hash, dependencies, created_by, version, status, created_at, updated_at FROM artifacts WHERE id = $1`
	artifact := &models.Artifact{}
	err := r.db.QueryRow(query, id).Scan(&artifact.ID, &artifact.OrgID, &artifact.TaskID, &artifact.IntentID, &artifact.TraceID, &artifact.Type, &artifact.Title, &artifact.Description, &artifact.ContentRef, &artifact.ContentHash, &artifact.Dependencies, &artifact.CreatedBy, &artifact.Version, &artifact.Status, &artifact.CreatedAt, &artifact.UpdatedAt)
	return artifact, err
}

func (r *ArtifactRepository) ListByTask(taskID string) ([]*models.Artifact, error) {
	query := `SELECT id, org_id, task_id, intent_id, trace_id, type, title, description, content_ref, content_hash, dependencies, created_by, version, status, created_at, updated_at FROM artifacts WHERE task_id = $1 ORDER BY version DESC`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artifacts []*models.Artifact
	for rows.Next() {
		artifact := &models.Artifact{}
		if err := rows.Scan(&artifact.ID, &artifact.OrgID, &artifact.TaskID, &artifact.IntentID, &artifact.TraceID, &artifact.Type, &artifact.Title, &artifact.Description, &artifact.ContentRef, &artifact.ContentHash, &artifact.Dependencies, &artifact.CreatedBy, &artifact.Version, &artifact.Status, &artifact.CreatedAt, &artifact.UpdatedAt); err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}

func (r *ArtifactRepository) Update(artifact *models.Artifact) error {
	query := `UPDATE artifacts SET title=$1, description=$2, content_ref=$3, content_hash=$4, dependencies=$5, version=$6, status=$7, updated_at=$8 WHERE id=$9`
	_, err := r.db.Exec(query, artifact.Title, artifact.Description, artifact.ContentRef, artifact.ContentHash, artifact.Dependencies, artifact.Version, artifact.Status, artifact.UpdatedAt, artifact.ID)
	return err
}

func (r *ArtifactRepository) UpdateStatus(id, status string) error {
	query := `UPDATE artifacts SET status=$1, updated_at=NOW() WHERE id=$2`
	_, err := r.db.Exec(query, status, id)
	return err
}

// Review Repository
type ReviewRepository struct {
	db *DB
}

func NewReviewRepository(db *DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (id, org_id, artifact_id, task_id, intent_id, trace_id, reviewer_agent_id, score, is_approved, comments, rejection_reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(query, review.ID, review.OrgID, review.ArtifactID, review.TaskID, review.IntentID, review.TraceID, review.ReviewerAgentID, review.Score, review.IsApproved, review.Comments, review.RejectionReason, review.CreatedAt, review.UpdatedAt)
	return err
}

func (r *ReviewRepository) Get(id string) (*models.Review, error) {
	query := `SELECT id, org_id, artifact_id, task_id, intent_id, trace_id, reviewer_agent_id, score, is_approved, comments, rejection_reason, created_at, updated_at FROM reviews WHERE id = $1`
	review := &models.Review{}
	err := r.db.QueryRow(query, id).Scan(&review.ID, &review.OrgID, &review.ArtifactID, &review.TaskID, &review.IntentID, &review.TraceID, &review.ReviewerAgentID, &review.Score, &review.IsApproved, &review.Comments, &review.RejectionReason, &review.CreatedAt, &review.UpdatedAt)
	return review, err
}

func (r *ReviewRepository) ListByTask(taskID string) ([]*models.Review, error) {
	query := `SELECT id, org_id, artifact_id, task_id, intent_id, trace_id, reviewer_agent_id, score, is_approved, comments, rejection_reason, created_at, updated_at FROM reviews WHERE task_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		review := &models.Review{}
		if err := rows.Scan(&review.ID, &review.OrgID, &review.ArtifactID, &review.TaskID, &review.IntentID, &review.TraceID, &review.ReviewerAgentID, &review.Score, &review.IsApproved, &review.Comments, &review.RejectionReason, &review.CreatedAt, &review.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// Memory Repository
type MemoryRepository struct {
	db *DB
}

func NewMemoryRepository(db *DB) *MemoryRepository {
	return &MemoryRepository{db: db}
}

func (r *MemoryRepository) Create(memory *models.Memory) error {
	query := `
		INSERT INTO memories (id, org_id, type, title, content, embedding_id, related_entities, source, validity, created_at, updated_at, last_retrieved_at, retrieval_count, citation_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.Exec(query, memory.ID, memory.OrgID, memory.Type, memory.Title, memory.Content, memory.EmbeddingID, memory.RelatedEntities, memory.Source, memory.Validity, memory.CreatedAt, memory.UpdatedAt, memory.LastRetrievedAt, memory.RetrievalCount, memory.CitationCount)
	return err
}

func (r *MemoryRepository) Get(id string) (*models.Memory, error) {
	query := `SELECT id, org_id, type, title, content, embedding_id, related_entities, source, validity, created_at, updated_at, last_retrieved_at, retrieval_count, citation_count FROM memories WHERE id = $1`
	memory := &models.Memory{}
	err := r.db.QueryRow(query, id).Scan(&memory.ID, &memory.OrgID, &memory.Type, &memory.Title, &memory.Content, &memory.EmbeddingID, &memory.RelatedEntities, &memory.Source, &memory.Validity, &memory.CreatedAt, &memory.UpdatedAt, &memory.LastRetrievedAt, &memory.RetrievalCount, &memory.CitationCount)
	return memory, err
}

func (r *MemoryRepository) List(orgID, memType string) ([]*models.Memory, error) {
	query := `SELECT id, org_id, type, title, content, embedding_id, related_entities, source, validity, created_at, updated_at, last_retrieved_at, retrieval_count, citation_count FROM memories WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if orgID != "" {
		query += fmt.Sprintf(" AND org_id = $%d", argIdx)
		args = append(args, orgID)
		argIdx++
	}
	if memType != "" {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, memType)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*models.Memory
	for rows.Next() {
		memory := &models.Memory{}
		if err := rows.Scan(&memory.ID, &memory.OrgID, &memory.Type, &memory.Title, &memory.Content, &memory.EmbeddingID, &memory.RelatedEntities, &memory.Source, &memory.Validity, &memory.CreatedAt, &memory.UpdatedAt, &memory.LastRetrievedAt, &memory.RetrievalCount, &memory.CitationCount); err != nil {
			return nil, err
		}
		memories = append(memories, memory)
	}
	return memories, nil
}

func (r *MemoryRepository) Update(memory *models.Memory) error {
	query := `UPDATE memories SET title=$1, content=$2, embedding_id=$3, related_entities=$4, validity=$5, updated_at=$6, last_retrieved_at=$7, retrieval_count=$8, citation_count=$9 WHERE id=$10`
	_, err := r.db.Exec(query, memory.Title, memory.Content, memory.EmbeddingID, memory.RelatedEntities, memory.Validity, memory.UpdatedAt, memory.LastRetrievedAt, memory.RetrievalCount, memory.CitationCount, memory.ID)
	return err
}

// Bid Repository
type BidRepository struct {
	db *DB
}

func NewBidRepository(db *DB) *BidRepository {
	return &BidRepository{db: db}
}

func (r *BidRepository) Create(bid *models.Bid) error {
	query := `
		INSERT INTO bids (id, org_id, task_id, agent_id, estimated_time_min, estimated_cost, confidence, proposal, created_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, bid.ID, bid.OrgID, bid.TaskID, bid.AgentID, bid.EstimatedTimeMin, bid.EstimatedCost, bid.Confidence, bid.Proposal, bid.CreatedAt, bid.Status)
	return err
}

func (r *BidRepository) GetBidsByTask(taskID string) ([]*models.Bid, error) {
	query := `SELECT id, org_id, task_id, agent_id, estimated_time_min, estimated_cost, confidence, proposal, created_at, status FROM bids WHERE task_id = $1 AND status = 'pending'`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.Bid
	for rows.Next() {
		bid := &models.Bid{}
		if err := rows.Scan(&bid.ID, &bid.OrgID, &bid.TaskID, &bid.AgentID, &bid.EstimatedTimeMin, &bid.EstimatedCost, &bid.Confidence, &bid.Proposal, &bid.CreatedAt, &bid.Status); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (r *BidRepository) UpdateStatus(id, status string) error {
	query := `UPDATE bids SET status=$1 WHERE id=$2`
	_, err := r.db.Exec(query, status, id)
	return err
}

type ArbitrationRepository struct {
	db *DB
}

func NewArbitrationRepository(db *DB) *ArbitrationRepository {
	return &ArbitrationRepository{db: db}
}

func (r *ArbitrationRepository) Create(arb *models.Arbitration) error {
	query := `
		INSERT INTO arbitrations (id, org_id, type, applicant_id, respondent_id, related_entity_ids, claim, evidence, arbitrator_agent_id, ruling, is_applicant_win, penalty_decision, is_final, status, created_at, ruled_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.Exec(query, arb.ID, arb.OrgID, arb.Type, arb.ApplicantID, arb.RespondentID, arb.RelatedEntityIDs, arb.Claim, arb.Evidence, arb.ArbitratorAgentID, arb.Ruling, arb.IsApplicantWin, arb.PenaltyDecision, arb.IsFinal, arb.Status, arb.CreatedAt, arb.RuledAt)
	return err
}

func (r *ArbitrationRepository) Get(id string) (*models.Arbitration, error) {
	query := `SELECT id, org_id, type, applicant_id, respondent_id, related_entity_ids, claim, evidence, arbitrator_agent_id, ruling, is_applicant_win, penalty_decision, is_final, status, created_at, ruled_at FROM arbitrations WHERE id = $1`
	arb := &models.Arbitration{}
	err := r.db.QueryRow(query, id).Scan(&arb.ID, &arb.OrgID, &arb.Type, &arb.ApplicantID, &arb.RespondentID, &arb.RelatedEntityIDs, &arb.Claim, &arb.Evidence, &arb.ArbitratorAgentID, &arb.Ruling, &arb.IsApplicantWin, &arb.PenaltyDecision, &arb.IsFinal, &arb.Status, &arb.CreatedAt, &arb.RuledAt)
	return arb, err
}

func (r *ArbitrationRepository) ListPending(agentID string) ([]*models.Arbitration, error) {
	query := `SELECT id, org_id, type, applicant_id, respondent_id, related_entity_ids, claim, evidence, arbitrator_agent_id, ruling, is_applicant_win, penalty_decision, is_final, status, created_at, ruled_at FROM arbitrations WHERE status = 'pending'`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var arbitrations []*models.Arbitration
	for rows.Next() {
		arb := &models.Arbitration{}
		if err := rows.Scan(&arb.ID, &arb.OrgID, &arb.Type, &arb.ApplicantID, &arb.RespondentID, &arb.RelatedEntityIDs, &arb.Claim, &arb.Evidence, &arb.ArbitratorAgentID, &arb.Ruling, &arb.IsApplicantWin, &arb.PenaltyDecision, &arb.IsFinal, &arb.Status, &arb.CreatedAt, &arb.RuledAt); err != nil {
			return nil, err
		}
		arbitrations = append(arbitrations, arb)
	}
	return arbitrations, nil
}

func (r *ArbitrationRepository) Update(arb *models.Arbitration) error {
	query := `UPDATE arbitrations SET arbitrator_agent_id=$1, ruling=$2, is_applicant_win=$3, penalty_decision=$4, is_final=$5, status=$6, ruled_at=$7 WHERE id=$8`
	_, err := r.db.Exec(query, arb.ArbitratorAgentID, arb.Ruling, arb.IsApplicantWin, arb.PenaltyDecision, arb.IsFinal, arb.Status, arb.RuledAt, arb.ID)
	return err
}
