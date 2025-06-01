#!/bin/bash

# NameTidy Development Watch Script
# This script provides hot reload functionality for development

set -e

echo "🚀 Starting NameTidy development watch..."
echo "📁 Working directory: $(pwd)"
echo "🔍 Watching for Go file changes..."

# Build function
build_and_run() {
    echo "🔨 Building NameTidy..."
    if go build -o nametidy .; then
        echo "✅ Build successful!"
        echo "🧪 Running basic test..."
        ./nametidy --help > /dev/null && echo "✅ Binary works correctly!"
    else
        echo "❌ Build failed!"
        return 1
    fi
}

# Initial build
build_and_run

# Watch for changes
if command -v air > /dev/null 2>&1; then
    echo "🌪️  Using Air for hot reload..."
    air
elif command -v inotifywait > /dev/null 2>&1; then
    echo "👁️  Using inotifywait for file watching..."
    while true; do
        inotifywait -e modify,create,delete -r . \
            --include='\.go$' \
            --exclude='\.git|tmp|vendor' 2>/dev/null
        
        echo "📝 File change detected, rebuilding..."
        build_and_run
        echo "⏳ Waiting for changes..."
    done
else
    echo "⚠️  No file watcher available. Manual rebuild required."
    echo "💡 Run 'go build -o nametidy .' to rebuild after changes."
fi