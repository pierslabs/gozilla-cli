# Contributing to Gozilla

Thank you for your interest in contributing to Gozilla!

## Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gozilla.git
   cd gozilla
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Build the project:
   ```bash
   go build -o gozilla ./cmd/gozilla
   ```

## Project Structure

```
gozilla/
├── cmd/gozilla/          # CLI entry point
├── internal/
│   ├── commands/         # Cobra commands
│   │   ├── generate/     # Generate subcommands
│   │   ├── new.go        # New project command
│   │   └── root.go       # Root command
│   ├── generators/       # Code generators
│   │   ├── project.go    # Project generator
│   │   ├── module.go     # Module generator
│   │   └── container.go  # AST-based container updater
│   └── templates/        # Template directories (future)
└── pkg/                  # Shared packages
```

## Testing

Run tests:
```bash
go test ./...
```

Test the CLI manually:
```bash
# Build
go build -o gozilla ./cmd/gozilla

# Create a test project
./gozilla new test-api

# Generate a module
cd test-api
../gozilla generate module users
```

## Code Style

- Follow standard Go conventions
- Run `gofmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions small and focused

## Commit Messages

Use conventional commits format:
- `feat: add new feature`
- `fix: fix bug`
- `docs: update documentation`
- `refactor: refactor code`
- `test: add tests`
- `chore: update dependencies`

## Pull Requests

1. Create a feature branch from `main`
2. Make your changes
3. Add tests if applicable
4. Update documentation
5. Submit a pull request

## Questions?

Open an issue or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
