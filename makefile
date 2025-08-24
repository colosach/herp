APP_NAME=herp
FRONTEND_DIR=frontend
BACKEND_DIR=backend
PUBLIC_DIR=public

.PHONY: build frontend backend run serve-backend clean dev-frontend

# Default build: frontend + backend
build: frontend

# Build Nuxt (SSR off, static output)
frontend:
	cd $(FRONTEND_DIR) && git pull origin main && npm install && npm run generate
	rm -rf $(PUBLIC_DIR)
	mkdir -p $(PUBLIC_DIR)
	cp -r $(FRONTEND_DIR)/.output/public/* $(PUBLIC_DIR)/

# Build Go backend (binary in project root)
backend:
	cd $(BACKEND_DIR) && go build -o ../$(APP_NAME
	
# Just run backend (useful for API dev while Nuxt runs with npm run dev)
serve-backend:
	cd $(BACKEND_DIR) && go run main.go

# Run after build
run: build serve-backend

# Clean up build artifacts
clean:
	rm -rf $(PUBLIC_DIR) $(APP_NAME)

# Frontend dev mode (hot reload on :3000)
dev-frontend:
	cd $(FRONTEND_DIR) && npm run dev
