# Human Control

## ADDED Requirements

### Requirement: Observer Role
The system SHALL provide read-only access to observers.

#### Scenario: Observer views data
- **WHEN** an observer sends GET requests to view data
- **THEN** the system returns data without modification capabilities

#### Scenario: Observer attempts modification
- **WHEN** an observer sends POST/PUT/DELETE requests
- **THEN** the system returns HTTP 403 Forbidden with error "Observers cannot modify data"

### Requirement: Governor Role
The system SHALL allow governors to configure and manage.

#### Scenario: Governor configures rules
- **WHEN** a governor updates governance rules
- **THEN** the system validates and applies the new rules

#### Scenario: Governor pauses intent
- **WHEN** a governor pauses an executing intent
- **THEN** the system pauses intent and all related tasks

### Requirement: Admin Role
The system SHALL provide full administrative control.

#### Scenario: Admin terminates agent
- **WHEN** an admin terminates an agent
- **THEN** the system terminates agent, revokes access, logs action

#### Scenario: Admin overrides system decision
- **WHEN** an admin overrides a system decision (with justification)
- **THEN** the system records override, logs justification, executes override

### Requirement: Emergency Intervention
The system SHALL allow emergency intervention by admins.

#### Scenario: Emergency intent cancellation
- **WHEN** an admin triggers emergency cancellation with reason
- **THEN** the system cancels intent immediately, logs intervention, notifies all affected agents

### Requirement: Audit Logging
The system SHALL log all human actions.

#### Scenario: Human action logged
- **WHEN** a human user performs any action
- **THEN** the system records action, user_id, timestamp, justification, affected entities

### Requirement: Dashboard Access
The system SHALL provide human-accessible dashboards.

#### Scenario: View organization dashboard
- **WHEN** an authorized human views the dashboard
- **THEN** the system displays: active intents, task progress, agent status summary, recent alerts

#### Scenario: View intent details
- **WHEN** a human views intent progress
- **THEN** the system displays: intent details, task graph visualization, artifact timeline, agent activities
