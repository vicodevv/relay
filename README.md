# Relay

Event-driven workflow orchestration engine for coordinating multi-step business processes with automatic retries, state persistence, and rollback logic. Built with Go.

## Overview

Relay helps you build reliable workflows that survive failures. Think of it as the conductor that makes sure all parts of a complex operation happen in the right order, with the right data, even when things break.

### The Problem

Every company has complex processes:
- **Fintech:** KYC verification → payment processing → settlement
- **E-commerce:** Order → inventory check → payment → shipping → notifications  
- **Logistics:** Pickup → driver assignment → route optimization → delivery

**Without Relay, you're writing:**
- Retry logic everywhere
- Manual failure handling
- State tracking across services
- Race condition debugging

**With Relay, you get:**
- ✅ Durable execution (workflows survive crashes)
- ✅ Automatic retries (network failed? Try again)
- ✅ State persistence (remembers exactly where it stopped)
- ✅ Rollback on failure (auto-undo completed steps)
- ✅ Human approvals (pause for manual review)
- ✅ Event sourcing (complete audit trail)

## Quick Example

**Define a workflow in YAML:**

```yaml
name: loan_approval
steps:
  - id: check_credit
    service: credit_api
    endpoint: /check
    retry: 3
    timeout: 30s
    
  - id: verify_employment
    service: employment_api
    endpoint: /verify
    depends_on: check_credit
    
  - id: manual_review
    type: human_approval
    condition: credit_score < 650
    timeout: 24h
    
  - id: disburse_funds
    service: payment_api
    endpoint: /disburse
    depends_on: [verify_employment, manual_review]
    compensate: refund_payment
    
compensations:
  - id: refund_payment
    service: payment_api
    endpoint: /refund
```

**What Relay does automatically:**
- Retries network failures (3 attempts with backoff)
- Waits for dependencies (Step B waits for Step A)
- Pauses for human approval
- Rolls back on failure (refunds payment if shipping fails)
- Tracks every state change

## Tech Stack

- **Language:** Go 1.23
- **Web Framework:** Gin
- **Database:** PostgreSQL (workflow state)
- **Queue:** Redis (async execution)
- **Scheduling:** Cron (time-based triggers)

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 15+

### Installation

```bash
# Clone the repository
git clone https://github.com/vicodevv/relay.git
cd relay

# Install dependencies
go mod download

# Setup environment
cp .env.example .env

# Start PostgreSQL and Redis
docker-compose up -d

# Run the server
go run cmd/server/main.go
```

The API runs at `http://localhost:8080`

### Quick Test

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Create workflow definition
curl -X POST http://localhost:8080/api/v1/workflows/definitions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "loan_approval",
    "steps": [
      {
        "id": "check_credit",
        "service": "credit_api",
        "endpoint": "/check",
        "retry": 3
      }
    ]
  }'

# Start a workflow
curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "definition_name": "loan_approval",
    "input": {"loan_id": "L123", "amount": 50000}
  }'
```

## Project Structure
relay/
├── cmd/
│   ├── server/main.go          # HTTP server
│   ├── worker/main.go          # Workflow worker (coming soon)
│   └── cli/main.go             # CLI tool (coming soon)
├── internal/
│   ├── engine/                 # Workflow execution engine (coming soon)
│   ├── storage/                # PostgreSQL repository
│   ├── queue/                  # Redis queue (coming soon)
│   ├── http/                   # API handlers
│   └── scheduler/              # Cron triggers (coming soon)
├── pkg/
│   └── workflow/               # Workflow definitions
├── workflows/                  # Example YAML files
├── migrations/                 # SQL schemas
└── docker/

## Core Concepts

### Workflow Definition
A YAML file that describes the steps in a process.

### Workflow Instance
A running execution of a workflow definition.

### Step
A single unit of work (API call, database query, etc.).

### State Machine
Workflows transition through states: PENDING → RUNNING → COMPLETED (or FAILED).

### Durable Execution
After each step, state is saved to PostgreSQL. On server restart, incomplete workflows resume automatically.

### Event Sourcing
Instead of storing current state only, Relay stores ALL events (workflow_started, step_completed, etc.) for complete audit trail.

### Compensation
When a step fails, run compensations in REVERSE order to undo completed work.

### Human-in-the-Loop
Workflows can pause for manual approval, then resume when approved.

## API Endpoints

### Workflow Definitions
- `POST /api/v1/workflows/definitions` - Create workflow definition

### Workflow Instances
- `POST /api/v1/workflows` - Start a workflow
- `GET /api/v1/workflows/:id` - Get workflow status
- `GET /api/v1/workflows` - List all workflows
- `GET /api/v1/workflows/:id/events` - Get workflow event history

### Health
- `GET /api/v1/health` - Health check

## Roadmap

### Phase 1: HTTP API ✅
- [x] Workflow definitions
- [x] Workflow instances
- [x] Event sourcing
- [x] REST API

### Phase 2: Execution Engine (In Progress)
- [ ] Step executor
- [ ] Retry handler
- [ ] State machine
- [ ] HTTP client for service calls

### Phase 3: Advanced Features
- [ ] Human approvals
- [ ] Compensation/rollback
- [ ] Parallel execution
- [ ] Cron scheduling
- [ ] Visual dashboard

## Use Cases

- **Fintech:** Payment processing, KYC verification, loan approvals
- **E-commerce:** Order fulfillment, inventory management, refunds
- **Logistics:** Delivery coordination, driver assignment, route optimization
- **SaaS:** User onboarding, subscription management, billing

## Why Relay?

**vs Temporal:**
- Simpler (YAML vs complex SDK)
- Lighter (single binary vs cluster)
- African-friendly (works with flaky networks)

**vs Airflow:**
- Real-time (not just batch/scheduled)
- Built for APIs (not just data pipelines)
- Easier to deploy

## Development

```bash
# Run server
go run cmd/server/main.go

# Run tests (coming soon)
go test ./...

# Build binary
go build -o relay cmd/server/main.go
```

## Environment Variables

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5433
DB_USER=relay
DB_PASSWORD=relay123
DB_NAME=relay_db
REDIS_HOST=localhost
REDIS_PORT=6380
```

## Contributing

Contributions welcome! Open issues or PRs.

## License

MIT License
