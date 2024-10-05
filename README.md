# High-Performance Distributed Rate Limiter

## Project Overview

This project implements a **High-Performance Distributed Rate Limiter** for an API gateway. The solution addresses the need to **limit the number of requests per second** per user in a **distributed environment**. The rate limiter ensures global rate limits are respected across multiple service instances, preventing individual nodes from over-serving users and ensuring fairness in handling traffic.

### Key Features
- **Concurrency**: Capable of handling multiple requests from different users concurrently using Go's concurrency primitives.
- **Distributed**: The rate limiter works across multiple instances of the service, ensuring global rate limits.
- **Efficiency**: Memory-efficient with minimized locking overhead to handle high traffic.
- **Customizable**: Each user can have a unique, dynamic rate limit.
- **Precision**: Implements the **sliding window algorithm** to maintain precision and fairness.
- **Persistence**: Redis is used as a globally distributed state to ensure rate limits are respected even after failures.

### Problem Statement

The goal is to build a rate-limiting service that satisfies the following requirements:
- **Concurrency**: Multiple requests from various users handled concurrently.
- **Distributed**: Rate limits must be respected globally across multiple instances of the service.
- **Efficiency**: Efficient memory and CPU utilization.
- **Customizable**: Support dynamic rate limits per user.
- **Precision**: Implement a sliding window algorithm.
- **Persistence**: Recover from failures while maintaining accurate request counts via Redis.

## Architecture

The application is split into different components to maintain separation of concerns:

1. **gRPC Service Layer**: This handles incoming gRPC requests and invokes the rate-limiting logic.
2. **Rate Limiting Logic**: Implements the core rate-limiting logic using the sliding window algorithm.
3. **Redis Integration**: Redis is used to maintain global consistency across service instances.
4. **Postgres Integration**: Rate limit configurations and user data are persisted in Postgres for long-term storage and recovery.
5. **Testing**: Unit tests and benchmarks are provided to ensure that the solution performs well under load and respects rate limits globally.

## gRPC Services

The project uses **gRPC** to expose rate-limiting functionality across distributed services.

### gRPC Files
The `api/proto/rate/v1/rate_service.proto` defines the following service and message types:

#### Service: `RateLimiterService`
- **`CheckRateLimit`**: This API checks if a request from a given user should be allowed based on the rate limit.
  - **Request**: `CheckRateLimitRequest`
    - Contains the user ID and rate limit to verify if the user is allowed to make the request.
  - **Response**: `CheckRateLimitResponse`
    - Returns a boolean indicating whether the request is allowed or denied based on the rate limit.
  
- **`GetUserRateLimit`**: This API retrieves the current rate limit for a specific user.
  - **Request**: `GetUserRateLimitRequest`
    - Contains the user ID.
  - **Response**: `GetUserRateLimitResponse`
    - Returns the current rate limit for the user.

- **`UpdateUserRateLimit`**: This API updates the rate limit for a specific user, allowing admins to dynamically adjust the rate limits.
  - **Request**: `UpdateUserRateLimitRequest`
    - Contains the user ID and the new rate limit to be applied.
  - **Response**: `UpdateUserRateLimitResponse`
    - Confirms the update of the user's rate limit.

These gRPC methods provide the core interface for interacting with the rate limiter from external systems, ensuring that rate limits are respected across distributed instances.

## Algorithms

### Sliding Window Algorithm

The **sliding window** algorithm ensures that rate limits are enforced precisely and fairly. The window "slides" over time, meaning that we account for recent requests and allow bursts, but not beyond the allowed rate. This ensures fairness, preventing a user from sending too many requests in a short period and gaming the system.

The rate-limiting logic is handled in the `RateLimitService`:
- Requests from the same user are counted.
- If the number of requests exceeds the configured rate limit within the defined time window, the request is rejected.
- The state of requests is stored in Redis to ensure persistence and global distribution across instances.

### Redis for Distributed Rate Limiting
Redis is utilized to store rate limit data for each user across distributed instances. This ensures that rate limits are respected globally even if multiple instances of the rate limiter are running.

## Setup Instructions

### Step 1: Environment Setup

Create an `.env` file in the project root with the following variables:

```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=rates
APP_PORT=8080
APP_IP=0.0.0.0
APP_IMAGE=app-service:latest
CONTAINER_NAME=app-container
APP_NAME=APP_NAME
APP_NETWORK_NAME=APP_NETWORK
REDIS_URL=redis:6379
WINDOW_MILI_SEC=100
PG_MIGRATION_FILES=file://internal/adapter/driven/db/migration
```

This file configures the database connection, application settings, and Redis connection details.

### Step 2: Running the Application

1. **Start the services with Docker**:

   ```bash
   docker-compose up --build
   ```

   This command will:
   - Set up a PostgreSQL database.
   - Spin up a Redis instance for distributed rate limiting.
   - Build and run the application service.

2. **Migrate Database Schema**:
   Ensure that your PostgreSQL database is correctly initialized with the necessary tables. The migrations are located in `internal/adapter/driven/db/migration`.

3. **Start the Application**:
   The application should be available at the configured port (e.g., `8080`).

### Step 3: Testing the Application

1. **API Testing**:
   The project includes OpenAPI definitions (`api/openapi/user.yaml`). These can be used to generate client libraries or test the API with tools like **Postman** or **BloomRPC**.

2. **Unit Tests and Benchmarks**:
   Run the tests using the following command:
   
   ```bash
   go test ./internal/adapter/driver/service/ -v
   ```

   This will execute all unit tests, including those that test the sliding window algorithm, Redis integration, and distributed consistency.

### Testing with gRPC Clients

To test the rate limiter using gRPC, you can use tools like **BloomRPC**. The `proto/rate/v1/rate_service.proto` defines the methods you can test.

1. **Start the application**:
   Ensure the application is running on the specified port (e.g., `8080`).

2. **Use BloomRPC**:
   - Import the `proto/rate/v1/rate_service.proto` file into **BloomRPC**.
   - Test methods like `CheckRateLimit` by providing the `userId` and `limit`.

## Test Coverage

We have implemented comprehensive tests, including:
- **Normal usage scenarios**: Testing rate limits for different users and configurations.
- **Edge cases**: Ensuring behavior when limits are exceeded or when no prior user data exists.
- **Distributed environment**: Ensuring that the rate limits are respected across multiple service instances using Redis.
- **Performance benchmarks**: Stress testing the application under high traffic to ensure the rate limiter is performant and efficient.

