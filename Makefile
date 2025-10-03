.PHONY: backend frontend dev dev-bg stop clean install test

# Start backend server
backend:
	go run .

# Start frontend development server  
frontend:
	cd frontend && npm run dev

# Install frontend dependencies
install:
	cd frontend && npm install

# Run both backend and frontend in development mode (foreground)
dev:
	@echo "Starting FlipAssistant development servers..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:5173"
	@echo ""
	@echo "Press Ctrl+C to stop both servers"
	@echo ""
	@trap 'make stop' INT; \
	go run . & \
	BACKEND_PID=$$!; \
	cd frontend && npm run dev & \
	FRONTEND_PID=$$!; \
	wait

# Run both servers in background (useful for testing)
dev-bg:
	@echo "Starting FlipAssistant servers in background..."
	@echo "Use 'make stop' to terminate them"
	go run . & echo $$! > .backend.pid
	cd frontend && npm run dev & echo $$! > ../.frontend.pid
	@echo "Backend PID: $$(cat .backend.pid)"
	@echo "Frontend PID: $$(cat .frontend.pid)"
	@echo "Servers started. Use 'make stop' to terminate."

# Stop all running servers
stop:
	@echo "Stopping FlipAssistant servers..."
	@pkill -f "go run" || true
	@pkill -f "vite" || true
	@if [ -f .backend.pid ]; then kill $$(cat .backend.pid) 2>/dev/null || true; rm -f .backend.pid; fi
	@if [ -f .frontend.pid ]; then kill $$(cat .frontend.pid) 2>/dev/null || true; rm -f .frontend.pid; fi
	@echo "All servers stopped."

# Check status of servers
status:
	@echo "Checking server status..."
	@if lsof -ti:8080 >/dev/null 2>&1; then echo "✅ Backend running on port 8080"; else echo "❌ Backend not running"; fi
	@if lsof -ti:5173 >/dev/null 2>&1; then echo "✅ Frontend running on port 5173"; else echo "❌ Frontend not running"; fi

# Build frontend for production
build:
	cd frontend && npm run build

# Clean build artifacts and stop servers
clean: stop
	cd frontend && rm -rf dist node_modules
	rm -f flips.db .backend.pid .frontend.pid

# Run tests (placeholder for future)
test:
	go test ./...
