#!/bin/bash
#
# This script automates the building, packaging, and uploading of sshman binaries to a GitHub release.
#
# Usage:
# ./build-release.sh <tag_name>
#
# Example:
# ./build-release.sh v1.3.0
#
# Prerequisites:
# - Go (1.21+)
# - Node.js (18+) and npm
# - GitHub CLI (gh) installed and authenticated
# - The script must be run from the root of the repository

set -e

# --- Configuration ---
# Target platforms (GOOS GOARCH)
PLATFORMS=(
    "linux amd64"
    "linux arm64"
    "windows amd64"
    "darwin amd64"
    "darwin arm64"
    "freebsd amd64"
    "openbsd amd64"
)

# --- Pre-flight Checks ---
if [ -z "$1" ]; then
    echo "Error: No release tag specified."
    echo "Usage: $0 <tag_name>"
    exit 1
fi

TAG_NAME=$1
RELEASE_DIR="release-artifacts"

# Check for required commands
for cmd in go npm gh; do
    if ! command -v $cmd &> /dev/null; then
        echo "Error: $cmd is not installed. Please install it and try again."
        exit 1
    fi
done

# --- Build Steps ---
echo "--- Preparing for release: $TAG_NAME ---"

# 1. Build the frontend
echo "Building frontend assets..."
(cd frontend && npm install && npm run build)

# 2. Create a directory for release artifacts
echo "Creating release directory: $RELEASE_DIR"
rm -rf $RELEASE_DIR
mkdir -p $RELEASE_DIR

# 3. Build and package binaries for each platform
for platform in "${PLATFORMS[@]}"; do
    read -r GOOS GOARCH <<<"$platform"
    
    echo "Building for $GOOS/$GOARCH..."
    
    # Determine binary name and extension
    BINARY_NAME="sshman"
    if [ "$GOOS" = "windows" ]; then
        BINARY_NAME="sshman.exe"
    fi
    
    # Build the binary
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -ldflags="-s -w" -o "$RELEASE_DIR/$BINARY_NAME" .
    
    # Package the binary
    ASSET_NAME="sshman-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        ASSET_PATH="$RELEASE_DIR/${ASSET_NAME}.zip"
        (cd $RELEASE_DIR && zip "${ASSET_NAME}.zip" "$BINARY_NAME")
    else
        ASSET_PATH="$RELEASE_DIR/${ASSET_NAME}.tar.gz"
        (cd $RELEASE_DIR && tar -czvf "${ASSET_NAME}.tar.gz" "$BINARY_NAME")
    fi
    
    echo "  -> Packaged as $ASSET_PATH"
    
    # Clean up the binary after packaging
    rm "$RELEASE_DIR/$BINARY_NAME"
done

# --- Upload ---
echo "--- Uploading artifacts to release $TAG_NAME ---"
MAX_RETRIES=50
RETRY_DELAY=2
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    # Use --clobber to overwrite existing files if the script is re-run
    if gh release upload --clobber "$TAG_NAME" "$RELEASE_DIR"/*; then
        echo "Upload successful."
        break
    else
        RETRY_COUNT=$((RETRY_COUNT + 1))
        if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
            echo "Error: Upload failed after $MAX_RETRIES attempts."
            exit 1
        fi
        echo "Upload failed. Retrying in $RETRY_DELAY seconds... (Attempt $RETRY_COUNT of $MAX_RETRIES)"
        sleep $RETRY_DELAY
    fi
done

# --- Cleanup ---
echo "Cleaning up release directory..."
rm -rf $RELEASE_DIR

echo "--- Done! ---"
echo "All binaries have been built, packaged, and uploaded to release $TAG_NAME."

