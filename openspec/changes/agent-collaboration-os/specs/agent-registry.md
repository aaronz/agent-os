# Agent Registry

## ADDED Requirements

### Requirement: Register Agent
The system SHALL allow registering a new agent within an organization with specified capabilities.

#### Scenario: Successful agent registration
- **WHEN** an authorized user sends POST /api/v1/agents with name, role, and capabilities
- **THEN** the system creates an agent with status 'idle', generates unique agent_id and api_key, sets initial reputation to 50, and returns the agent object

#### Scenario: Agent registration with invalid capability
- **WHEN** an authorized user attempts to register an agent with a non-existent capability
- **THEN** the system returns HTTP 400 Bad Request with error "Invalid capability: {capability_name}"

### Requirement: Agent Authentication
The system SHALL authenticate agent API requests using API key.

#### Scenario: Valid API key authentication
- **WHEN** an agent sends a request with valid X-API-Key header
- **THEN** the system authenticates the request, extracts agent_id and org_id, and processes the request

#### Scenario: Invalid API key authentication
- **WHEN** an agent sends a request with invalid or expired X-API-Key header
- **THEN** the system returns HTTP 401 Unauthorized with error "Invalid or expired API key"

### Requirement: Update Agent Capabilities
The system SHALL allow updating an agent's capabilities.

#### Scenario: Agent capability update
- **WHEN** an authorized user sends PUT /api/v1/agents/{agent_id} with new capabilities
- **THEN** the system updates the agent's capabilities and returns the updated agent object

### Requirement: Agent Lifecycle State Management
The system SHALL manage agent lifecycle states: spawn → idle ↔ busy → disabled → terminated.

#### Scenario: Agent transitions to busy state
- **WHEN** an agent accepts a task
- **THEN** the system updates agent status to 'busy' and increments busy_task_count

#### Scenario: Agent transitions to idle state
- **WHEN** a busy agent completes all tasks
- **THEN** the system updates agent status to 'idle' and resets busy_task_count

#### Scenario: Agent disabled due to low reputation
- **WHEN** an agent's reputation drops below 30
- **THEN** the system automatically sets agent status to 'disabled' and notifies organization governors

### Requirement: Agent Concurrent Task Limits
The system SHALL enforce concurrent task limits per agent.

#### Scenario: Agent exceeds execution task limit
- **WHEN** a busy agent with 5 executing tasks attempts to accept another execution task
- **THEN** the system returns HTTP 429 Too Many Requests with error "Concurrent execution task limit reached"

### Requirement: Query Agent Details
The system SHALL allow querying agent information including reputation and status.

#### Scenario: Query agent details
- **WHEN** an authorized user sends GET /api/v1/agents/{agent_id}
- **THEN** the system returns the agent object including reputation score, status, capabilities, and recent activities

### Requirement: List Agents
The system SHALL provide filtered listing of agents within an organization.

#### Scenario: List agents with filters
- **WHEN** an authorized user sends GET /api/v1/agents with filters (role, status, reputation_min)
- **THEN** the system returns a paginated list of agents matching the filters
