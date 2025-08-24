#!/bin/bash

APP_NAME="herp"
BINARY_NAME="herp"
SERVICE_NAME="herp"

echo "🚀 Deploying $APP_NAME..."

# Step 2: Pull latest code (optional, if using git)
# echo "📥 Pulling latest code..."
git pull origin main

# Step 3: Build the binary
echo "⚙️ Building the Go binary..."
go build -o "$BINARY_NAME" . || { echo "❌ Build failed"; exit 1; }

# Step 4: Restart systemd service
echo "🔁 Restarting $SERVICE_NAME service..."
sudo systemctl restart "$SERVICE_NAME"

# Step 5: Check service status
echo "📋 Status:"
sudo systemctl status "$SERVICE_NAME" --no-pager

echo "✅ Deployment complete."
