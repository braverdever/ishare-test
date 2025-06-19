#!/bin/bash

# iSHARE Task API Test Script
# This script demonstrates the complete OAuth 2.0 flow and task management

BASE_URL="http://localhost:8080"
CLIENT_ID="test-client"
CLIENT_SECRET="test-secret"
REDIRECT_URI="http://localhost:8080/oauth/callback"

echo "🚀 iSHARE Task API Test Script"
echo "================================"

# Test 1: Health Check
echo ""
echo "1. Testing Health Check..."
curl -s "$BASE_URL/health" | jq .

# Test 2: API Documentation
echo ""
echo "2. Testing API Documentation..."
curl -s "$BASE_URL/" | jq .

# Test 3: Register a new user
echo ""
echo "3. Registering a new user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/oauth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
echo $USER_RESPONSE | jq .

# Test 4: OAuth Authorization Flow
echo ""
echo "4. Starting OAuth Authorization Flow..."
echo "Visit this URL in your browser to authorize:"
echo "$BASE_URL/oauth/authorize?response_type=code&client_id=$CLIENT_ID&redirect_uri=$REDIRECT_URI&scope=tasks:read%20tasks:write&state=test-state"
echo ""

# For automated testing, we'll simulate the flow
echo "5. Simulating OAuth Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/oauth/login" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=test@example.com&password=password123&client_id=$CLIENT_ID&redirect_uri=$REDIRECT_URI&scope=tasks:read%20tasks:write&state=test-state")

# Extract authorization code from the redirect URL
AUTH_CODE=$(echo $LOGIN_RESPONSE | grep -o 'code=[^&]*' | cut -d'=' -f2)
echo "Authorization Code: $AUTH_CODE"

# Test 6: Exchange authorization code for access token
echo ""
echo "6. Exchanging authorization code for access token..."
TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/oauth/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code&code=$AUTH_CODE&redirect_uri=$REDIRECT_URI&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET")

echo $TOKEN_RESPONSE | jq .

# Extract access token
ACCESS_TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.access_token')
echo "Access Token: $ACCESS_TOKEN"

# Test 7: Create a task
echo ""
echo "7. Creating a new task..."
CREATE_TASK_RESPONSE=$(curl -s -X POST "$BASE_URL/tasks" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "title": "Complete API Documentation",
    "description": "Write comprehensive documentation for the iSHARE Task API",
    "status": "pending"
  }')
echo $CREATE_TASK_RESPONSE | jq .

# Extract task ID
TASK_ID=$(echo $CREATE_TASK_RESPONSE | jq -r '.id')
echo "Task ID: $TASK_ID"

# Test 8: Get the created task
echo ""
echo "8. Retrieving the created task..."
curl -s -X GET "$BASE_URL/tasks/$TASK_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 9: List all tasks
echo ""
echo "9. Listing all tasks..."
curl -s -X GET "$BASE_URL/tasks" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 10: Update the task
echo ""
echo "10. Updating the task..."
UPDATE_TASK_RESPONSE=$(curl -s -X PUT "$BASE_URL/tasks/$TASK_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "title": "Complete API Documentation - Updated",
    "description": "Write comprehensive documentation for the iSHARE Task API with OAuth 2.0 and JWS",
    "status": "in_progress"
  }')
echo $UPDATE_TASK_RESPONSE | jq .

# Test 11: Get the updated task
echo ""
echo "11. Retrieving the updated task..."
curl -s -X GET "$BASE_URL/tasks/$TASK_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 12: Create another task
echo ""
echo "12. Creating another task..."
curl -s -X POST "$BASE_URL/tasks" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "title": "Implement Rate Limiting",
    "description": "Add rate limiting to protect the API from abuse",
    "status": "pending"
  }' | jq .

# Test 13: List tasks with filtering
echo ""
echo "13. Listing tasks with status filter..."
curl -s -X GET "$BASE_URL/tasks?status=pending" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 14: Test pagination
echo ""
echo "14. Testing pagination..."
curl -s -X GET "$BASE_URL/tasks?page=1&limit=5" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 15: Delete the first task
echo ""
echo "15. Deleting the first task..."
curl -s -X DELETE "$BASE_URL/tasks/$TASK_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 16: Verify task is deleted
echo ""
echo "16. Verifying task is deleted..."
curl -s -X GET "$BASE_URL/tasks/$TASK_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Test 17: Test invalid token
echo ""
echo "17. Testing invalid token..."
curl -s -X GET "$BASE_URL/tasks" \
  -H "Authorization: Bearer invalid-token" | jq .

# Test 18: Test missing token
echo ""
echo "18. Testing missing token..."
curl -s -X GET "$BASE_URL/tasks" | jq .

# Test 19: Cleanup expired tokens
echo ""
echo "19. Cleaning up expired tokens..."
curl -s -X POST "$BASE_URL/oauth/cleanup" | jq .

echo ""
echo "✅ All tests completed!"
echo ""
echo "📚 API Documentation available at: $BASE_URL/swagger/index.html"
echo "🔗 Health check: $BASE_URL/health"
echo ""
echo "💡 Tips:"
echo "- Use the Swagger UI to explore the API interactively"
echo "- The OAuth flow requires user interaction in a browser"
echo "- All task endpoints require valid Bearer token authentication"
echo "- JWS tokens are signed and verified for security" 