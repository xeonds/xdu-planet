all: build/xdu-planet

.PHONY: frontend

build/xdu-planet: main.go frontend
	go mod tidy && go build -o build/xdu-planet

frontend:
	cd frontend && pnpm i && pnpm run build
