#!/bin/bash

APP_NAME="herp"
BINARY_NAME="herp"
SERVICE_NAME="herp"
FRONTEND_DIR="frontend"
PUBLIC_DIR="public"
BACKEND_DIR="backend"

echo "🚀 starting $APP_NAME..."

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
make start || { echo "❌ Backend build failed"; exit 1; }

