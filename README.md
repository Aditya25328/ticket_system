# Ticket System Backend in Go

A robust and lightweight RESTful Ticket System API built using **Golang**, the **Gin web framework**, and **GORM** with a **SQLite** database.

## Features

- **User Authentication**: Secure register and login system using bcrypt for password hashing and stateless **JWT** (JSON Web Tokens) for request authorization.
- **Ownership-based Access**: Users can only create, view, and update status for tickets they personally created.
- **Ticket Status Flow**: Enforces a strict one-way lifecycle flow:
  - `open` -> `in_progress` -> `closed`
  - Prevents reopening of `closed` tickets.
- **Dockerized Ready**: Fully containerized multi-stage Docker build config for deployment.
- **Robust Edge Case Handling**: Case-insensitive Authorization header parsing and exact JSON key matching.

---

## Submission Details

- **GitHub Repository Link**: https://github.com/Aditya25328/ticket_system.git
- **Deployed Application URL**: https://ticket-system-1qss.onrender.com
- **Public Health Check URL**: https://ticket-system-1qss.onrender.com/health

---

## Assumptions & Design Choices

1. **Pure-Go SQLite Driver**: Utilizes `github.com/glebarez/sqlite` which is a CGO-free SQLite implementation. This allows for lightweight multi-stage Docker builds without requiring a GCC toolchain in Alpine.
2. **Linear Status Flow**: A ticket can transition from `open` directly to `closed` or via `in_progress` to `closed`. Any transitions moving backwards (e.g. `closed -> open`, `in_progress -> open`) are rejected with `400 Bad Request`.
3. **Stateless JWT**: Authentication uses a stateless JWT claim validation structure, retrieving the secret key from the environment variable (`JWT_SECRET`).

---

## Tech Stack

- **Language**: Go 1.26
- **Router**: Gin Web Framework (v1.12.0)
- **ORM**: GORM (v1.31.2)
- **Database**: SQLite (CGO-free wrapper)
- **Auth**: JWT (golang-jwt/jwt/v5) and Bcrypt (golang.org/x/crypto/bcrypt)

---

## API Endpoints

### Public Endpoints

| Method | Endpoint | Purpose | Request Body | Response |
| :--- | :--- | :--- | :--- | :--- |
| **GET** | `/health` | Health Check | *None* | `{"status": "ok"}` |
| **POST** | `/auth/register` | Register new user | `{"name": "...", "email": "...", "password": "..."}` | Registered user details |
| **POST** | `/auth/login` | Login and return JWT | `{"email": "...", "password": "..."}` | `{"token": "JWT_TOKEN"}` |

### Protected Endpoints (Requires `Authorization: Bearer <token>`)

| Method | Endpoint | Purpose | Request Body | Response |
| :--- | :--- | :--- | :--- | :--- |
| **POST** | `/tickets` | Create a new ticket | `{"title": "...", "description": "..."}` | Created ticket details |
| **GET** | `/tickets` | List logged-in user tickets | *None* | Array of user tickets |
| **GET** | `/tickets/{id}` | Get ticket by ID | *None* | Ticket details |
| **PATCH** | `/tickets/{id}/status` | Update ticket status | `{"status": "..."}` | Updated ticket details |

---

## Getting Started

### Local Setup & Execution

1. **Clone the repository**:
   ```bash
   git clone <repo_url>
   cd ticket-system
   ```

2. **Configure Environment Variables**:
   Create a `.env` file based on `.env.example`:
   ```bash
   cp .env.example .env
   ```

3. **Run the application**:
   ```bash
   go run cmd/main.go
   ```
   The server will start on port `8080`.

4. **Verify server is running**:
   ```bash
   curl http://localhost:8080/health
   ```

### Docker Setup

To run the application inside a Docker container:

1. **Build the Docker Image**:
   ```bash
   docker build -t ticket-system .
   ```

2. **Run the Container**:
   ```bash
   docker run -p 8080:8080 ticket-system
   ```

3. **Verify the containerized app**:
   ```bash
   curl http://localhost:8080/health
   ```
