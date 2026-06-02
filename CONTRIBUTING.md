# Contributing to go-kiro-gateway

Thanks for your interest in contributing! This guide covers everything you need to get started.

## Ways to Contribute

- **Report bugs** — open a [Bug Report](https://github.com/chasedputnam/go-kiro-gateway/issues/new?template=bug_report.yml)
- **Request features** — open a [Feature Request](https://github.com/chasedputnam/go-kiro-gateway/issues/new?template=feature_request.yml)
- **Submit pull requests** — bug fixes, new features, documentation improvements
- **Improve docs** — fix typos, clarify setup steps, add examples

## Development Setup

### Prerequisites

- Go 1.25+
- Docker (optional, for container testing)
- `make` (for build targets)

### Getting Started

```bash
git clone https://github.com/chasedputnam/go-kiro-gateway.git
cd go-kiro-gateway

# Copy and configure environment
cp .env.example .env
# Edit .env with your Kiro credentials

# Download dependencies
cd gateway
go mod download

# Run tests
go test -race ./...

# Run linter
go vet ./...
```

### Project Structure

```
go-kiro-gateway/
├── gateway/            # All Go source code
│   ├── main.go         # Entry point, server setup
│   ├── handlers/       # HTTP route handlers (OpenAI, Anthropic endpoints)
│   ├── auth/           # Credential loading and token refresh logic
│   ├── proxy/          # Upstream request forwarding
│   ├── models/         # Model name normalization and registry
│   ├── middleware/      # Auth, logging, payload management middleware
│   └── Makefile        # Build, test, lint targets
├── .github/
│   ├── workflows/      # CI, release, and Docker GitHub Actions
│   ├── ISSUE_TEMPLATE/ # Bug report and feature request forms
│   └── PULL_REQUEST_TEMPLATE.md
├── docs/               # Additional documentation
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yml  # Local Docker Compose setup
└── .env.example        # Environment variable reference
```

### Makefile Targets

```bash
make build          # Build binary for local platform
make build-all      # Cross-compile for Linux, macOS, Windows
make test           # Run tests with race detection
make lint           # Run go vet
make clean          # Remove build artifacts
make docker         # Build Docker image locally
```

## Making Changes

1. **Fork** the repository and create a branch from `main`:
   ```bash
   git checkout -b fix/describe-your-change
   ```

2. **Make your changes.** Keep commits focused — one logical change per commit.

3. **Run tests and lint** before pushing:
   ```bash
   cd gateway
   go vet ./...
   go test -race ./...
   go mod tidy
   ```

4. **Open a pull request** against `main`. Fill out the PR template — describe what changed and how you tested it.

## Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Bug fix | `fix/short-description` | `fix/token-refresh-race` |
| Feature | `feat/short-description` | `feat/gemini-model-support` |
| Docs | `docs/short-description` | `docs/docker-compose-example` |
| Chore | `chore/short-description` | `chore/update-dependencies` |

## Code Conventions

- Follow standard Go formatting — run `gofmt` before committing
- Use `zerolog` for all logging (already imported); avoid `fmt.Println` in production paths
- Keep handler functions focused; extract helpers for reuse
- Add or update tests for any logic changes
- Avoid adding new dependencies unless necessary; discuss in an issue first

## Commit Messages

Use the [Conventional Commits](https://www.conventionalcommits.org/) style:

```
feat: add SOCKS5 proxy support
fix: handle empty refresh token gracefully
docs: clarify AWS SSO setup in README
chore: bump golang base image to 1.25
```

## Security Issues

Do **not** open a public issue for security vulnerabilities. See [SECURITY.md](SECURITY.md) for the private reporting process.

## License

By contributing, you agree that your contributions will be licensed under the [AGPL-3.0 License](LICENSE).
