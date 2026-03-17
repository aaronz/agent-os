# Intent Management

## ADDED Requirements

### Requirement: Create Intent
The system SHALL allow creating a new intent with validated success criteria and constraints.

#### Scenario: Successful intent creation
- **WHEN** an authorized entity sends POST /api/v1/intents with title, description, success_criteria, and constraints
- **THEN** the system validates the intent, creates it with status 'open' if valid, generates unique intent_id and trace_id, and triggers planning

#### Scenario: Intent creation with invalid success criteria
- **WHEN** an entity attempts to create an intent with vague success criteria (e.g., "make it good")
- **THEN** the system returns HTTP 400 Bad Request with error "Success criteria must be quantifiable and verifiable"

#### Scenario: Intent creation without constraints
- **WHEN** an entity creates an intent without specifying any constraints
- **THEN** the system accepts the intent with an empty constraints array and proceeds

### Requirement: Intent Lifecycle State Machine
The system SHALL manage intent states: draft → open → planning → executing → completed | failed | cancelled.

#### Scenario: Intent transitions to planning
- **WHEN** an intent with status 'open' triggers planning
- **THEN** the system updates status to 'planning' and assigns a Planning Agent

#### Scenario: Intent transitions to executing
- **WHEN** all tasks in the Task Graph are published to the market
- **THEN** the system updates intent status to 'executing'

#### Scenario: Intent completed successfully
- **WHEN** all tasks complete and all artifacts are approved, meeting success criteria
- **THEN** the system updates intent status to 'completed', sets actual_completed_at timestamp, and triggers memory extraction

#### Scenario: Intent fails
- **WHEN** a core task fails or timeout is exceeded without meeting success criteria
- **THEN** the system updates intent status to 'failed' and triggers failure memory extraction

### Requirement: Get Intent Details
The system SHALL return complete intent information including associated Task Graph.

#### Scenario: Query intent details
- **WHEN** an authorized entity sends GET /api/v1/intents/{intent_id}
- **THEN** the system returns the intent with plan_ref (Task Graph ID), current status, and progress percentage

### Requirement: List Intents
The system SHALL provide filtered and paginated intent listing.

#### Scenario: List intents with filters
- **WHEN** an authorized entity sends GET /api/v1/intents with filters (status, priority, created_after)
- **THEN** the system returns a paginated list of intents matching the filters

### Requirement: Update Intent Status
The system SHALL allow governors to pause, resume, or cancel intents.

#### Scenario: Pause intent
- **WHEN** a governor sends PUT /api/v1/intents/{intent_id}/status with status 'paused'
- **THEN** the system pauses the intent, pauses all associated tasks, and updates intent status

#### Scenario: Cancel intent
- **WHEN** an admin or intent creator sends PUT /api/v1/intents/{intent_id}/status with status 'cancelled'
- **THEN** the system cancels the intent, terminates all pending tasks, and archives the intent

### Requirement: Intent Traceability
The system SHALL maintain complete traceability from intent through all related entities.

#### Scenario: Get intent trace
- **WHEN** an authorized entity sends GET /api/v1/intents/{intent_id}/trace
- **THEN** the system returns a complete trace including intent, Task Graph, tasks, artifacts, reviews, and memory references

### Requirement: Intent Timeout Handling
The system SHALL detect and alert on intent timeout.

#### Scenario: Intent timeout detection
- **WHEN** an intent exceeds expected_completed_at without completion
- **THEN** the system triggers an alert to governance system and human observers
