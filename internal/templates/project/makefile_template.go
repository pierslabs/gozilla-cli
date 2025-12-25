package templates

func MakefileTemplate(data ProjectData) string {
	return `.PHONY: run build test docker-up docker-down migrate-up migrate-down

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	@echo "Migrations not configured yet"

migrate-down:
	@echo "Migrations not configured yet"

clean:
	rm -rf bin/
`
}
