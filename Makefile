# Constants

PROJECT_NAME = 'go-backend-template'
DB_URL = 'postgres://go-backend-template:go-backend-template@localhost:5454/go-backend-template?sslmode=disable'

ifeq ($(OS),Windows_NT) 
    DETECTED_OS := Windows
else
    DETECTED_OS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

# Help

.SILENT: help
help:
	@echo
	@echo "Usage: make [command]"
	@echo
	@echo "Commands:"
	@echo " rename-project name={name}    Rename project"	
	@echo	
	@echo " build-http                    Build http server"
	@echo " dev                           Watch mode"
	@echo
	@echo " migration-create name={name}  Create migration"
	@echo " migration-up                  Up migrations"
	@echo " migration-down                Down last migration"
	@echo
	@echo " docker-up                     Up docker services"
	@echo " docker-down                   Down docker services"
	@echo
	@echo " fmt                           Format source code"
	@echo " test                          Run unit tests"
	@echo

# Build

.SILENT: rename-project
rename-project:
    ifeq ($(name),)
		@echo 'new project name not set'
    else
        ifeq ($(DETECTED_OS),Darwin)
			@grep -RiIl '$(PROJECT_NAME)' | xargs sed -i '' 's/$(PROJECT_NAME)/$(name)/g'
        endif

        ifeq ($(DETECTED_OS),Linux)
			@grep -RiIl '$(PROJECT_NAME)' | xargs sed -i 's/$(PROJECT_NAME)/$(name)/g'
        endif

        ifeq ($(DETECTED_OS),Windows)
			@grep 'target is not implemented on Windows platform'
        endif
    endif

.SILENT: build-http
build:
	@go build -o ./bin/http-server ./src/main.go
	@echo executable file \"http-server\" saved in ./bin/http-server

# Test

.SILENT: test
test:
	@go test ./... -v

# Create migration

.SILENT: migration-create
migration-create:
	@migrate create -ext sql -dir ./migrations -seq $(name)

# Up migration

.SILENT: migration-up
migration-up:
	@migrate -database $(DB_URL) -path ./migrations up

# Down migration

.SILENT: "migration-down"
migration-down:
	@migrate -database $(DB_URL) -path ./migrations down 1

# Docker

.SILENT: docker-up
docker-up:
	@docker-compose up -d

.SILENT: docker-down
docker-down:
	@docker-compose down

# Format

.SILENT: fmt
fmt:
	@go fmt ./...

# Watch
.SILENT: dev
dev:
	ENT_PATH=.env gin run ./src/main.go --appPort 3000 --port 5000

# Start
.SILENT: start
start:
	ENT_PATH=.env go run ./src/main.go

# Default

.DEFAULT_GOAL := help
