# Employee Module API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All endpoints require a valid Bearer token in the `Authorization` header:
```
Authorization: Bearer <your_access_token>
```

## Roles
- `USER`: Regular user
- `ADMIN`: Can manage employees
- `SUPER_USER`: Full access including delete

---

## Endpoints

### 1. Create Employee

Create a new employee with user account.

**Endpoint:** `POST /employees`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "Password123!",
  "position": "Software Engineer",
  "phoneNumber": "+628123456789",
  "address": "Jl. Sudirman No. 1, Jakarta",
  "salaryBase": 15000000,
  "joinDate": "2024-01-15",
  "bankName": "BCA",
  "bankAccountNumber": "1234567890",
  "bankAccountHolder": "John Doe",
  "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527"
}
```

**Field Descriptions:**
- `name` (required): Employee's full name
- `email` (required): Employee's email address (must be unique)
- `password` (required): Account password (min 6 characters)
- `position` (required): Job position
- `phoneNumber` (optional): Phone number
- `address` (optional): Residential address
- `salaryBase` (required): Base salary in IDR
- `joinDate` (required): Join date (format: YYYY-MM-DD)
- `bankName` (optional): Bank name
- `bankAccountNumber` (optional): Bank account number
- `bankAccountHolder` (optional): Bank account holder name
- `scheduleId` (optional): Schedule ID (UUID format)

**Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Employee created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "userId": "660e8400-e29b-41d4-a716-446655440000",
    "userName": "John Doe",
    "userEmail": "john.doe@example.com",
    "position": "Software Engineer",
    "phoneNumber": "+628123456789",
    "address": "Jl. Sudirman No. 1, Jakarta",
    "salaryBase": 15000000,
    "joinDate": "2024-01-15",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "John Doe",
    "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
    "createdAt": "2024-01-28T10:30:00Z",
    "updatedAt": "2024-01-28T10:30:00Z"
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "Email already exists",
  "error": null
}
```

---

### 2. Get All Employees

Retrieve paginated list of employees with optional search.

**Endpoint:** `GET /employees`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 15): Items per page (max: 100)
- `search` (optional): Search by name or email

**Example Request:**
```
GET /employees?page=1&per_page=15&search=John
```

**Response (200 OK):**
```json
{
  "currentPage": 1,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "userId": "660e8400-e29b-41d4-a716-446655440000",
      "userName": "John Doe",
      "userEmail": "john.doe@example.com",
      "position": "Software Engineer",
      "phoneNumber": "+628123456789",
      "address": "Jl. Sudirman No. 1, Jakarta",
      "salaryBase": 15000000,
      "joinDate": "2024-01-15",
      "bankName": "BCA",
      "bankAccountNumber": "1234567890",
      "bankAccountHolder": "John Doe",
      "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
      "createdAt": "2024-01-28T10:30:00Z",
      "updatedAt": "2024-01-28T10:30:00Z"
    }
  ],
  "firstPageUrl": "http://localhost:8080/api/v1/employees?page=1&per_page=15",
  "from": 1,
  "lastPage": 5,
  "lastPageUrl": "http://localhost:8080/api/v1/employees?page=5&per_page=15",
  "links": [
    {
      "url": null,
      "label": "pagination.previous",
      "active": false
    },
    {
      "url": "http://localhost:8080/api/v1/employees?page=1&per_page=15",
      "label": "1",
      "active": true
    },
    {
      "url": "http://localhost:8080/api/v1/employees?page=2&per_page=15",
      "label": "2",
      "active": false
    }
  ],
  "nextPageUrl": "http://localhost:8080/api/v1/employees?page=2&per_page=15",
  "path": "http://localhost:8080/api/v1/employees",
  "perPage": 15,
  "prevPageUrl": null,
  "to": 15,
  "total": 67
}
```

---

### 3. Get Employee by ID

Retrieve a specific employee by ID.

**Endpoint:** `GET /employees/:id`

**Roles Required:** `ADMIN`, `SUPER_USER`

**URL Parameters:**
- `id` (required): Employee ID (UUID)

**Example Request:**
```
GET /employees/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employee fetched successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "userId": "660e8400-e29b-41d4-a716-446655440000",
    "userName": "John Doe",
    "userEmail": "john.doe@example.com",
    "position": "Software Engineer",
    "phoneNumber": "+628123456789",
    "address": "Jl. Sudirman No. 1, Jakarta",
    "salaryBase": 15000000,
    "joinDate": "2024-01-15",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "John Doe",
    "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
    "createdAt": "2024-01-28T10:30:00Z",
    "updatedAt": "2024-01-28T10:30:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "employee not found",
  "error": null
}
```

---

### 4. Update Employee

Update an existing employee's information.

**Endpoint:** `PATCH /employees/:id`

**Roles Required:** `ADMIN`, `SUPER_USER`

**URL Parameters:**
- `id` (required): Employee ID (UUID)

**Request Body:**
```json
{
  "position": "Senior Software Engineer",
  "phoneNumber": "+628987654321",
  "address": "Jl. Sudirman No. 2, Jakarta",
  "salaryBase": 18000000,
  "bankName": "Mandiri",
  "bankAccountNumber": "0987654321",
  "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527"
}
```

**Note:** All fields are optional. Only provided fields will be updated.

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employee updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "userId": "660e8400-e29b-41d4-a716-446655440000",
    "userName": "John Doe",
    "userEmail": "john.doe@example.com",
    "position": "Senior Software Engineer",
    "phoneNumber": "+628987654321",
    "address": "Jl. Sudirman No. 2, Jakarta",
    "salaryBase": 18000000,
    "joinDate": "2024-01-15",
    "bankName": "Mandiri",
    "bankAccountNumber": "0987654321",
    "bankAccountHolder": "John Doe",
    "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
    "createdAt": "2024-01-28T10:30:00Z",
    "updatedAt": "2024-01-28T11:00:00Z"
  }
}
```

---

### 5. Delete Employee

Delete an employee and their associated user account.

**Endpoint:** `DELETE /employees/:id`

**Roles Required:** `SUPER_USER` only

**URL Parameters:**
- `id` (required): Employee ID (UUID)

**Example Request:**
```
DELETE /employees/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employee deleted successfully",
  "data": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "employee not found",
  "error": null
}
```

**Error Response (403 Forbidden):**
```json
{
  "success": false,
  "statusCode": 403,
  "message": "insufficient permissions",
  "error": null
}
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found |
| 409 | Conflict - Email already exists |
| 422 | Validation Error |
| 500 | Internal Server Error |

---

## Validation Rules

### Create Employee
- `name`: Required, min 3 characters
- `email`: Required, valid email format, must be unique
- `password`: Required, min 6 characters
- `position`: Required, min 3 characters
- `salaryBase`: Required, must be greater than 0
- `joinDate`: Required, format: YYYY-MM-DD
- `scheduleId`: Optional, valid UUID format

### Update Employee
- `position`: Optional, min 3 characters
- `salaryBase`: Optional, must be greater than 0
- `scheduleId`: Optional, valid UUID format

---

## Transaction Behavior

### Create Employee
The create operation uses **Saga Pattern** for distributed transaction:

1. **Step 1**: Create User in MongoDB
   - If fails → Return error immediately
   
2. **Step 2**: Create Employee in PostgreSQL
   - If fails → **Compensating Action**: Delete User from MongoDB
   - Return error "Failed to create employee profile"

3. **Step 3**: Return success with combined data

### Delete Employee
1. Delete Employee record from PostgreSQL
2. Delete associated User from MongoDB
3. If MongoDB deletion fails → Log error but don't fail (employee already deleted)

---

## Postman Collection JSON

Import this JSON into Postman:

```json
{
  "info": {
    "name": "Employee Module API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api/v1",
      "type": "string"
    },
    {
      "key": "accessToken",
      "value": "your_access_token_here",
      "type": "string"
    }
  ],
  "item": [
    {
      "name": "Create Employee",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": "{{baseUrl}}/employees",
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"John Doe\",\n  \"email\": \"john.doe@example.com\",\n  \"password\": \"Password123!\",\n  \"position\": \"Software Engineer\",\n  \"phoneNumber\": \"+628123456789\",\n  \"address\": \"Jl. Sudirman No. 1, Jakarta\",\n  \"salaryBase\": 15000000,\n  \"joinDate\": \"2024-01-15\",\n  \"bankName\": \"BCA\",\n  \"bankAccountNumber\": \"1234567890\",\n  \"bankAccountHolder\": \"John Doe\",\n  \"scheduleId\": \"8918b101-f313-4ef8-bbcb-c72ebfae3527\"\n}"
        }
      }
    },
    {
      "name": "Get All Employees",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/employees?page=1&per_page=15&search=",
          "query": [
            {
              "key": "page",
              "value": "1"
            },
            {
              "key": "per_page",
              "value": "15"
            },
            {
              "key": "search",
              "value": ""
            }
          ]
        }
      }
    },
    {
      "name": "Get Employee by ID",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": "{{baseUrl}}/employees/:id"
      }
    },
    {
      "name": "Update Employee",
      "request": {
        "method": "PATCH",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": "{{baseUrl}}/employees/:id",
        "body": {
          "mode": "raw",
          "raw": "{\n  \"position\": \"Senior Software Engineer\",\n  \"phoneNumber\": \"+628987654321\",\n  \"address\": \"Jl. Sudirman No. 2, Jakarta\",\n  \"salaryBase\": 18000000,\n  \"bankName\": \"Mandiri\",\n  \"bankAccountNumber\": \"0987654321\"\n}"
        }
      }
    },
    {
      "name": "Delete Employee",
      "request": {
        "method": "DELETE",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": "{{baseUrl}}/employees/:id"
      }
    }
  ]
}
```

---

## Notes

1. **Authentication**: First, you need to login using `/api/v1/auth/login` to get the access token
2. **Schedule ID**: Use the schedule ID from `/api/v1/schedules` endpoint
3. **Pagination**: Use Laravel-style pagination for consistent API behavior
4. **Search**: Search is case-insensitive and searches in both name and email fields
5. **Distributed Transaction**: Employee creation spans two databases (MongoDB for User, PostgreSQL for Employee)
