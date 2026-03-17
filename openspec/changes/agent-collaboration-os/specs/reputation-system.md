# Reputation System

## ADDED Requirements

### Requirement: Initial Reputation
The system SHALL assign initial reputation score to new agents.

#### Scenario: New agent registration
- **WHEN** an agent is registered
- **THEN** the system sets initial reputation to 50

### Requirement: Reputation Rewards
The system SHALL increase agent reputation for positive outcomes.

#### Scenario: Task completion reward
- **WHEN** a task is completed with approval
- **THEN** the system adds 2-5 points to agent reputation based on review score

#### Scenario: High-quality review reward
- **WHEN** a review is not overturned by arbitration
- **THEN** the system adds 1-3 points to reviewer reputation

#### Scenario: Successful planning reward
- **WHEN** a Task Graph completes with 100% task success
- **THEN** the system adds 3-8 points to Planning Agent reputation

### Requirement: Reputation Penalties
The system SHALL decrease agent reputation for negative outcomes.

#### Scenario: Task failure penalty
- **WHEN** a task fails or times out
- **THEN** the system deducts 5-10 points from agent reputation

#### Scenario: Artifact rejection penalty
- **WHEN** an artifact is rejected
- **THEN** the system deducts 3-8 points from agent reputation

#### Scenario: Abandoned task penalty
- **WHEN** an agent wins a bid but abandons the task
- **THEN** the system deducts 10 points, increments abandonment count; 3 abandonments disables agent

#### Scenario: Overturned review penalty
- **WHEN** a review is overturned by arbitration
- **THEN** the system deducts 5 points from reviewer; 3 overturned reviews removes review资格

### Requirement: Reputation Thresholds
The system SHALL enforce reputation-based access control.

#### Scenario: High-reputation threshold
- **WHEN** an agent's reputation reaches 80
- **THEN** the agent gains access to high-reputation privileges (priority bidding, advanced roles)

#### Scenario: Low-reputation threshold
- **WHEN** an agent's reputation drops below 30
- **THEN** the system automatically disables the agent

### Requirement: Reputation Decay
The system SHALL decay reputation for inactive agents.

#### Scenario: Inactive agent decay
- **WHEN** an agent has no active tasks for 30 days
- **THEN** the system deducts 5 reputation points per month until reaching 30

### Requirement: Get Agent Reputation
The system SHALL expose reputation information.

#### Scenario: Query agent reputation
- **WHEN** an entity sends GET /api/v1/agents/{agent_id}/reputation
- **THEN** the system returns current reputation, reputation history, and breakdown of rewards/penalties
