# Room Planner Backend

Go-based backend service for the **Room Planner** application. Handles business logic, database operations, and file storage.

## Tech Stack

- **Language:** Go 1.25.1
- **Database:** MariaDB (see [rp-database](https://github.com/kvilmos/rp-database))
- **Object Storage:** MinIO (S3 compatible)
- **Cache:** Redis
- **Infra:** Docker & Docker Compose

## Prerequisites

- [Go](https://go.dev/doc/install) (1.25.1+)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Docker Compose](https://docs.docker.com/compose/)
- Running instance of [rp-database](https://github.com/kvilmos/rp-database)

## Getting Started

### 1. Start the Database

Ensure the main database is running as described in the `rp-database` repository.

- [Link to rp-database](https://github.com/kvilmos/rp-database)

### 2. Start Dependencies (MinIO & Redis)

The backend requires MinIO for file storage and Redis for caching. Start the pre-configured containers:

**Start MinIO:**
A helper service will automatically create the required buckets and webhooks.

```bash
cd dev-minio
docker-compose up -d
```

_Console: [http://localhost:9001](http://localhost:9001) (User/Pass: `minioadmin`)_

**Start Redis:**

```bash
cd dev-redis
docker-compose up -d
```

### 3. Run the Backend

Once MariaDB, MinIO, and Redis are running:

```bash
# From project root
go mod tidy
go run .
```

---

## Default Configuration

The application attempts to connect to services using these default settings:

| Service   | Connection String / Host                           |
| :-------- | :------------------------------------------------- |
| **MySQL** | `develop:develop@tcp(127.0.0.1:3306)/room-planner` |
| **MinIO** | `127.0.0.1:9000`                                   |
| **Redis** | `localhost:6379`                                   |

---

## Useful Commands

| Context        | Command               | Description                    |
| :------------- | :-------------------- | :----------------------------- |
| **Start**      | `go run .`            | Run the Go application locally |
| **Stop MinIO** | `docker-compose down` | Execute inside `dev-minio/`    |
| **Stop Redis** | `docker-compose down` | Execute inside `dev-redis/`    |
