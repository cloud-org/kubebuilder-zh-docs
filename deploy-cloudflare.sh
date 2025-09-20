#!/bin/bash

#  Copyright 2020 The Kubernetes Authors.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

set -e

echo "🚀 Starting Cloudflare deployment process..."

# Get the directory that this script file is in
THIS_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
cd "$THIS_DIR"

# Build the project first using the existing install-and-build script
echo "📦 Building the project..."
chmod +x ./book/install-and-build.sh
./book/install-and-build.sh build

# Check if book/book directory exists (mdBook output)
if [ ! -d "book/book" ]; then
    echo "❌ Error: book/book directory not found. Build may have failed."
    exit 1
fi

# Install Wrangler if not available
if ! command -v wrangler &> /dev/null; then
    echo "📥 Installing Wrangler CLI..."
    npm install -g wrangler
fi

# Deploy to Cloudflare Pages
echo "🌐 Deploying to Cloudflare Pages..."

# Check if wrangler.toml exists, if not create a basic one
if [ ! -f "wrangler.toml" ]; then
    echo "📝 Creating wrangler.toml configuration..."
    cat > wrangler.toml << EOF
name = "kubebuilder-zh-docs"
compatibility_date = "2024-01-01"

[env.production]
name = "kubebuilder-zh-docs"

[[pages_builds]]
destination = "book/book"
build_command = "./book/install-and-build.sh build"
build_watch_dirs = ["book/src", "book"]
EOF
fi

# Deploy using Wrangler
wrangler pages deploy book/book --project-name=kubebuilder-zh-docs

echo "✅ Deployment completed successfully!"
echo "🔗 Your site should be available at: https://kubebuilder-zh-docs.pages.dev"