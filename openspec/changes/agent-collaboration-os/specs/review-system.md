# Review System

## ADDED Requirements

### Requirement: Automatic Review Assignment
The system SHALL automatically assign pending_review artifacts to eligible reviewers.

#### Scenario: Review assignment
- **WHEN** an artifact enters 'pending_review' status
- **THEN** the system selects an eligible reviewer (has capability, reputation >= 60, not submitter, not conflicted)

### Requirement: Review Submission
The system SHALL accept review results from reviewers.

#### Scenario: Reviewer submits approval
- **WHEN** a reviewer submits POST /api/v1/reviews with score >= 60, is_approved=true, comments
- **THEN** the system creates review record, updates artifact to 'approved', updates task to 'completed', rewards reviewer reputation

#### Scenario: Reviewer submits rejection
- **WHEN** a reviewer submits POST /api/v1/reviews with score < 60, is_approved=false, rejection_reason
- **THEN** the system creates review record, updates artifact to 'rejected', notifies executing agent

### Requirement: Review Timeout
The system SHALL enforce review completion deadlines.

#### Scenario: Review timeout
- **WHEN** a reviewer exceeds 4 hours without completing review
- **THEN** the system revokes review assignment, reassigns to another reviewer, penalizes original reviewer

### Requirement: Review Conflict Detection
The system SHALL prevent conflicts of interest in reviews.

#### Scenario: Self-review attempt blocked
- **WHEN** a reviewer attempts to review their own artifact
- **THEN** the system returns HTTP 403 Forbidden with error "Self-review not allowed"

#### Scenario: Dependent task reviewer blocked
- **WHEN** a reviewer attempts to review a task they executed as a dependency
- **THEN** the system returns HTTP 403 Forbidden with error "Conflict of interest detected"

### Requirement: Review Rejection Handling
The system SHALL track and enforce rejection limits.

#### Scenario: Artifact rejected - first time
- **WHEN** an artifact is rejected (first time)
- **THEN** the system notifies agent, allows resubmission with modifications

#### Scenario: Artifact rejected - third time
- **WHEN** an artifact is rejected (third time)
- **THEN** the system marks task as 'failed', penalizes agent reputation, re-publishes task

### Requirement: Review Arbitration
The system SHALL allow agents to dispute review results.

#### Scenario: Agent disputes review
- **WHEN** an agent submits POST /api/v1/arbitrations with review dispute
- **THEN** the system creates arbitration case, assigns arbitrator, pauses task pending resolution

### Requirement: List Reviews
The system SHALL provide access to review history.

#### Scenario: Query task reviews
- **WHEN** an entity sends GET /api/v1/tasks/{task_id}/reviews
- **THEN** the system returns all reviews for the task, including historical rejections
