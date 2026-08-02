# StreamFlix Backend

A production-oriented AI-powered movie recommendation backend built with Go, PostgreSQL, and pgvector.

StreamFlix provides semantic movie search and content-based recommendations by leveraging vector embeddings. Instead of relying solely on keyword matching, it enables natural language search and similarity-based recommendations using embedding vectors stored directly in PostgreSQL.

---

# Features

* Semantic movie search using vector embeddings
* Content-based movie recommendations
* PostgreSQL with pgvector integration
* Modular monolith architecture
* Repository → Service → Handler architecture
* Dependency Injection
* Concurrent embedding generation pipeline
* Configurable worker pool
* Configurable batch processing
* Retry mechanism with exponential backoff
* SQL-based database migrations
* Prometheus metrics
* Structured JSON logging
* Graceful shutdown
* Production-ready HTTP server configuration

---

# Tech Stack

| Category        | Technology                 |
| --------------- | -------------------------- |
| Language        | Go                         |
| HTTP Framework  | Gin                        |
| Database        | PostgreSQL                 |
| Vector Database | pgvector                   |
| ORM             | GORM                       |
| AI Embeddings   | External Embedding Service |
| Monitoring      | Prometheus                 |
| Logging         | slog                       |

---

# Architecture

The project follows a layered architecture.

```text
                HTTP Request
                      │
                      ▼
                 Gin Router
                      │
                      ▼
                  Handler Layer
                      │
                      ▼
                  Service Layer
                      │
                      ▼
                Repository Layer
                      │
                      ▼
             PostgreSQL + pgvector
```

Each layer has a single responsibility.

* **Handlers** process HTTP requests and responses.
* **Services** contain business logic.
* **Repositories** encapsulate database access.
* **PostgreSQL** stores both movie metadata and embedding vectors.

---

# Project Structure

```text
.
├── cmd
│   ├── api
│   └── ingest
├── data
├── embedding-service
├── internal
│   ├── config
│   ├── database
│   ├── embedding
│   ├── ingestion
│   ├── metrics
│   ├── movie
│   ├── recommendation
│   ├── shared
│   └── watchhistory
├── migrations
└── scripts
```

---

# Semantic Search

Semantic search converts a natural language query into an embedding vector and retrieves the closest movie embeddings using PostgreSQL's cosine similarity operator.

Example queries:

* Space exploration with emotional storytelling
* Mind-bending science fiction
* Crime thriller with psychological drama

Endpoint:

```http
GET /search?q=<query>
```

---

# Recommendation Engine

Recommendations are generated using stored movie embeddings.

The backend retrieves the selected movie's embedding and performs a nearest-neighbor search using pgvector.

Endpoint:

```http
GET /movies/{id}/recommendations?limit=5
```

---

# Embedding Pipeline

The ingestion pipeline is designed for large datasets.

Features include:

* Configurable worker pool
* Batch processing
* Context propagation
* Concurrent embedding generation
* Retry mechanism
* Structured logging

Configuration:

```env
EMBEDDING_WORKERS=5
EMBEDDING_BATCH_SIZE=100
```

---

# API Endpoints

| Method | Endpoint                       | Description                   |
| ------ | ------------------------------ | ----------------------------- |
| GET    | `/health`                      | Health check                  |
| GET    | `/metrics`                     | Prometheus metrics            |
| GET    | `/movies`                      | Retrieve movies               |
| GET    | `/search?q=`                   | Semantic movie search         |
| GET    | `/movies/{id}/recommendations` | Content-based recommendations |

---

# Database

The application uses PostgreSQL with pgvector.

Movie embeddings are stored directly within the `movies` table using the `vector(384)` data type.

Similarity search is performed using cosine distance.

---

# Configuration

Create a `.env` file based on `.env.example`.

Example:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=streamflix
DB_SSLMODE=disable

EMBEDDING_SERVICE_URL=http://localhost:8001

EMBEDDING_WORKERS=5
EMBEDDING_BATCH_SIZE=100
```

---

# Running the Project

## 1. Clone the repository

```bash
git clone https://github.com/<your-username>/streamflix-backend.git
cd streamflix-backend
```

## 2. Start PostgreSQL

Ensure PostgreSQL with the pgvector extension is running.

## 3. Configure environment variables

```bash
cp .env.example .env
```

Update the database credentials in `.env`.

## 4. Start the embedding service

```bash
cd embedding-service

python -m venv venv

source venv/bin/activate

pip install -r requirements.txt

python main.py
```

## 5. Run database migrations

The application automatically executes pending SQL migrations during startup.

## 6. Import the dataset

```bash
go run scripts/import_movies.go
```

## 7. Generate embeddings

```bash
go run cmd/ingest/main.go
```

## 8. Start the API

```bash
go run cmd/api/main.go
```

---

# Observability

Prometheus metrics are exposed at:

```text
/metrics
```

Current metrics include:

* HTTP request count
* HTTP request duration

---

# Design Principles

* Modular Monolith
* Clean Layered Architecture
* Dependency Injection
* Repository Pattern
* Context Propagation
* Production-Oriented Configuration
* SQL-Based Migrations
* Separation of Concerns

---

# Future Improvements

* AI-specific Prometheus metrics
* Redis caching
* Authentication
* Personalized recommendations
* Hybrid search
* Docker production image
* GitHub Actions CI/CD
* OpenTelemetry tracing
* Kubernetes deployment

---

# License

This project is licensed under the MIT License.
