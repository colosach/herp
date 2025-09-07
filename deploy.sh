#!/bin/bash

APP_NAME="herp"
BINARY_NAME="herp"
SERVICE_NAME="herp"
FRONTEND_DIR="frontend"
PUBLIC_DIR="public"
BACKEND_DIR="backend"

echo "🚀 Deploying $APP_NAME..."

# Step 1: Pull latest main repo and submodules
echo "📥 Pulling latest code (backend + frontend submodule)..."
git pull origin main || { echo "❌ Git pull failed"; exit 1; }
git submodule update --init --remote --recursive || { echo "❌ Submodule update failed"; exit 1; }

# Step 2: Build frontend
echo "🌐 Building frontend..."
cd "$FRONTEND_DIR" || { echo "❌ Frontend directory not found"; exit 1; }
npm install
npm run generate || { echo "❌ Frontend build failed"; exit 1; }
cd ..

# Replace public folder with fresh build
rm -rf "$PUBLIC_DIR"
mkdir -p "$PUBLIC_DIR"
cp -r "$FRONTEND_DIR/.output/public/"* "$PUBLIC_DIR"/

# Step 3: Build backend and drop binary in root
echo "⚙️ Building backend..."
cd "$BACKEND_DIR" || { echo "❌ Backend directory not found"; exit 1; }
go build -o "../$BINARY_NAME" . || { echo "❌ Backend build failed"; exit 1; }
cd ..

# Step 4: Restart systemd service
echo "🔁 Restarting $SERVICE_NAME service..."
sudo systemctl restart "$SERVICE_NAME"

# Step 5: Check service status
echo "📋 Status:"
sudo systemctl status "$SERVICE_NAME" --no-pager

echo "✅ Deployment complete."
