# Project Structure

This repository serves as a template for building Golang applications using Hexagonal Architecture. It includes Docker support for easy containerized deployment. Follow the instructions below to get started.



# Table of Contents

1. [Project Structure](#project-structure)
   - [Layers and Folder Structure](#layers-and-folder-structure)
     - [API Layer](#api-layer)
     - [CMD Layer](#cmd-layer)
     - [Docs Layer](#docs-layer)
     - [Internal Layer](#internal-layer)
       - [Adapters Layer](#adapters-layer)
       - [Domain Layer](#domain-layer)
       - [Ports Layer](#ports-layer)
     - [Pkg Layer](#pkg-layer)
     - [Script Layer](#script-layer)
     - [Test Layer](#test-layer)
     - [Root Files](#root-files)
2. [Makefile](#makefile)
   - [Variables and Environment](#variables-and-environment)
   - [Targets](#targets)
     - [install](#install)
     - [buf](#buf)
     - [buf-win](#buf-win)
     - [run](#run)
     - [lint](#lint)
     - [test-run](#test-run)
     - [docker-build](#docker-build)
     - [docker-run](#docker-run)
     - [docker-compose-up](#docker-compose-up)
     - [docker-compose-down](#docker-compose-down)
     - [create-tree](#create-tree)
     - [fill-tree](#fill-tree)
     - [docker-network-up](#docker-network-up)
     - [docker-network-down](#docker-network-down)
   - [PHONY Targets](#phony-targets)
3. [Docker Compose File](#docker-compose-file)
   - [Version](#version)
   - [Services](#services)
     - [postgres](#postgres)
     - [app-service](#app-service)
   - [Networks](#networks)
   - [Volumes](#volumes)
   - [Environment Variables](#environment-variables)


## Layers and Folder Structure

### API Layer
The API layer contains the API definitions, including OpenAPI and Protocol Buffers (protobuf) specifications.

```
api
|--- openapi # OpenAPI specifications
|   |--- README.md # Documentation for OpenAPI specifications
|   |--- user.yaml # OpenAPI definition for User service
|--- proto # Protocol buffer definitions
|   |--- user
|   |   |--- v1
|   |   |   |--- user_service.proto # gRPC service definition for User service
|   |   |   |--- user_type.proto # gRPC message definitions for User service
|   |--- buf.gen.yaml # buf generate configuration
|   |--- buf.yaml # buf configuration
|   |--- README.md # Documentation for proto files
```

### CMD Layer
The `cmd` directory contains the entry point of the application. It includes the main application file and related documentation.

```
cmd
|--- main.go # Main application file
|--- README.md # Documentation for cmd directory
```

### Docs Layer
The `docs` directory contains the project documentation, including the directory tree structure.

```
docs
|--- README.md # Documentation for docs directory
|--- tree.md # Directory tree structure
```

### Internal Layer
The `internal` directory is where the core application logic resides. This includes the adapters, domain, and ports layers.

#### Adapters Layer
The adapters layer is divided into `driven` and `driver` adapters.

```
internal
|--- adapter # Adapters layer
|   |--- driven # Driven adapters (secondary)
|   |   |--- db # Database related code
|   |   |   |--- migration # Database migration files
|   |   |   |   |--- 000001_create_users_table.down.sql # SQL for down migration
|   |   |   |   |--- 000001_create_users_table.up.sql # SQL for up migration
|   |   |   |   |--- init.sql # SQL script for initial setup
|   |   |   |   |--- README.md # Documentation for migration files
|   |   |   |--- repository # Repositories for database interactions
|   |   |   |   |--- README.md # Documentation for repository
|   |   |   |   |--- user.go # User repository implementation
|   |   |   |   |--- user_mock.go # Mock implementation for user repository
|   |   |   |--- db_handler.go # Database handler
|   |   |   |--- postgres_transaction.go # Postgres transaction implementation
|   |   |   |--- postgres_transaction_mock.go # Mock implementation for postgres transaction
|   |   |--- passowrd.go # Password handling utilities
|   |--- driver # Driver adapters (primary)
|   |   |--- grpc # gRPC server code
|   |   |   |--- proto
|   |   |   |   |--- user
|   |   |   |   |   |--- v1
|   |   |   |   |   |   |--- user_service.pb.go # Generated gRPC code
|   |   |   |   |   |   |--- user_service_grpc.pb.go # Generated gRPC code
|   |   |   |   |   |   |--- user_type.pb.go # Generated gRPC code
|   |   |   |--- README.md # Documentation for gRPC server code
|   |   |   |--- user_service.go # gRPC service implementation
|   |   |--- http # HTTP server code
|   |   |   |--- README.md # Documentation for HTTP server code
|   |   |--- model # Models used in the application
|   |   |   |--- README.md # Documentation for models
|   |   |--- service # Application services
|   |   |   |--- README.md # Documentation for services
|   |   |   |--- user_service.go # User service implementation
|   |   |   |--- user_service_test.go # Tests for user service
|   |   |--- README.md # Documentation for adapters layer
```

#### Domain Layer
The domain layer contains the core business logic, including domain models and errors.

```
|--- domain # Domain layer (core business logic)
|   |--- error # Domain errors
|   |   |--- user.go # User-related errors
|   |--- model # Domain models
|   |   |--- user.go # User domain model
|   |--- README.md # Documentation for domain layer
```

#### Ports Layer
The ports layer defines the interfaces for the application, separating the core logic from the external components.

```
|--- port # Ports layer (interfaces)
|   |--- driven # Interfaces for driven adapters
|   |   |--- db # Database interfaces
|   |   |   |--- repository # Repository interfaces
|   |   |   |   |--- README.md # Documentation for repository interfaces
|   |   |   |   |--- user.go # User repository interface
|   |   |   |--- db_handler.go # Database handler interface
|   |   |   |--- db_transaction.go # Database transaction interface
|   |   |--- passowrd.go # Password interface
|   |--- driver # Interfaces for driver adapters
|   |   |--- model
|   |   |   |--- README.md # Documentation for model interfaces
|   |   |   |--- user.go # User model interface
|   |   |--- service
|   |   |   |--- README.md # Documentation for service interfaces
|   |   |   |--- user.go # User service interface
```

### Pkg Layer
The `pkg` directory contains package-level utilities and shared code.

```
|--- pkg # Package level utilities
|   |--- README.md # Documentation for package level utilities
```

### Script Layer
The `script` directory contains utility scripts for tasks like generating the directory tree.

```
|--- script # Utility scripts
|   |--- create_tree.ps1 # Script to create directory tree
|   |--- fill_tree.ps1 # Script to fill tree
```

### Test Layer
The `test` directory is for test-related files and documentation.

```
|--- test # Test related files
|   |--- README.md # Documentation for test files
```

### Root Files
The root of the repository contains various configuration files and the main README.

```
|--- .env # Environment variables
|--- .env.example # Example environment variables
|--- .gitignore # Git ignore file
|--- docker-compose.yml # Docker Compose configuration
|--- Dockerfile # Dockerfile for building the application
|--- go.mod # Go module file
|--- go.sum # Go dependencies
|--- Makefile # Makefile for task automation
|--- README.md # Project documentation
```

 Certainly! Below is an explanation of the Makefile, which you can add to your README:

---

## Makefile

The Makefile in this project automates various tasks to streamline the development and deployment process. Below is an explanation of each target in the Makefile.

### Variables and Environment

- **Loading Environment Variables**: If a `.env` file exists in the root directory, it will be included to load environment variables.
```makefile
ifneq (,$(wildcard ./.env))
    include .env
endif
```

### Targets

- **install**: Installs necessary Go modules and tools for protobuf and gRPC.
```makefile
install:
	@go mod tidy
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

- **buf**: Generates protobuf files using Buf. This target is for Unix-like systems.
```makefile
buf:
	@env PATH="$$PATH:$$(go env GOPATH)/bin" buf generate --template api/proto/buf.gen.yaml api/proto
	@echo "✅ buf done!"
```

- **buf-win**: Generates protobuf files using Buf for Windows systems.
```makefile
buf-win:
	@set PATH=%PATH%;%GOPATH%\bin
	@buf generate --template api\proto\buf.gen.yaml api/proto
	@echo "✅ buf done!"
```

- **run**: Runs the main Go application.
```makefile
run:
	go run ./cmd
```

- **lint**: Formats the Go code and runs the linter.
```makefile
lint:
	gofumpt -l -w .
	golangci-lint run -v
```

- **test-run**: Runs the Go tests.
```makefile
test-run:
	go test ./...
```

- **docker-build**: Builds a Docker image with the name specified by the `APP_NAME` environment variable.
```makefile
docker-build:
	docker build -t $(APP_NAME) .
```

- **docker-run**: Runs a Docker container from the built image and maps port 8080.
```makefile
docker-run:
	docker run -p 8080:8080 $(APP_NAME)
```

- **docker-compose-up**: Brings up the Docker Compose services and builds them if necessary.
```makefile
docker-compose-up:
	docker-compose up --build
```

- **docker-compose-down**: Brings down the Docker Compose services and removes volumes.
```makefile
docker-compose-down:
	docker-compose down --volumes
```

- **create-tree**: Creates a directory tree structure and saves it to `docs/tree.md`. This target is specific to Windows and uses PowerShell.
```makefile
create-tree: 
	powershell.exe -NoProfile -ExecutionPolicy Bypass -File ./script/create_tree.ps1 -path "." >> .\docs\tree.md
```

- **fill-tree**: Fills the tree structure for a specified path. This target is specific to Windows and uses PowerShell.
```makefile
fill-tree: 
	powershell.exe -NoProfile -ExecutionPolicy Bypass -File ./script/fill_tree.ps1 -path "../." 
```

- **docker-network-up**: Creates a Docker network with the name specified by the `APP_NETWORK_NAME` environment variable.
```makefile
docker-network-up:
	docker network create -d bridge $(APP_NETWORK_NAME)
```

- **docker-network-down**: Removes the Docker network with the name specified by the `APP_NETWORK_NAME` environment variable.
```makefile
docker-network-down:
	docker network rm $(APP_NETWORK_NAME)
```

### PHONY Targets
The `.PHONY` declaration specifies targets that are not associated with files.
```makefile
.PHONY: dev-run install buf lint run test-run docker-build docker-run docker-compose-up docker-compose-down create-tree fill-tree docker-network-up docker-network-down
```

This Makefile is designed to be cross-platform, supporting both Unix-like systems and Windows. It facilitates common development tasks such as code generation, linting, testing, Docker operations, and environment setup.

---

## Docker Compose File

This Docker Compose file defines the services, networks, and volumes required for the application. It uses environment variables to provide flexibility and reusability.

### Version

Specifies the version of the Docker Compose file format.
```yaml
version: '3.8'
```

### Services

#### postgres

- **image**: Uses the official PostgreSQL image.
- **container_name**: Names the container `postgres`.
- **environment**: Sets the `POSTGRES_PASSWORD` environment variable using the `DB_PASSWORD` environment variable from the `.env` file.
- **ports**: Maps the port specified by `DB_PORT` to `5432` on the container.
- **networks**: Connects to the network specified by the `APP_NETWORK_NAME` environment variable.
- **volumes**: 
  - `postgres-user-data` volume is used to persist PostgreSQL data.
  - Mounts the `init.sql` script to the container to initialize the database.

```yaml
  postgres:
    image: postgres
    container_name: postgres
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "${DB_PORT}:5432"
    networks:
      - ${APP_NETWORK_NAME}
    volumes:
      - postgres-user-data:/var/lib/postgresql/data
      - ./internal/adapter/driven/db/migration/init.sql:/docker-entrypoint-initdb.d/init.sql  # Mount init.sql into the container
```

#### app-service

- **image**: Uses the image specified by the `APP_IMAGE` environment variable.
- **container_name**: Names the container using the `CONTAINER_NAME` environment variable.
- **build**:
  - **dockerfile**: Specifies the Dockerfile to use for building the image.
  - **context**: Specifies the build context.
- **environment**: Sets environment variables for the application using values from the `.env` file.
- **ports**: Maps the port specified by `APP_PORT` to the same port on the container.
- **depends_on**: Ensures the `postgres` service is started before `app-service`.
- **restart**: Always restarts the container on failure.
- **networks**: Connects to the network specified by the `APP_NETWORK_NAME` environment variable.

```yaml
  app-service:
    image: ${APP_IMAGE}  # Specify the image name and tag
    container_name: ${CONTAINER_NAME} 
    build:
      dockerfile: Dockerfile
      context: .
    environment:
      DB_HOST: ${DB_HOST}
      DB_PORT: ${DB_PORT}
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      PORT: ${APP_PORT}
      IP: ${APP_IP}
    ports:
      - "${APP_PORT}:${APP_PORT}"
    depends_on:
      - postgres
    restart: always
    networks:
      - ${APP_NETWORK_NAME}
```

### Networks

Defines the network `APP_NETWORK`, using the `APP_NETWORK_NAME` environment variable for flexibility. It uses the `bridge` driver and is marked as external.

```yaml
networks:
  APP_NETWORK:
    driver: bridge
    external: true
```

### Volumes

Defines the `postgres-user-data` volume to persist PostgreSQL data.

```yaml
volumes:
  postgres-user-data:
```

### Environment Variables

To use this Docker Compose file effectively, ensure you have the following variables defined in your `.env` file:

```dotenv
DB_PASSWORD=your_db_password
DB_PORT=5432
APP_NETWORK_NAME=your_app_network_name
APP_IMAGE=your_app_image
CONTAINER_NAME=your_container_name
DB_HOST=your_db_host
DB_USER=your_db_user
DB_NAME=your_db_name
APP_PORT=8081
APP_IP=0.0.0.0
```

This setup allows you to configure your services flexibly using environment variables, making it easier to manage different environments (development, staging, production).

---