#!/bin/bash

# Build the plugin
echo "Building Obsidian Novel Maker plugin..."

# Change to plugin directory
cd "$(dirname "$0")"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install --silent
fi

# Build with esbuild
echo "Bundling plugin..."
node build.js

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    exit 0
else
    echo "❌ Build failed!"
    exit 1
fi
