# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Chinese translation repository for the Kubebuilder documentation. It's built using mdBook and contains comprehensive documentation for the Kubebuilder project, including tutorials, reference materials, and migration guides.

## Architecture

- **mdBook Documentation**: Built using Rust's mdBook with custom preprocessors
- **Custom Go Utilities**: Located in `book/utils/` for processing documentation
- **Preprocessors**:
  - `litgo.sh` - Processes Go code examples using literate-go utility
  - `markerdocs.sh` - Generates marker documentation
- **Content Structure**: Main content in `book/src/` with Chinese documentation
- **Build Artifacts**: Generated site outputs to `book/` directory

## Essential Commands

### Building and Serving

```bash
# Build the documentation (includes downloading dependencies)
./book/install-and-build.sh build

# Serve locally for development (after mdBook is installed)
cd book && mdbook serve
# Then visit http://localhost:3000

# Deploy to Cloudflare (requires environment variables)
./deploy-cloudflare.sh
```

### Development

```bash
# Build Go utilities (preprocessors)
cd book/utils
go build -o ../../bin/literate-go ./litgo
go build -o ../../bin/marker-docs ./markerdocs

# Manual mdBook operations (after install-and-build.sh)
cd book
mdbook build  # Build only
mdbook clean  # Clean build artifacts
```

## Key Dependencies

- **Go 1.23+**: Required for building utilities and examples
- **mdBook v0.4.40**: Rust-based documentation generator
- **controller-gen v0.19.0**: Kubernetes code generation tool
- **Node.js 18+**: For Cloudflare deployment via Wrangler

## Build Process

1. Downloads and installs mdBook binary for current platform
2. Builds Go utilities in `book/utils/`
3. Installs controller-gen tool
4. Runs mdBook with custom preprocessors
5. Outputs static site to `book/` directory

## Deployment

- **Cloudflare Pages**: Automated via GitHub Actions on push to main/master
- **Environment Variables Required**:
  - `CLOUDFLARE_API_TOKEN`
  - `CLOUDFLARE_ACCOUNT_ID`
- **Local Deploy**: Use `./deploy-cloudflare.sh` script

## File Structure

- `book/src/SUMMARY.md` - Table of contents for Chinese documentation
- `book/book.toml` - mdBook configuration
- `book/utils/` - Go utilities for documentation processing
- `book/install-and-build.sh` - Main build script
- `deploy-cloudflare.sh` - Cloudflare deployment script
- `.github/workflows/` - CI/CD automation

## Important Notes

- This is a translation project - content follows the original Kubebuilder documentation structure
- Custom preprocessors handle Go code examples and marker documentation
- Build script handles cross-platform mdBook installation automatically
- All scripts include proper error handling and platform detection
