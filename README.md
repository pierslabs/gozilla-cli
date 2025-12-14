<div align="center" >
  <div style="display: flex; align-items: center; justify-content: center; gap: 10px;">
  <h1>GoZilla</h1>
  <img src=".github/images/gozilla-logo.png" alt="GoZilla Logo" width="30" />

  </div>
  <p><strong>Generate production-ready Go projects with Clean Architecture in seconds</strong></p>

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

</div>

---

## 🚧 Work in Progress

A CLI tool for rapid Go development with automatic dependency injection.

## What it does

```bash
$ gozilla new my-api
✓ Clean Architecture structure
✓ Gin framework configured
✓ Docker & PostgreSQL setup
✓ Auto-wired DI container
✓ Example module with tests

$ gozilla g mod users
✓ Domain layer (entities, repositories)
✓ Application layer (use cases, DTOs)
✓ Infrastructure layer (handlers, DB)
✓ Auto-wired in DI container
✓ Full CRUD + tests generated

$ gozilla g mod orders --depends=users
✓ Module with dependencies resolved
✓ Ready to code immediately
```

## Why?

**The problem:** Setting up Go projects with proper architecture takes hours.
Copy-pasting boilerplate. Manually wiring dependencies. Inconsistent structure
across projects.

**The solution:** Generate everything automatically. Focus on business logic,
not infrastructure.

## Key Features

- ⚡ **Instant setup** - From zero to running API in 30 seconds
- 🏗️ **Clean Architecture** - Domain-driven design with clear boundaries
- 🔌 **Auto DI wiring** - Dependencies injected automatically
- 🧩 **Modular** - Each feature is a self-contained module
- 🧪 **Test ready** - Unit tests generated with every module
- 📦 **Framework agnostic** - Gin, Fiber, Echo, Chi support (coming soon)

## Generated Structure

```
my-api/
├── cmd/api/main.go
├── internal/
│   ├── domain/                    # Shared domain
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── http/
│   │   └── container/             # DI container ✨
│   └── modules/
│       └── users/
│           ├── user.module.go     # Module DI
│           ├── domain/
│           ├── application/
│           └── infra/
├── docker-compose.yaml
└── Makefile
```

### Can I use a different database?

Yes! PostgreSQL comes pre-configured for convenience, but you can:

1. Replace the repository implementation in `infra/repository.go`
2. Update `infrastructure/database/` with your DB of choice
3. Modify `docker-compose.yaml`

The repository pattern makes switching databases straightforward.

Future versions will support `--db` flag during project creation.

## Join Early Access

👉 [Get notified when it launches](https://gozilla.vercel.app)

## Installation

### From Source

```bash
git clone https://github.com/pierslabs/gozilla.git
cd gozilla
go build -o gozilla ./cmd/gozilla
sudo mv gozilla /usr/local/bin/
```

### Verify Installation

```bash
gozilla --version
```

## Quick Start

### 1. Create a new project

```bash
gozilla new my-api
cd my-api
```

### 2. Start the database

```bash
docker-compose up -d
```

### 3. Run the application

```bash
make run
```

Your API is now running at `http://localhost:8080`

Test the health endpoint:

```bash
curl http://localhost:8080/api/v1/health
```

### 4. Generate your first module

```bash
gozilla generate module users
```

This creates a complete CRUD module with:

- Domain entities and repository interface
- Use cases (Create, Get, List, Update, Delete)
- HTTP handlers and routes
- Auto-wired in the DI container

## Development

### Build

```bash
go build -o gozilla ./cmd/gozilla
```

### Test

```bash
go test ./...
```

## Roadmap

- [x] Core CLI commands (`new`, `generate module`)
- [x] Gin framework support
- [x] Auto dependency wiring
- [ ] Test generation for modules
- [ ] Migration generation
- [ ] Multi-framework support (Fiber, Echo, Chi)
- [ ] Custom templates
- [ ] GitHub Actions workflows

## Status

MVP completed! **Star ⭐ this repo** if you find this useful!

## Contributing

Not accepting contributions yet. Follow for updates.

## License

MIT
