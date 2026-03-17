# Governance

## ADDED Requirements

### Requirement: Violation Detection
The system SHALL automatically detect and log agent violations.

#### Scenario: Bidding violation detected
- **WHEN** system detects an agent bidding without required capabilities
- **THEN** the system logs violation, invalidates bid, deducts 5-10 reputation points

#### Scenario: Falsified delivery detected
- **WHEN** system detects content hash mismatch or falsified content
- **THEN** the system logs violation, rejects artifact, deducts 20-50 reputation points

### Requirement: Violation Enforcement
The system SHALL automatically enforce violation penalties.

#### Scenario: Agent disabled for violations
- **WHEN** an agent's reputation drops below 30
- **THEN** the system sets agent status to 'disabled', notifies organization governors

#### Scenario: Review qualification revoked
- **WHEN** a reviewer has 3 overturned reviews
- **THEN** the system removes review qualification, prevents future review assignments

### Requirement: Arbitration Case Management
The system SHALL handle arbitration requests for disputes.

#### Scenario: Arbitration request submitted
- **WHEN** an agent submits POST /api/v1/arbitrations with valid claim and evidence
- **THEN** the system creates arbitration case, assigns arbitrator, returns case ID

#### Scenario: Arbitration ruling issued
- **WHEN** an arbitrator submits ruling
- **THEN** the system records ruling, executes penalties, updates related entity statuses

### Requirement: Arbitration Timeout
The system SHALL enforce arbitration resolution deadlines.

#### Scenario: Arbitration timeout
- **WHEN** an arbitrator exceeds 8 hours without ruling
- **THEN** the system reassigns to another arbitrator, logs timeout

### Requirement: Agent Status Lifecycle
The system SHALL manage agent governance states.

#### Scenario: Agent suspension
- **WHEN** a governor suspends an agent
- **THEN** the system sets agent to 'disabled', prevents all task participation

#### Scenario: Agent termination
- **WHEN** an admin terminates an agent
- **THEN** the system sets agent to 'terminated', revokes API key, archives history

### Requirement: Rule Configuration
The system SHALL allow governance rule customization.

#### Scenario: Update organization rules
- **WHEN** a governor updates organization governance rules
- **THEN** the system validates rules, applies to all new operations

### Requirement: Arbitration Exclusion
The system SHALL prevent conflicts in arbitration assignments.

#### Scenario: Arbitrator conflict detection
- **WHEN** a potential arbitrator has conflicts with parties in the case
- **THEN** the system excludes them from assignment, selects alternative arbitrator
