# iSHARE Task Management API Documentation

## Overview

The iSHARE Task Management API is a secure REST API that implements OAuth 2.0 authorization with JSON Web Signatures (JWS) for token signing and verification. This API provides full CRUD operations for task management with enterprise-grade security.

## Base URL

```
http://localhost:8080
```

## Authentication

The API uses OAuth 2.0 Authorization Code Flow with JWS token signing. All task endpoints require a valid Bearer token in the Authorization header.

### OAuth 2.0 Flow

1. **Authorization Request**: Redirect user to `/oauth/authorize`
2. **User Consent**: User grants permission via login form
3. **Authorization Code**: Server returns authorization code
4. **Token Exchange**: Exchange code for access token at `/oauth/token`
5. **API Access**: Use access token for API requests

### Token Format

```
Authorization: Bearer <jws_token>
```

## Endpoints

### OAuth 2.0 Endpoints

#### 1. OAuth Authorization

**GET** `/oauth/authorize`

Initiates the OAuth 2.0 authorization code flow.

**Query Parameters:**

- `response_type` (required): Must be "code"
- `client_id` (required): OAuth client ID
- `redirect_uri` (required): Redirect URI after authorization
- `scope` (optional): Requested scopes (e.g., "tasks:read tasks:write")
- `state` (optional): State parameter for CSRF protection

**Example:**

```bash
curl "http://localhost:8080/oauth/authorize?response_type=code&client_id=test-client&redirect_uri=http://localhost:8080/oauth/callback&scope=tasks:read%20tasks:write&state=random-state"
```

**Response:** HTML login form

#### 2. OAuth Token Exchange

**POST** `/oauth/token`

Exchanges authorization code for access token.

**Form Data:**

- `grant_type` (required): Must be "authorization_code"
- `code` (required): Authorization code from previous step
- `redirect_uri` (required): Same redirect URI used in authorization
- `client_id` (required): OAuth client ID
- `client_secret` (required): OAuth client secret

**Example:**

```bash
curl -X POST "http://localhost:8080/oauth/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code&code=AUTH_CODE&redirect_uri=http://localhost:8080/oauth/callback&client_id=test-client&client_secret=test-secret"
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "scope": "tasks:read tasks:write"
}
```

#### 3. User Registration

**POST** `/oauth/register`

Creates a new user account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Example:**

```bash
curl -X POST "http://localhost:8080/oauth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Task Management Endpoints

All task endpoints require Bearer token authentication.

#### 1. Create Task

**POST** `/tasks`

Creates a new task.

**Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "title": "Complete API Documentation",
  "description": "Write comprehensive documentation for the API",
  "status": "pending"
}
```

**Example:**

```bash
curl -X POST "http://localhost:8080/tasks" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete API Documentation",
    "description": "Write comprehensive documentation for the API",
    "status": "pending"
  }'
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete API Documentation",
  "description": "Write comprehensive documentation for the API",
  "status": "pending",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 2. List Tasks

**GET** `/tasks`

Retrieves all tasks with optional filtering and pagination.

**Headers:**

```
Authorization: Bearer <access_token>
```

**Query Parameters:**

- `status` (optional): Filter by status (e.g., "pending", "completed")
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Example:**

```bash
curl -X GET "http://localhost:8080/tasks?status=pending&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**

```json
{
  "tasks": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Complete API Documentation",
      "description": "Write comprehensive documentation for the API",
      "status": "pending",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

#### 3. Get Task

**GET** `/tasks/{id}`

Retrieves a specific task by ID.

**Headers:**

```
Authorization: Bearer <access_token>
```

**Example:**

```bash
curl -X GET "http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete API Documentation",
  "description": "Write comprehensive documentation for the API",
  "status": "pending",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 4. Update Task

**PUT** `/tasks/{id}`

Updates a specific task.

**Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "title": "Complete API Documentation - Updated",
  "description": "Write comprehensive documentation with OAuth 2.0 and JWS",
  "status": "in_progress"
}
```

**Example:**

```bash
curl -X PUT "http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete API Documentation - Updated",
    "description": "Write comprehensive documentation with OAuth 2.0 and JWS",
    "status": "in_progress"
  }'
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete API Documentation - Updated",
  "description": "Write comprehensive documentation with OAuth 2.0 and JWS",
  "status": "in_progress",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### 5. Delete Task

**DELETE** `/tasks/{id}`

Deletes a specific task.

**Headers:**

```
Authorization: Bearer <access_token>
```

**Example:**

```bash
curl -X DELETE "http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**

```json
{
  "message": "Task deleted successfully"
}
```

### Utility Endpoints

#### Health Check

**GET** `/health`

Returns the health status of the API.

**Example:**

```bash
curl -X GET "http://localhost:8080/health"
```

**Response:**

```json
{
  "status": "ok",
  "message": "iSHARE Task API is running"
}
```

#### API Information

**GET** `/`

Returns API information and available endpoints.

**Example:**

```bash
curl -X GET "http://localhost:8080/"
```

**Response:**

```json
{
  "name": "iSHARE Task Management API",
  "version": "1.0.0",
  "description": "A secure REST API for task management with OAuth 2.0 and JWS token signing",
  "documentation": "/swagger/index.html",
  "endpoints": {
    "oauth": {
      "authorize": "GET /oauth/authorize - OAuth 2.0 authorization endpoint",
      "token": "POST /oauth/token - OAuth 2.0 token endpoint",
      "callback": "GET /oauth/callback - OAuth callback endpoint",
      "register": "POST /oauth/register - User registration"
    },
    "tasks": {
      "create": "POST /tasks - Create a new task",
      "list": "GET /tasks - List all tasks",
      "get": "GET /tasks/{id} - Get a specific task",
      "update": "PUT /tasks/{id} - Update a task",
      "delete": "DELETE /tasks/{id} - Delete a task"
    }
  },
  "authentication": "All task endpoints require Bearer token authentication"
}
```

## Data Models

### Task Model

```json
{
  "id": "uuid",
  "title": "string (required)",
  "description": "string",
  "status": "string (pending|in_progress|completed)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### User Model

```json
{
  "id": "uuid",
  "email": "string (required, unique)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Token Response

```json
{
  "access_token": "string",
  "token_type": "Bearer",
  "expires_in": "integer (seconds)",
  "scope": "string"
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required or invalid token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## JWS Token Structure

The API uses JSON Web Signatures (JWS) for secure token signing. Tokens have the following structure:

```
header.payload.signature
```

### Header

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Payload

```json
{
  "sub": "user_id",
  "email": "user@example.com",
  "scope": "tasks:read tasks:write",
  "iss": "ishare-task-api",
  "aud": "ishare-clients",
  "exp": 1640995200,
  "iat": 1640908800,
  "nbf": 1640908800
}
```

## Security Features

1. **OAuth 2.0 Authorization Code Flow**: Industry-standard authentication
2. **JWS Token Signing**: HMAC-SHA256 signature verification
3. **Token Expiration**: Configurable token lifetime
4. **Scope-based Authorization**: Fine-grained access control
5. **Password Hashing**: bcrypt password hashing
6. **Input Validation**: Comprehensive request validation
7. **SQL Injection Protection**: Parameterized queries with GORM

## Rate Limiting

For production deployments, implement rate limiting to protect against abuse.

## CORS

Configure CORS headers for web client access:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Testing

Use the provided test script to verify API functionality:

```bash
chmod +x test_api.sh
./test_api.sh
```

## Swagger Documentation

Interactive API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

## Support

For questions and support, refer to the iSHARE developer portal at [dev.ishare.eu](https://dev.ishare.eu).
