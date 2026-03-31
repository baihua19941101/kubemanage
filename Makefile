SHELL := /bin/bash

.PHONY: backend-run backend-test frontend-dev frontend-build

backend-run:
	cd backend && go run ./cmd/server

backend-test:
	cd backend && go test ./...

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build
