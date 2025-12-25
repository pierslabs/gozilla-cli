package templates

import (
	"fmt"
	"strings"
)

func ReadmeTemplate(data ProjectData) string {
	projectName := strings.Title(data.ProjectName)
	return fmt.Sprintf(`# %s

Generated with [Gozilla](https://github.com/pierslabs/gozilla)

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make

### Installation

1. Start the database:
`+"```bash"+`
make docker-up
`+"```"+`

2. Run the application:
`+"```bash"+`
make run
`+"```"+`

The API will be available at http://localhost:8080

### API Endpoints

- `+"`GET /api/v1/health`"+` - Health check

## Development

### Generate a new module

`+"```bash"+`
gozilla generate module users
`+"```"+`

### Run tests

`+"```bash"+`
make test
`+"```"+`

### Build

`+"```bash"+`
make build
`+"```"+`

## Project Structure

`+"```"+`
.
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                # Shared domain
│   ├── infrastructure/        # Infrastructure layer
│   │   ├── config/           # Configuration
│   │   ├── database/         # Database connection
│   │   ├── http/             # HTTP server
│   │   └── container/        # DI container
│   └── modules/              # Feature modules
│       └── health/           # Health check module
├── migrations/               # Database migrations
└── docker-compose.yaml       # Docker setup
`+"```"+`

## License

MIT
`, projectName)
}
