# Observability

## ADDED Requirements

### Requirement: Trace Management
The system SHALL maintain distributed traces across all operations.

#### Scenario: Trace context propagation
- **WHEN** any API request is made with X-Trace-Id
- **THEN** the system propagates trace context to all downstream operations

#### Scenario: New trace creation
- **WHEN** a new intent is created without trace_id
- **THEN** the system generates a new trace_id and associates with all related entities

### Requirement: Event Logging
The system SHALL log all significant events.

#### Scenario: State transition logged
- **WHEN** any entity changes state
- **THEN** the system logs: entity_type, entity_id, from_state, to_state, trigger, timestamp

#### Scenario: Agent action logged
- **WHEN** an agent performs an action (bid, execute, submit)
- **THEN** the system logs: agent_id, action_type, target_entity, outcome, timestamp

### Requirement: Metrics Collection
The system SHALL collect and expose operational metrics.

#### Scenario: Task metrics available
- **WHEN** metrics are queried
- **THEN** the system exposes: task success rate, average completion time, bid acceptance rate

#### Scenario: Agent metrics available
- **WHEN** metrics are queried
- **THEN** the system exposes: active agents, reputation distribution, task throughput per agent

### Requirement: Audit Trail
The system SHALL maintain immutable audit trails.

#### Scenario: Query audit trail
- **WHEN** an authorized entity queries GET /api/v1/audit?entity_id={id}
- **THEN** the system returns all actions performed on the entity, with actor, timestamp, details

### Requirement: Alert Generation
The system SHALL generate alerts for abnormal conditions.

#### Scenario: Intent timeout alert
- **WHEN** an intent exceeds expected_completed_at
- **THEN** the system generates alert with severity 'warning', notifies governors

#### Scenario: Reputation anomaly alert
- **WHEN** an agent's reputation changes significantly
- **THEN** the system generates alert with details, triggers review if threshold exceeded

#### Scenario: Task failure cascade alert
- **WHEN** 3+ tasks fail in the same intent
- **THEN** the system generates alert with severity 'critical', notifies governors

### Requirement: Health Monitoring
The system SHALL expose health check endpoints.

#### Scenario: Health check request
- **WHEN** system receives GET /health
- **THEN** the system returns: status, database_health, queue_health, uptime

### Requirement: Performance Tracking
The system SHALL track API performance.

#### Scenario: API latency recorded
- **WHEN** an API request completes
- **THEN** the system records: endpoint, latency_ms, status_code, agent_id (if applicable)
