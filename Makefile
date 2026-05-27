.PHONY: all frontend backend clean dev

all: frontend backend

frontend:
	cd frontend && npm install && npm run build-only

backend: frontend
	go mod tidy
	go build -o stronglifts .

clean:
	rm -rf frontend/dist frontend/node_modules stronglifts stronglifts.db

# Dev mode: run Go server (frontend must be built first)
run: all
	./stronglifts

# Dev frontend with hot reload (proxy to Go backend)
dev-frontend:
	cd frontend && npm run dev

# Dev backend only
dev-backend:
	go run .
