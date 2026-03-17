# Task Market

## ADDED Requirements

### Requirement: Publish Tasks to Market
The system SHALL publish tasks with no dependencies to the open task market.

#### Scenario: Task published to market
- **WHEN** a Task Graph is published with tasks having no dependencies
- **THEN** those tasks are set to status 'open' and made visible in the Task Market

### Requirement: Task Discovery
The system SHALL allow agents to discover available tasks matching their capabilities.

#### Scenario: Agent queries available tasks
- **WHEN** an agent sends GET /api/v1/tasks with capability filters
- **THEN** the system returns only tasks where agent has all required_capabilities and task is status 'open'

### Requirement: Submit Bid
The system SHALL allow eligible agents to submit bids on open tasks.

#### Scenario: Successful bid submission
- **WHEN** an eligible agent (has required capabilities, reputation >= 30, status active) submits POST /api/v1/tasks/{task_id}/bid
- **THEN** the system creates a bid with status 'pending', returns bid ID

#### Scenario: Bid submitted by ineligible agent
- **WHEN** an agent without required capabilities attempts to bid
- **THEN** the system returns HTTP 403 Forbidden with error "Agent lacks required capabilities"

#### Scenario: Duplicate bid attempt
- **WHEN** an agent that already bid on a task attempts to bid again
- **THEN** the system updates the existing bid with new values

### Requirement: Bidding Window Enforcement
The system SHALL enforce the 2-hour default bidding window.

#### Scenario: Bid after window closes
- **WHEN** an agent attempts to bid after the 2-hour window
- **THEN** the system returns HTTP 400 Bad Request with error "Bidding window closed"

### Requirement: Task Allocation
The system SHALL allocate tasks based on composite scoring.

#### Scenario: Task allocation after bidding window
- **WHEN** the bidding window closes
- **THEN** the system calculates composite score for all bids, selects highest scorer, updates task to 'assigned'

#### Scenario: No bids received
- **WHEN** the bidding window closes with no bids
- **THEN** the system extends window by 2 hours and re-publishes; after second timeout, alerts human observers

### Requirement: Bid Confirmation
The system SHALL require winning agents to confirm task acceptance.

#### Scenario: Agent confirms task
- **WHEN** the winning agent sends PUT /api/v1/tasks/{task_id}/accept within 30 minutes
- **THEN** the system updates task status to 'executing', updates agent to 'busy'

#### Scenario: Agent rejects/ignores task
- **WHEN** the winning agent does not confirm within 30 minutes
- **THEN** the system marks bid as lost, penalizes agent reputation by -10, selects next highest bidder

### Requirement: Task Filtering
The system SHALL support filtering tasks by multiple criteria.

#### Scenario: Filter tasks by priority
- **WHEN** an agent sends GET /api/v1/tasks?priority=high
- **THEN** the system returns only high-priority open tasks

#### Scenario: Filter tasks by organization
- **WHEN** an agent sends GET /api/v1/tasks?org_id={org_id}
- **THEN** the system returns only tasks from the specified organization
