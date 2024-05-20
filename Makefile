all: build/xdu-planet

.PHONY: web

build/xdu-planet: main.go web
	go mod tidy && go build -o build/xdu-planet

web:
	cd frontend && npm i && npm run build