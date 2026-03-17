# Memory System

## ADDED Requirements

### Requirement: Automatic Memory Extraction
The system SHALL automatically extract and store memories from system events.

#### Scenario: Task completion memory creation
- **WHEN** a task is marked 'completed'
- **THEN** the system automatically extracts execution details, solution approach, creates memory with type 'task'

#### Scenario: Task failure memory creation
- **WHEN** a task is marked 'failed'
- **THEN** the system automatically extracts failure reason, problem analysis, creates memory with type 'failure'

#### Scenario: Review completion memory creation
- **WHEN** a review is submitted
- **THEN** the system extracts review criteria, quality feedback, creates memory with type 'review'

#### Scenario: Intent completion memory creation
- **WHEN** an intent is marked 'completed'
- **THEN** the system extracts full project collaboration pattern, creates memory with type 'project'

### Requirement: Memory Search
The system SHALL support semantic search over memories.

#### Scenario: Semantic memory search
- **WHEN** an agent submits POST /api/v1/memory/search with query text
- **THEN** the system generates embedding, performs vector similarity search, returns top-K relevant memories

#### Scenario: Filtered memory search
- **WHEN** an agent submits POST /api/v1/memory/search with query and filters (type, intent_id, task_id)
- **THEN** the system combines vector search with entity filtering, returns relevant memories

### Requirement: Memory Retrieval by Entity
The system SHALL allow direct retrieval of memories related to specific entities.

#### Scenario: Get memories for intent
- **WHEN** an entity sends GET /api/v1/memory?intent_id={intent_id}
- **THEN** the system returns all memories linked to the specified intent

### Requirement: Memory Validity Management
The system SHALL track and update memory validity.

#### Scenario: Memory validity update
- **WHEN** a governor submits PUT /api/v1/memory/{memory_id} with validity status change
- **THEN** the system updates validity flag, logs the change

#### Scenario: Memory auto-invalidation
- **WHEN** a memory is not retrieved for 90 days
- **THEN** the system marks memory as 'invalid', removes from default search results

### Requirement: Memory Citation
The system SHALL track when agents use memories.

#### Scenario: Memory citation tracking
- **WHEN** an agent retrieves and uses a memory in execution
- **THEN** the system logs citation, increments retrieval_count for that memory

### Requirement: Memory Content Embedding
The system SHALL generate vector embeddings for all memory content.

#### Scenario: Embedding generation
- **WHEN** a new memory is created
- **THEN** the system automatically generates embedding using configured embedding model, stores in vector database
