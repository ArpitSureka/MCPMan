#!/bin/bash

set -e

PROJECT_NAME="mcpman"
BUILD_DIR="bin"
FORMULA_FILE="$BUILD_DIR/mcpman.rb"
VERSION="0.1.0"

# Step 1: Build the Go binary
echo "🔨 Building the Go project..."
mkdir -p "$BUILD_DIR"
GOOS=darwin GOARCH=amd64 go build -o "$BUILD_DIR/$PROJECT_NAME" cmd/mcpman/main.go

# Step 2: Create Homebrew Formula
echo "📝 Creating Homebrew formula in $FORMULA_FILE..."

cat <<EOF > "$FORMULA_FILE"
class Mcpman < Formula
  desc "Your CLI tool description"
  homepage "https://github.com/ArpitSureka/mcpman"
  url "file://$(pwd)/$BUILD_DIR/$PROJECT_NAME"
  version "$VERSION"
  sha256 "$(shasum -a 256 $BUILD_DIR/$PROJECT_NAME | awk '{print $1}')"
  license "MIT"

  def install
    bin.install "$PROJECT_NAME"
  end
end
EOF

# Step 3: Install via brew
echo "🍺 Installing the CLI tool using Homebrew..."
brew install --formula "$FORMULA_FILE"

echo "✅ Installed $PROJECT_NAME successfully!"
