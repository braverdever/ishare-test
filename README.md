# iSHARE Task Management API

A secure REST API built in Go that implements OAuth 2.0 authorization with JSON Web Signatures (JWS) for token signing and verification. This API provides task management functionality with secure authentication and authorization.

## Features

- **OAuth 2.0 Authorization Code Flow**: Secure authentication using industry-standard OAuth 2.0
- **JWS Token Signing**: JSON Web Signatures for secure token verification
- **Task Management**: Full CRUD operations for tasks
- **PostgreSQL Database**: Reliable data persistence
- **Swagger Documentation**: Interactive API documentation
- **Security**: JWT-based authentication with signature verification

## Architecture

This API follows the iSHARE protocol principles for secure API communication where the consumer and provider don't necessarily have pre-shared secrets. Instead, tokens are signed using JWS for verification.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Setup Instructions

### 1. Clone the Repository

```bash
git clone <repository-url>
cd ishare-task-api
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Database Setup

Create a PostgreSQL database and update the connection string in `config/config.go`:

```bash
# Create database
createdb ishare_tasks

# Or using psql
psql -c "CREATE DATABASE ishare_tasks;"
```

### 4. Environment Configuration

Create a `.env` file in the root directory:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ishare_tasks
DB_SSL_MODE=disable

JWT_SECRET=your_super_secret_jwt_key_here
JWT_ISSUER=ishare-task-api
JWT_AUDIENCE=ishare-clients

OAUTH_CLIENT_ID=your_oauth_client_id
OAUTH_CLIENT_SECRET=your_oauth_client_secret
OAUTH_REDIRECT_URI=http://localhost:8080/oauth/callback

SERVER_PORT=8080
```

### 5. Generate Swagger Documentation

```bash
swag init -g main.go
```

### 6. Run the Application

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### OAuth 2.0 Endpoints

- `GET /oauth/authorize` - OAuth authorization endpoint
- `POST /oauth/token` - OAuth token endpoint
- `GET /oauth/callback` - OAuth callback endpoint

### Task Management Endpoints

All task endpoints require valid JWT token in Authorization header: `Authorization: Bearer <token>`

- `POST /tasks` - Create a new task
- `GET /tasks` - List all tasks
- `GET /tasks/{id}` - Get a specific task
- `PUT /tasks/{id}` - Update a task
- `DELETE /tasks/{id}` - Delete a task

### API Documentation

- Swagger UI: `http://localhost:8080/swagger/index.html`
- API Docs: `http://localhost:8080/swagger/doc.json`

## Authentication Flow

1. **Authorization Request**: Client redirects user to `/oauth/authorize`
2. **User Consent**: User grants permission to the application
3. **Authorization Code**: Server returns authorization code to client
4. **Token Exchange**: Client exchanges code for access token at `/oauth/token`
5. **API Access**: Client uses access token for API requests

## JWS Token Structure

Tokens are signed using JWS with the following structure:

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_id",
    "iss": "ishare-task-api",
    "aud": "ishare-clients",
    "exp": 1640995200,
    "iat": 1640908800,
    "scope": "tasks:read tasks:write"
  },
  "signature": "base64_encoded_signature"
}
```

## Database Schema

### Tasks Table

```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Security Considerations

1. **JWS Signing**: All tokens are signed using HMAC-SHA256
2. **Token Expiration**: Tokens have configurable expiration times
3. **Scope-based Authorization**: Fine-grained access control
4. **HTTPS**: Production deployment should use HTTPS
5. **Rate Limiting**: Implement rate limiting for production use
6. **Input Validation**: All inputs are validated and sanitized

## Assumptions Made

1. **OAuth Provider**: The implementation assumes a simple OAuth provider. In production, you might integrate with external providers like Google, GitHub, etc.
2. **Database**: PostgreSQL is used, but the code can be easily adapted for other databases
3. **Token Storage**: Tokens are stored in memory for simplicity. Production should use Redis or similar
4. **User Management**: Basic user management is implemented. Production should include password reset, email verification, etc.
5. **Error Handling**: Basic error handling is implemented. Production should include comprehensive logging and monitoring

## Development

### Project Structure

```
ishare-task-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── jwt.go
│   │   ├── oauth.go
│   │   └── middleware.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── auth.go
│   │   └── tasks.go
│   ├── models/
│   │   ├── task.go
│   │   └── user.go
│   └── routes/
│       └── routes.go
├── docs/
├── go.mod
├── go.sum
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o bin/server cmd/server/main.go
```

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Support

For support and questions, please refer to the iSHARE developer portal at [dev.ishare.eu](https://dev.ishare.eu). 