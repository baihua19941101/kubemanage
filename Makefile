SHELL := /bin/bash

.PHONY: backend-run backend-test frontend-dev frontend-build stage2-qa

backend-run:
	cd backend && go run ./cmd/server

backend-test:
	cd backend && go test ./...

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

stage2-qa:
	./scripts/p207_stage2_qa.sh
