# Kubebuilder Documentation

This is the Chinese translation of the Kubebuilder project documentation, based on the [official documentation](https://github.com/kubernetes-sigs/kubebuilder).

## Translation Base

This translation is based on the following commit:
- **Source commit**: `3fd04199acb4556376dc17a7393bfd43bcd40c26`

## Local Development

### Quick Start

```bash
# Build and start local server
./book/install-and-build.sh build
cd book && mdbook serve

# Visit http://localhost:3000
```

### Manual Installation (Optional)

If you need to install dependencies manually:

1. Install [mdBook](https://rust-lang.github.io/mdBook/guide/installation.html)
2. Make sure [controller-gen](https://pkg.go.dev/sigs.k8s.io/controller-tools/cmd/controller-gen) is installed
3. Run `mdbook serve`

## Deployment

### Automatic Deployment

Pushing to `main` branch will automatically trigger GitHub Actions deployment to Cloudflare Pages.

### Manual Deployment

```bash
# Required environment variables:
# CLOUDFLARE_API_TOKEN
# CLOUDFLARE_ACCOUNT_ID

./deploy-cloudflare.sh
```
