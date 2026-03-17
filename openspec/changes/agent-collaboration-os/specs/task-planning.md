# Task Planning

## ADDED Requirements

### Requirement: Generate Task Graph
The system SHALL generate a directed acyclic graph (DAG) of tasks from an intent.

#### Scenario: Planning agent generates Task Graph
- **WHEN** a Planning Agent with valid capability submits POST /api/v1/planning/task-graph
- **THEN** the system validates the graph, stores it as published, and returns the Task Graph ID

#### Scenario: Task Graph contains circular dependency
- **WHEN** a Planning Agent submits a Task Graph with circular dependencies
- **THEN** the system returns HTTP 400 Bad Request with error "Circular dependency detected in Task Graph"

### Requirement: Task Graph Validation
The system SHALL validate Task Graph completeness and correctness.

#### Scenario: Task Graph missing success criteria coverage
- **WHEN** a Task Graph is submitted but doesn't cover all intent success_criteria
- **THEN** the system returns HTTP 400 Bad Request with error "Task Graph incomplete: missing coverage for {criteria}"

#### Scenario: Task Graph has invalid capability references
- **WHEN** a Task Graph references a capability that doesn't exist in the system
- **THEN** the system returns HTTP 400 Bad Request with error "Invalid capability: {capability_name}"

### Requirement: Task Graph Dependency Validation
The system SHALL ensure all task dependencies are valid.

#### Scenario: Task Graph has dangling dependency
- **WHEN** a task references a dependency that doesn't exist in the graph
- **THEN** the system returns HTTP 400 Bad Request with error "Dangling dependency: task {id} references non-existent task {dep_id}"

#### Scenario: Task Graph has self-dependency
- **WHEN** a task depends on itself
- **THEN** the system returns HTTP 400 Bad Request with error "Self-dependency detected in task {id}"

### Requirement: Progressive Task Publishing
The system SHALL publish tasks progressively as dependencies complete.

#### Scenario: Dependent task publishing
- **WHEN** a task's all dependencies are completed
- **THEN** the system automatically updates task status to 'open' and publishes to Task Market

### Requirement: Task Graph Versioning
The system SHALL support Task Graph versioning for modifications.

#### Scenario: Task Graph version update
- **WHEN** a Planning Agent submits PUT /api/v1/planning/task-graph/{graph_id} with modifications
- **THEN** the system creates a new version, validates, and publishes if valid

### Requirement: Get Task Graph Details
The system SHALL provide complete Task Graph information.

#### Scenario: Query Task Graph
- **WHEN** an authorized entity sends GET /api/v1/planning/task-graph/{graph_id}
- **THEN** the system returns the Task Graph with all tasks, dependencies, and status

### Requirement: Planning Agent Timeout
The system SHALL enforce planning timeout rules.

#### Scenario: Planning timeout
- **WHEN** a Planning Agent exceeds 4 hours without completing Task Graph
- **THEN** the system revokes planning permission, updates intent to 'open', and reassigns to another Planning Agent
