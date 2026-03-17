# Design: Web Dashboard for Agent Collaboration OS

## 1. Technology Stack

- **Frontend Framework**: Next.js 14 (React)
- **UI Library**: Tailwind CSS + shadcn/ui
- **State Management**: React Query + Zustand
- **Real-time**: WebSocket client
- **Charts**: Recharts
- **Forms**: React Hook Form + Zod

## 2. Application Structure

```
/                       # Redirects to /dashboard
/login                  # Authentication
/dashboard              # Main dashboard (default view)
/intents               # Intent management
  /[id]                # Intent detail
  /new                 # Create new Intent
/tasks                 # Task monitoring
/agents                # Agent registry
/artifacts             # Artifact browser
/memory                # Memory search
/governance            # Rule configuration
/reputation            # Reputation management
/arbitration           # Arbitration cases
/alerts                # Alert center
/settings              # Organization settings
```

## 3. Key Components

### Navigation
- Sidebar navigation with role-based menu items
- Breadcrumb trail for deep navigation
- Global search (Cmd+K)

### Dashboard View
- System health summary cards
- Active intents count
- Tasks in progress / completed
- Agent status distribution (pie chart)
- Recent activity feed
- Alert notifications

### Intent Management
- Intent list with filters (status, priority, date)
- Intent creation form with validation
- Intent detail view with task graph visualization
- Timeline of status changes

### Agent Registry
- Agent cards with reputation score
- Capability badges
- Activity history
- Status toggle (enable/disable)

### Task Monitor
- Kanban-style task board (optional)
- Task detail modal
- Progress timeline
- Artifact preview

### Governance Panel
- Rule configuration forms
- Reputation score charts
- Arbitration case queue
- Audit log viewer

## 4. API Integration

All endpoints reuse existing REST APIs:
- `GET /api/v1/organizations` → Organization list
- `GET /api/v1/intents` → Intent list
- `GET /api/v1/tasks` → Task list
- `GET /api/v1/agents` → Agent list
- `GET /api/v1/memory/search` → Memory search
- `GET /api/v1/metrics` → System metrics
- `GET /api/v1/audit` → Audit logs

## 5. Authentication

- JWT-based authentication
- Login stores token in httpOnly cookie
- Role-based route protection
- Session timeout (24 hours)

## 6. Real-time Updates

- WebSocket connection to backend event stream
- Event types: task.completed, task.failed, intent.status_changed, alert.created
- Optimistic UI updates

## 7. UI/UX Guidelines

- Dark mode support
- Consistent color palette (blue primary, green success, red error)
- Loading skeletons for async content
- Toast notifications for actions
- Confirmation dialogs for destructive actions
