.PHONY: test-backend test-frontend lint dev-api

test-backend:
	cd apps/backend && go test ./... -v

test-frontend:
	pnpm --dir apps/game-web test && pnpm --dir apps/admin-web test

lint:
	pnpm --dir apps/game-web lint && pnpm --dir apps/admin-web lint

dev-api:
	cd apps/backend && go run ./cmd/api
