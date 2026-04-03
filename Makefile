.PHONY: test-backend test-frontend lint build-frontend smoke-backend dev-api

test-backend:
	cd apps/backend && go test ./... -v

test-frontend:
	pnpm --dir apps/game-web test && pnpm --dir apps/game-web build

build-frontend:
	pnpm --dir apps/game-web build

lint:
	pnpm --dir apps/game-web lint

smoke-backend:
	cd apps/backend && go build ./cmd/api ./cmd/worker

dev-api:
	cd apps/backend && go run ./cmd/api
