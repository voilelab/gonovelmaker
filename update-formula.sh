#!/bin/bash

set -e

# Check if version argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v0.0.7"
    exit 1
fi

VERSION="$1"
FORMULA_FILE="Formula/novelmaker-obs.rb"

# Remove 'v' prefix if present for consistency
VERSION_NUM="${VERSION#v}"
VERSION_TAG="v${VERSION_NUM}"

# GitHub URL for the tarball
URL="https://github.com/voilelab/gonovelmaker/archive/refs/tags/${VERSION_TAG}.tar.gz"

echo "Downloading tarball from: $URL"

# Download the tarball to a temporary location
TEMP_FILE=$(mktemp)
if ! curl -L -o "$TEMP_FILE" "$URL"; then
    echo "Error: Failed to download tarball from $URL"
    rm -f "$TEMP_FILE"
    exit 1
fi

echo "Calculating SHA256..."
# Calculate SHA256
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    SHA256=$(shasum -a 256 "$TEMP_FILE" | awk '{print $1}')
else
    # Linux
    SHA256=$(sha256sum "$TEMP_FILE" | awk '{print $1}')
fi

# Clean up temporary file
rm -f "$TEMP_FILE"

echo "Version: $VERSION_TAG"
echo "SHA256: $SHA256"

# Use sed to update the version and SHA256
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS sed requires different syntax
    sed -i '' "s|archive/refs/tags/v[0-9.]*\.tar\.gz|archive/refs/tags/${VERSION_TAG}.tar.gz|" "$FORMULA_FILE"
    sed -i '' "s/sha256 \"[a-f0-9]*\"/sha256 \"$SHA256\"/" "$FORMULA_FILE"
else
    # Linux sed
    sed -i "s|archive/refs/tags/v[0-9.]*\.tar\.gz|archive/refs/tags/${VERSION_TAG}.tar.gz|" "$FORMULA_FILE"
    sed -i "s/sha256 \"[a-f0-9]*\"/sha256 \"$SHA256\"/" "$FORMULA_FILE"
fi

echo "✓ Formula updated successfully!"
echo ""
echo "Changes made:"
echo "  URL: $URL"
echo "  SHA256: $SHA256"
echo ""
echo "To verify the changes, run:"
echo "  brew install --build-from-source $FORMULA_FILE"
echo "  brew test novelmaker-obs"
