#!/bin/bash

set -e

# Configuration
BASE_URL="${SITE_URL:-https://kubebuilder-zh.vibecodinghub.org/}"
BOOK_DIR="$(dirname "$0")/book"
SITEMAP_TOOL="$(dirname "$0")/utils/sitemap/main.go"

# Build the sitemap generator
echo "Building sitemap generator..."
cd "$(dirname "$0")/utils/sitemap"
go build -o ../../bin/sitemap-gen main.go
cd - > /dev/null

# Generate sitemap.xml
echo "Generating sitemap.xml..."
./bin/sitemap-gen "$BOOK_DIR" "$BASE_URL" > "$BOOK_DIR/sitemap.xml"

echo "Sitemap generated at $BOOK_DIR/sitemap.xml"