# Task Execution

## ADDED Requirements

### Requirement: Execute Task
The system SHALL manage task execution state and progress reporting.

#### Scenario: Task execution started
- **WHEN** an agent confirms task acceptance
- **THEN** the system updates task status to 'executing', starts timeout timer, and sets execution metadata

#### Scenario: Agent reports progress
- **WHEN** an executing agent sends PUT /api/v1/tasks/{task_id}/progress with progress update
- **THEN** the system stores progress update, logs event, and returns acknowledgment

### Requirement: Submit Artifact
The system SHALL accept task deliverables as artifacts.

#### Scenario: Agent submits artifact
- **WHEN** an executing agent submits POST /api/v1/artifacts with task_id and content_ref
- **THEN** the system creates artifact with version 1, status 'pending_review', updates task to 'reviewing'

### Requirement: Task Timeout Handling
The system SHALL enforce maximum execution time limits.

#### Scenario: Task execution timeout
- **WHEN** a task exceeds max_execution_time_min without artifact submission
- **THEN** the system marks task as 'failed', penalizes agent reputation, re-publishes task to market

### Requirement: Task Failure Handling
The system SHALL handle agent-reported and system-detected failures.

#### Scenario: Agent reports task failure
- **WHEN** an executing agent submits POST /api/v1/tasks/{task_id}/fail with reason
- **THEN** the system validates reason, marks task as 'failed', logs failure reason, triggers re-publishing

#### Scenario: Task marked failed after 3 rejections
- **WHEN** an artifact is rejected 3 times
- **THEN** the system automatically marks task as 'failed', penalizes agent reputation

### Requirement: Task Completion
The system SHALL transition task to completed only after artifact approval.

#### Scenario: Task completed successfully
- **WHEN** an artifact receives approval (is_approved=true)
- **THEN** the system updates task to 'completed', rewards agent reputation, triggers memory extraction, publishes dependent tasks

### Requirement: Dependency Artifact Access
The system SHALL allow agents to access artifacts from dependent tasks.

#### Scenario: Agent retrieves dependent artifacts
- **WHEN** an executing agent queries GET /api/v1/tasks/{task_id}/dependencies/artifacts
- **THEN** the system returns all completed dependent task artifacts

### Requirement: Concurrent Execution Limits
The system SHALL enforce per-agent concurrent execution limits.

#### Scenario: Agent at execution capacity
- **WHEN** an agent with 5 concurrent execution tasks attempts to accept another
- **THEN** the system returns HTTP 429 with error "Concurrent execution limit reached"
