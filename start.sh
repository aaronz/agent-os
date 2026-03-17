#!/bin/bash

set -e

echo "🚀 Starting Agent OS..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed${NC}"
    exit 1
fi

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$PROJECT_ROOT"

# Start backend
echo -e "${YELLOW}[1/3]${NC} Starting backend server..."
(
    cd "$PROJECT_ROOT"
    go run cmd/server/main.go &
    BACKEND_PID=$!
    echo $BACKEND_PID > /tmp/agent-os-backend.pid
    echo -e "${GREEN}Backend started on http://localhost:8080${NC}"
)

# Wait for backend to be ready
echo -e "${YELLOW}[2/3]${NC} Waiting for backend..."
for i in {1..30}; do
    if curl -s http://localhost:8080/api/v1/health &> /dev/null 2>&1 || curl -s http://localhost:8080/ &> /dev/null 2>&1; then
        echo -e "${GREEN}Backend is ready!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${YELLOW}Warning: Backend may not be ready yet${NC}"
    fi
    sleep 1
done

# Start frontend
echo -e "${YELLOW}[3/3]${NC} Starting frontend..."
(
    cd "$PROJECT_ROOT/dashboard"
    
    # Install dependencies if node_modules doesn't exist
    if [ ! -d "node_modules" ]; then
        echo "Installing frontend dependencies..."
        npm install
    fi
    
    npm run dev &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > /tmp/agent-os-frontend.pid
    echo -e "${GREEN}Frontend started on http://localhost:3000${NC}"
)

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Agent OS is running!${NC}"
echo -e "${GREEN}========================================${NC}"
echo "  Backend:  http://localhost:8080"
echo "  Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Save PIDs for cleanup
trap 'stop' SIGINT SIGTERM

stop() {
    echo ""
    echo -e "${YELLOW}Stopping Agent OS...${NC}"
    if [ -f /tmp/agent-os-backend.pid ]; then
        kill $(cat /tmp/agent-os-backend.pid) 2>/dev/null || true
        rm /tmp/agent-os-backend.pid
    fi
    if [ -f /tmp/agent-os-frontend.pid ]; then
        kill $(cat /tmp/agent-os-frontend.pid) 2>/dev/null || true
        rm /tmp/agent-os-frontend.pid
    fi
    echo -e "${GREEN}Stopped${NC}"
    exit 0
}

# Wait for any process to exit
wait
