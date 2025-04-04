# Makefile at root of musicbook-hub

.PHONY: build backend frontend dev

# Build React frontend and Go binary
build: frontend backend

# Build frontend (React with Tailwind)
frontend:
	cd web && npm install && npm run build

# Cross-compile Go backend for Linux (x86_64)
backend:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/musicbook-hub

# Dev mode: runs both frontend and backend locally
dev:
	@echo "Starting Go backend and React frontend..."
	@make -j2 run-backend run-frontend

run-backend:
	cd . && go run main.go

run-frontend:
	cd web && npm run dev

#scp musicbook-hub youruser@yourserver:/home/youruser/musicbook-hub/
#rsync -av web/dist/ youruser@yourserver:/home/youruser/musicbook-hub/web/dist/

