## Design and Architecture Overview

This project follows the **Hexagonal Architecture (Ports and Adapters)** design pattern, which separates the business logic from external concerns like databases, message queues, APIs, etc. This architecture provides clear boundaries between the core domain logic and the infrastructure, promoting a modular and flexible design. 

The project uses **Redis** and **PostgreSQL** for data persistence, and **gRPC** for service-to-service communication. The core focus is on providing a scalable and distributed rate-limiting service with high efficiency and precision.

### Hexagonal Architecture (Ports and Adapters)

The **Hexagonal Architecture**, also known as **Ports and Adapters Architecture**, emphasizes separation between the core business logic (domain) and the outer layers (infrastructure) to achieve flexibility, maintainability, and testability.

- **Ports**: Interfaces that define how the outside world interacts with the core application logic. Ports define the use cases of the application.
- **Adapters**: Implementations of the ports, providing the actual mechanisms for interacting with databases, external APIs, and other systems.
- **Core Domain**: Contains the business rules and logic, independent of any infrastructure or external frameworks.

This project is structured as follows:

### Layers:

1. **Domain (Core)**: The core logic of the application. This layer contains the **business rules** and **use cases**. It does not depend on any external frameworks, databases, or APIs.
   
   - **Domain Models**: Located in `internal/domain/model/`, where key entities like `UserRateLimit` are defined. These models represent the core business concepts.
   - **Domain Errors**: Located in `internal/domain/error/`, where business-related errors are managed.
   
2. **Ports**: These are interfaces that define how the core application interacts with the outside world. For example, how the rate limiter logic interacts with the database or cache.

   - **Driver Ports**: Define how external services interact with the core application logic (e.g., gRPC or HTTP services). 
     - Located in `internal/port/driver/service/`.
   - **Driven Ports**: Define the interfaces that the application uses to interact with external systems, such as the database or cache. 
     - Located in `internal/port/driven/db/` for database interactions and `internal/port/driven/cache.go` for cache interactions.

3. **Adapters (Infrastructure)**: These implement the interfaces defined in the Ports layer. They include:

   - **Driven Adapters**: These adapters implement interactions with external systems such as Redis and PostgreSQL.
     - **Database Adapters**: Located in `internal/adapter/driven/db/`, where actual database-related code is present (e.g., `PostgresTransaction`, `Repository`, etc.).
     - **Cache Adapters**: Located in `internal/adapter/driven/cache/`, where Redis-related interactions are implemented.
   
   - **Driver Adapters**: These adapters expose the application logic via different interfaces such as gRPC or HTTP.
     - **gRPC Services**: Located in `internal/adapter/driver/grpc/`, responsible for handling gRPC requests and mapping them to core domain logic.

### Redis and PostgreSQL Integration

- **PostgreSQL**: Acts as the **primary data store** where configuration and user-specific rate limits are stored persistently. This ensures that rate limits and user data can be persisted across system restarts and other failure conditions.
  
  - PostgreSQL schemas and migrations are located in `internal/adapter/driven/db/migration/`.
  - Rate limiting logic interacts with PostgreSQL via the **repository** pattern defined in `internal/adapter/driven/db/repository/rate.go`.
  
- **Redis**: Provides a **distributed in-memory store** to ensure global rate limits are enforced across all instances of the rate-limiting service. Redis is used to handle **global state** in a distributed environment, so that requests across different instances share the same rate-limiting information.
  
  - Redis interactions are handled in the `internal/adapter/driven/cache/` directory with a `redis.go` implementation. The global state of each user's request counts and timestamp is maintained here to ensure consistency across all instances.

### How the Rate Limiter Works
1. **RateLimit Check**: 
   - When a user makes a request, the **RateLimiterService** checks if the request should be allowed by evaluating the user's request count in the current time window. This logic is encapsulated in the service defined in `internal/adapter/driver/service/limiter.go`.
   - The sliding window algorithm is used to count requests within a time window, allowing for fairness and precision.

2. **Persistence**:
   - The PostgreSQL database stores each user's rate limit configuration and request count. This ensures that the state can be recovered if the system crashes or restarts.
   - Redis is used to ensure that the rate limit is globally consistent across multiple instances. The Redis cache is queried first to check if the user has exceeded the rate limit.

3. **Concurrency**:
   - Go's concurrency primitives (goroutines, channels) ensure that multiple requests can be handled efficiently and concurrently.
   - Redis ensures consistency of state across distributed instances of the service, so multiple instances can handle requests concurrently while respecting the global rate limit.

### Detailed Design

The rate limiter uses a **sliding window algorithm** to ensure that requests are counted within a specific time window. Here's a high-level breakdown:

- **Sliding Window**: The time window "slides" over time, and requests within the window are counted. Redis stores this state for each user, and the service checks if the user has exceeded their limit.
- **Persistence**: PostgreSQL stores persistent data on the user’s rate limit configuration and request history to handle failure recovery.

### gRPC APIs

The `RateLimiterService` provides the following gRPC APIs to interact with the rate-limiting service:

#### 1. `CheckRateLimit`
This API checks whether a user’s request should be allowed or denied based on their current rate limit.

**Request**:
```proto
message CheckRateLimitRequest {
    string user_id = 1;
    int32 limit = 2;
}
```

- `user_id`: The ID of the user making the request.
- `limit`: The rate limit to check. If this is `0`, the limit stored in the database is used.

**Response**:
```proto
message CheckRateLimitResponse {
    bool allowed = 1;
}
```

- `allowed`: A boolean indicating whether the request was allowed.

#### 2. `GetUserRateLimit`
Retrieves the current rate limit configuration for a specific user.

**Request**:
```proto
message GetUserRateLimitRequest {
    string user_id = 1;
}
```

**Response**:
```proto
message GetUserRateLimitResponse {
    int32 rate_limit = 1;
}
```

- `rate_limit`: The current rate limit for the user.

#### 3. `UpdateUserRateLimit`
Updates the rate limit for a specific user (useful for dynamic rate adjustments).

**Request**:
```proto
message UpdateUserRateLimitRequest {
    string user_id = 1;
    int32 new_limit = 2;
}
```

**Response**:
```proto
message UpdateUserRateLimitResponse {
    bool success = 1;
}
```

- `success`: Indicates whether the update was successful.

## Database Design

### PostgreSQL

The PostgreSQL database is used to store persistent information about the user's rate limits and configurations. The migrations for PostgreSQL can be found in `internal/adapter/driven/db/migration`.

The `user_rate_limits` table is designed as follows:

```sql
CREATE TABLE user_rate_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    request_count INT NOT NULL DEFAULT 0,
    rate_limit INT NOT NULL DEFAULT 100,  -- Default rate limit for users
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

This table stores:
- **user_id**: The unique identifier for the user.
- **request_count**: Tracks how many requests the user has made in the current time window.
- **rate_limit**: The user's rate limit (default is 100 requests per second).
- **timestamp**: Timestamp for the last request, which is used for sliding window calculations.

### Redis

Redis is used to store temporary data about user requests for efficient, distributed rate limiting. Redis stores:
- The user's request count within the sliding window.
- The timestamp of the user's last request to enforce the sliding window.

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

### Step 2: Running the Application

1. **Start the services with Docker**:

   ```bash
   docker-compose up --build
   ```

   This command will:
   - Set up a PostgreSQL database.
   - Spin up a Redis instance for distributed rate limiting.
   - Build and run the application service.

2. **Start the Application**:
   The application should be available at

 the configured port (e.g., `8080`).

### Step 3: Testing the Application

1. **gRPC Testing**:
   Use a tool like **BloomRPC** to test the gRPC APIs provided in the `proto/rate/v1/rate_service.proto` file.
   
   Import the proto file and call methods like `CheckRateLimit`, `GetUserRateLimit`, and `UpdateUserRateLimit`.

2. **Unit Tests and Benchmarks**:
   Run the tests using the following command:
   
   ```bash
   go test ./internal/adapter/driver/service/ -v
   ```

   This will execute all unit tests, including those that test the sliding window algorithm, Redis integration, and distributed consistency.

 