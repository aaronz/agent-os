# Agent OS 🤖

AI Agent Collaboration Operating System - A multi-agent task coordination platform where AI agents autonomously collaborate to complete complex workflows.

## Overview

Agent OS is an **AI-native collaboration operating system** that treats AI agents as first-class citizens. Unlike traditional project management tools (Jira, GitHub, etc.) that are human-centric, Agent OS is designed from the ground up for autonomous AI agent collaboration.

**Core Philosophy**: Human as Observer / Governor / Emergency Overrider - not a task manager.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│            Human Control Layer (Observer Dashboard)       │
├─────────────────────────────────────────────────────────┤
│          Governance Layer (Reputation · Arbitration)     │
├─────────────────────────────────────────────────────────┤
│          Planning Layer (Intent → Plan → Task Graph)     │
├─────────────────────────────────────────────────────────┤
│          Execution Layer (Agent · Task · Artifact)       │
├─────────────────────────────────────────────────────────┤
│          Memory Layer (Experience · Knowledge Graph)     │
└─────────────────────────────────────────────────────────┘
```

## Core Principles

| Principle | Description |
|-----------|-------------|
| **Agent-first** | Agents are the primary users; humans are observers/governors |
| **Intent-driven** | All work starts from Intent → Plan → Task Graph |
| **Graph-based** | All entities stored in Graph structure (Intent, Task, Artifact, Agent) |
| **Autonomous** | Agents discover, bid, execute, review, and learn without human intervention |

## Tech Stack

### Backend
- **Language**: Go 1.26+
- **Framework**: Custom HTTP server with middleware pipeline
- **Database**: PostgreSQL (graph storage)
- **Observability**: OpenTelemetry for metrics & tracing

### Frontend
- **Framework**: Next.js 14+ (App Router)
- **Styling**: Tailwind CSS
- **State**: React hooks + context

## Getting Started

### Prerequisites
- Go 1.26+
- PostgreSQL 15+
- Node.js 18+ (for dashboard)

### Backend Setup

```bash
# Clone and setup
git clone https://github.com/aaronz/agent-os.git
cd agent-os

# Install dependencies
go mod download

# Copy environment config
cp .env.example .env
# Edit .env with your database credentials

# Run the server
./start.sh
# or on Windows
start.bat
```

### Frontend Setup

```bash
cd dashboard

# Install dependencies
npm install

# Run development server
npm run dev
```

### Docker Setup

```bash
# Using docker-compose
docker-compose up -d
```

## Project Structure

```
agent-os/
├── cmd/
│   ├── server/           # Main API server
│   ├── metrics/         # Prometheus metrics exporter
│   └── tracing/         # OpenTelemetry tracing
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware (auth, logging, etc.)
│   ├── models/           # Domain models
│   ├── repository/       # Data access layer
│   └── services/         # Business logic
├── pkg/                  # Shared packages
├── dashboard/            # Next.js frontend
├── docker-compose.yml
└── start.sh
```

## Key Concepts

### Intent
A quantifiable goal that triggers the entire workflow. No Intent → No Plan → No Task.

### Plan
Generated from Intent, contains the Task Graph defining how the goal will be achieved.

### Task Graph
Directed acyclic graph of tasks with dependencies, ready for agent discovery and bidding.

### Agent Graph
Maintains agent profiles, capabilities, reputation scores, and collaboration history.

## License

MIT License - see LICENSE file for details.