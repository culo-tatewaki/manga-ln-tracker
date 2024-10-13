#!/bin/bash

# Remove previous build
rm -rf dist && rm -rf binaries

# Build server for Windows
GOOS=windows GOARCH=amd64 go build -o binaries/server-windows.exe

# Build server for Linux
GOOS=linux GOARCH=amd64 go build -o binaries/server-linux

# Build server for macOS
GOOS=darwin GOARCH=amd64 go build -o binaries/server-macos

# Build Electron App
npx electron-builder
