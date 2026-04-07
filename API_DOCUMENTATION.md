# Go-Fiber-Boilerplate API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All endpoints (except public auth endpoints) require a valid Bearer token in the `Authorization` header:
```
Authorization: Bearer <your_access_token>
```

## Standard Response Format

### Success Response (Single Data)
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Success message",
  "data": {
    "id": "uuid",
    "field": "value"
  }
}
```

### Success Response (With Pagination)
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Data retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "field": "value"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 10,
    "total": 100,
    "lastPage": 10,
    "firstPageUrl": "http://localhost:8080/api/v1/resource?page=1&per_page=10",
    "lastPageUrl": "http://localhost:8080/api/v1/resource?page=10&per_page=10",
    "nextPageUrl": "http://localhost:8080/api/v1/resource?page=2&per_page=10",
    "prevPageUrl": null
  }
}
```

### Error Response
```json
{
  "success": false,
  "statusCode": 400,
  "message": "Error message",
  "error": "Detailed error or null"
}
```

### Validation Error Response
```json
{
  "success": false,
  "statusCode": 422,
  "message": "Validation failed",
  "error": [
    {
      "field": "email",
      "message": "email is required"
    }
  ]
}
```

## Roles
- **USER**: Regular employee user
- **ADMIN**: Can manage users, employees, schedules, view all attendances and payrolls
- **SUPER_USER**: Full access including delete operations

---

## 1. AUTH MODULE

### 1.1 Register
**Endpoint:** `POST /auth/register`

**Access:** Public

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "USER"
}
```

**Field Validation:**
- `name`: Required, string, min 3 characters
- `email`: Required, valid email format, unique
- `password`: Required, string, min 6 characters
- `role`: Required, one of: `USER`, `ADMIN`, `SUPER_USER`

**Success Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "USER",
      "profileImage": null,
      "createdAt": "2024-01-30T10:00:00Z"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "email already exists",
  "error": null
}
```

---

### 1.2 Login
**Endpoint:** `POST /auth/login`

**Access:** Public

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "USER",
      "profileImage": null,
      "createdAt": "2024-01-30T10:00:00Z"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "statusCode": 401,
  "message": "invalid credentials",
  "error": null
}
```

---

### 1.3 Refresh Token
**Endpoint:** `POST /auth/refresh`

**Access:** Public

**Request Body:**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Token refreshed successfully",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600
  }
}
```

---

### 1.4 Logout
**Endpoint:** `POST /auth/logout`

**Access:** Authenticated (All roles)

**Request Body:**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Logout successful",
  "data": null
}
```

---

### 1.5 Logout All Devices
**Endpoint:** `POST /auth/logout-all`

**Access:** Authenticated (All roles)

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Logged out from all devices successfully",
  "data": null
}
```

---

## 2. USER MODULE

### 2.1 Get Own Profile
**Endpoint:** `GET /users/me`

**Access:** USER, ADMIN, SUPER_USER

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "USER",
    "profileImage": "http://localhost:9000/users/uuid/profile.jpg",
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 2.2 Get All Users (Paginated)
**Endpoint:** `GET /users`

**Access:** ADMIN, SUPER_USER

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 10, max: 100): Items per page

**Example Request:**
```
GET /users?page=1&per_page=10
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "USER",
      "profileImage": null,
      "createdAt": "2024-01-30T10:00:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 10,
    "total": 50,
    "lastPage": 5,
    "firstPageUrl": "http://localhost:8080/api/v1/users?page=1&per_page=10",
    "lastPageUrl": "http://localhost:8080/api/v1/users?page=5&per_page=10",
    "nextPageUrl": "http://localhost:8080/api/v1/users?page=2&per_page=10",
    "prevPageUrl": null
  }
}
```

---

### 2.3 Get User by ID
**Endpoint:** `GET /users/:id`

**Access:** ADMIN, SUPER_USER

**URL Parameters:**
- `id`: User UUID

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "USER",
    "profileImage": null,
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 2.4 Create User
**Endpoint:** `POST /users`

**Access:** ADMIN, SUPER_USER

**Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "password": "password123",
  "role": "USER"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "User created successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Jane Smith",
    "email": "jane@example.com",
    "role": "USER",
    "profileImage": null,
    "createdAt": "2024-01-30T11:00:00Z"
  }
}
```

---

### 2.5 Update User
**Endpoint:** `PATCH /users/:id`

**Access:** ADMIN, SUPER_USER

**Request Body (all fields optional):**
```json
{
  "name": "Jane Updated",
  "email": "jane.updated@example.com",
  "role": "ADMIN"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User updated successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Jane Updated",
    "email": "jane.updated@example.com",
    "role": "ADMIN",
    "profileImage": null,
    "createdAt": "2024-01-30T11:00:00Z"
  }
}
```

---

### 2.6 Delete User
**Endpoint:** `DELETE /users/:id`

**Access:** SUPER_USER only

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User deleted successfully",
  "data": null
}
```

---

### 2.7 Upload Profile Image
**Endpoint:** `POST /users/:id/profile-image`

**Access:**
- USER: Can upload own profile image only
- ADMIN, SUPER_USER: Can upload any user's profile image

**Content-Type:** `multipart/form-data`

**Form Data:**
- `image`: Image file (JPG, PNG, max 2MB)

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Profile image uploaded successfully",
  "data": {
    "imageUrl": "http://localhost:9000/users/uuid/profile.jpg"
  }
}
```

---

### 2.8 Update Profile Image
**Endpoint:** `PUT /users/:id/profile-image`

**Access:**
- USER: Can update own profile image only
- ADMIN, SUPER_USER: Can update any user's profile image

**Content-Type:** `multipart/form-data`

**Form Data:**
- `image`: Image file (JPG, PNG, max 2MB)

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Profile image updated successfully",
  "data": {
    "imageUrl": "http://localhost:9000/users/uuid/profile.jpg"
  }
}
```

---

### 2.9 Delete Profile Image
**Endpoint:** `DELETE /users/:id/profile-image`

**Access:**
- USER: Can delete own profile image only
- ADMIN, SUPER_USER: Can delete any user's profile image

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Profile image deleted successfully",
  "data": null
}
```

---

## 3. EMPLOYEE MODULE

### 3.1 Get All Employees (Paginated)
**Endpoint:** `GET /employees`

**Access:** ADMIN, SUPER_USER

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 10, max: 100): Items per page
- `search` (optional): Search by name or email

**Example Request:**
```
GET /employees?page=1&per_page=10&search=john
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employees retrieved successfully",
  "data": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "628123456789",
      "address": "Jl. Sudirman No. 1",
      "salaryBase": 15000000,
      "joinDate": "2024-01-01T00:00:00Z",
      "scheduleId": "880e8400-e29b-41d4-a716-446655440003",
      "scheduleName": "Office Hours",
      "createdAt": "2024-01-30T10:00:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 10,
    "total": 25,
    "lastPage": 3,
    "firstPageUrl": "http://localhost:8080/api/v1/employees?page=1&per_page=10",
    "lastPageUrl": "http://localhost:8080/api/v1/employees?page=3&per_page=10",
    "nextPageUrl": "http://localhost:8080/api/v1/employees?page=2&per_page=10",
    "prevPageUrl": null
  }
}
```

---

### 3.2 Get Employee by ID
**Endpoint:** `GET /employees/:id`

**Access:** ADMIN, SUPER_USER

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employee fetched successfully",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "628123456789",
    "address": "Jl. Sudirman No. 1",
    "salaryBase": 15000000,
    "joinDate": "2024-01-01T00:00:00Z",
    "scheduleId": "880e8400-e29b-41d4-a716-446655440003",
    "scheduleName": "Office Hours",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "John Doe",
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 3.3 Create Employee
**Endpoint:** `POST /employees`

**Access:** ADMIN, SUPER_USER

**Request Body:**
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "phone": "628123456789",
  "address": "Jl. Sudirman No. 1",
  "salaryBase": 15000000,
  "joinDate": "2024-01-01",
  "scheduleId": "880e8400-e29b-41d4-a716-446655440003",
  "bankName": "BCA",
  "bankAccountNumber": "1234567890",
  "bankAccountHolder": "John Doe"
}
```

**Field Validation:**
- `userId`: Required, valid UUID (must exist in users table)
- `phone`: Required, string
- `address`: Required, string
- `salaryBase`: Required, integer > 0
- `joinDate`: Required, date format YYYY-MM-DD
- `scheduleId`: Optional, valid UUID
- `bankName`: Required, string
- `bankAccountNumber`: Required, string
- `bankAccountHolder`: Required, string

**Success Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Employee created successfully",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "628123456789",
    "address": "Jl. Sudirman No. 1",
    "salaryBase": 15000000,
    "joinDate": "2024-01-01T00:00:00Z",
    "scheduleId": "880e8400-e29b-41d4-a716-446655440003",
    "scheduleName": "Office Hours",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "John Doe",
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 3.4 Update Employee
**Endpoint:** `PATCH /employees/:id`

**Access:** ADMIN, SUPER_USER

**Request Body (all fields optional):**
```json
{
  "phone": "628987654321",
  "address": "Jl. Thamrin No. 2",
  "salaryBase": 16000000,
  "scheduleId": "880e8400-e29b-41d4-a716-446655440004",
  "bankName": "Mandiri",
  "bankAccountNumber": "0987654321",
  "bankAccountHolder": "John Doe"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employee updated successfully",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "628987654321",
    "address": "Jl. Thamrin No. 2",
    "salaryBase": 16000000,
    "joinDate": "2024-01-01T00:00:00Z",
    "scheduleId": "880e8400-e29b-41d4-a716-446655440004",
    "scheduleName": "Flexible Hours",
    "bankName": "Mandiri",
    "bankAccountNumber": "0987654321",
    "bankAccountHolder": "John Doe",
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 3.5 Delete Employee
**Endpoint:** `DELETE /employees/:id`

**Access:** SUPER_USER only

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employee deleted successfully",
  "data": null
}
```

---

## 4. SCHEDULE MODULE

### 4.1 Get All Schedules (Paginated)
**Endpoint:** `GET /schedules`

**Access:** ADMIN, SUPER_USER

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 10, max: 100): Items per page

**Example Request:**
```
GET /schedules?page=1&per_page=10
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Schedules retrieved successfully",
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "name": "Office Hours",
      "timeIn": "09:00",
      "timeOut": "17:00",
      "allowedLateMinutes": 15,
      "officeLat": -6.2088,
      "officeLong": 106.8456,
      "allowedRadiusMeters": 100,
      "createdAt": "2024-01-30T10:00:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 10,
    "total": 5,
    "lastPage": 1,
    "firstPageUrl": "http://localhost:8080/api/v1/schedules?page=1&per_page=10",
    "lastPageUrl": "http://localhost:8080/api/v1/schedules?page=1&per_page=10",
    "nextPageUrl": null,
    "prevPageUrl": null
  }
}
```

---

### 4.2 Get Schedule by ID
**Endpoint:** `GET /schedules/:id`

**Access:** ADMIN, SUPER_USER

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Schedule fetched successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "name": "Office Hours",
    "timeIn": "09:00",
    "timeOut": "17:00",
    "allowedLateMinutes": 15,
    "officeLat": -6.2088,
    "officeLong": 106.8456,
    "allowedRadiusMeters": 100,
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 4.3 Create Schedule
**Endpoint:** `POST /schedules`

**Access:** SUPER_USER only

**Request Body:**
```json
{
  "name": "Office Hours",
  "timeIn": "09:00",
  "timeOut": "17:00",
  "allowedLateMinutes": 15,
  "officeLat": -6.2088,
  "officeLong": 106.8456,
  "allowedRadiusMeters": 100
}
```

**Field Validation:**
- `name`: Required, string
- `timeIn`: Required, time format HH:MM
- `timeOut`: Required, time format HH:MM
- `allowedLateMinutes`: Required, integer >= 0
- `officeLat`: Required, float (latitude)
- `officeLong`: Required, float (longitude)
- `allowedRadiusMeters`: Required, integer > 0

**Success Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Schedule created successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "name": "Office Hours",
    "timeIn": "09:00",
    "timeOut": "17:00",
    "allowedLateMinutes": 15,
    "officeLat": -6.2088,
    "officeLong": 106.8456,
    "allowedRadiusMeters": 100,
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 4.4 Update Schedule
**Endpoint:** `PATCH /schedules/:id`

**Access:** SUPER_USER only

**Request Body (all fields optional):**
```json
{
  "name": "Office Hours Updated",
  "timeIn": "08:00",
  "timeOut": "16:00",
  "allowedLateMinutes": 30,
  "officeLat": -6.2000,
  "officeLong": 106.8500,
  "allowedRadiusMeters": 150
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Schedule updated successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "name": "Office Hours Updated",
    "timeIn": "08:00",
    "timeOut": "16:00",
    "allowedLateMinutes": 30,
    "officeLat": -6.2000,
    "officeLong": 106.8500,
    "allowedRadiusMeters": 150,
    "createdAt": "2024-01-30T10:00:00Z"
  }
}
```

---

### 4.5 Delete Schedule
**Endpoint:** `DELETE /schedules/:id`

**Access:** SUPER_USER only

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Schedule deleted successfully",
  "data": null
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "schedule is in use by employees",
  "error": null
}
```

---

## 5. ATTENDANCE MODULE

### 5.1 Clock In
**Endpoint:** `POST /attendance/clock-in`

**Access:** USER only (Admin & Super User forbidden)

**Request Body:**
```json
{
  "lat": -6.2088,
  "long": 106.8456
}
```

**Field Validation:**
- `lat`: Required, float (latitude)
- `long`: Required, float (longitude)

**Success Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Clock in successful",
  "data": {
    "attendanceId": "990e8400-e29b-41d4-a716-446655440005",
    "employeeId": "770e8400-e29b-41d4-a716-446655440002",
    "clockInTime": "2024-01-30T09:05:00Z",
    "status": "PRESENT",
    "distance": 45.5,
    "scheduleName": "Office Hours"
  }
}
```

**Error Response (403 Forbidden):**
```json
{
  "success": false,
  "statusCode": 403,
  "message": "Admin and Super User do not need to clock in",
  "error": null
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "already clocked in today",
  "error": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "out of office range",
  "error": null
}
```

---

### 5.2 Clock Out
**Endpoint:** `POST /attendance/clock-out`

**Access:** USER only (Admin & Super User forbidden)

**Request Body:**
```json
{
  "lat": -6.2088,
  "long": 106.8456
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Clock out successful",
  "data": {
    "attendanceId": "990e8400-e29b-41d4-a716-446655440005",
    "clockOutTime": "2024-01-30T17:05:00Z",
    "distance": 48.2
  }
}
```

---

### 5.3 Get Own History (Paginated)
**Endpoint:** `GET /attendance/history`

**Access:** USER, ADMIN, SUPER_USER

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 10, max: 100): Items per page

**Example Request:**
```
GET /attendance/history?page=1&per_page=10
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Attendance history retrieved successfully",
  "data": [
    {
      "id": "990e8400-e29b-41d4-a716-446655440005",
      "date": "2024-01-30T00:00:00Z",
      "clockInTime": "2024-01-30T09:05:00Z",
      "clockOutTime": "2024-01-30T17:05:00Z",
      "status": "PRESENT",
      "notes": "",
      "scheduleName": "Office Hours"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 10,
    "total": 20,
    "lastPage": 2,
    "firstPageUrl": "http://localhost:8080/api/v1/attendance/history?page=1&per_page=10",
    "lastPageUrl": "http://localhost:8080/api/v1/attendance/history?page=2&per_page=10",
    "nextPageUrl": "http://localhost:8080/api/v1/attendance/history?page=2&per_page=10",
    "prevPageUrl": null
  }
}
```

---

### 5.4 Get All Attendances (Paginated) - Admin Only
**Endpoint:** `GET /attendance/all`

**Access:** ADMIN, SUPER_USER only

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 10, max: 100): Items per page
- `employee_id` (optional): Filter by employee UUID
- `schedule_id` (optional): Filter by schedule UUID
- `status` (optional): Filter by status (PRESENT, LATE)
- `date_from` (optional): Filter date from (YYYY-MM-DD)
- `date_to` (optional): Filter date to (YYYY-MM-DD)

**Example Requests:**
```
# Get all attendances
GET /attendance/all?page=1&per_page=10

# Filter by employee
GET /attendance/all?employee_id=770e8400-e29b-41d4-a716-446655440002

# Filter by status
GET /attendance/all?status=LATE

# Filter by date range
GET /attendance/all?date_from=2024-01-01&date_to=2024-01-31

# Combined filters
GET /attendance/all?employee_id=uuid&status=PRESENT&date_from=2024-01-01
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Attendances retrieved successfully",
  "data": [
    {
      "id": "990e8400-e29b-41d4-a716-446655440005",
      "date": "2024-01-30T00:00:00Z",
      "clockInTime": "2024-01-30T09:05:00Z",
      "clockOutTime": "2024-01-30T17:05:00Z",
      "status": "PRESENT",
      "notes": "",
      "scheduleName": "Office Hours"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 10,
    "total": 100,
    "lastPage": 10,
    "firstPageUrl": "http://localhost:8080/api/v1/attendance/all?page=1&per_page=10",
    "lastPageUrl": "http://localhost:8080/api/v1/attendance/all?page=10&per_page=10",
    "nextPageUrl": "http://localhost:8080/api/v1/attendance/all?page=2&per_page=10",
    "prevPageUrl": null
  }
}
```

---

## 6. PAYROLL MODULE

### 6.1 Generate Bulk Payroll
**Endpoint:** `POST /payrolls/generate`

**Access:** ADMIN, SUPER_USER

**Request Body:**
```json
{
  "periodMonth": 1,
  "periodYear": 2024
}
```

**Field Validation:**
- `periodMonth`: Required, integer 1-12
- `periodYear`: Required, integer 2020-2100

**Success Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Payroll generated successfully",
  "data": {
    "totalGenerated": 15,
    "periodStart": "2024-01-01",
    "periodEnd": "2024-01-31",
    "message": "Successfully generated 15 payrolls"
  }
}
```

---

### 6.2 Get All Payrolls (Paginated)
**Endpoint:** `GET /payrolls`

**Access:** ADMIN, SUPER_USER

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 15, max: 100): Items per page

**Example Request:**
```
GET /payrolls?page=1&per_page=15
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Payrolls retrieved successfully",
  "data": [
    {
      "id": "aa0e8400-e29b-41d4-a716-446655440006",
      "employeeName": "John Doe",
      "period": "2024-01 - 2024-01",
      "netSalary": 15800000,
      "status": "DRAFT",
      "generatedAt": "2024-01-30T10:30:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 15,
    "total": 30,
    "lastPage": 2,
    "firstPageUrl": "http://localhost:8080/api/v1/payrolls?page=1&per_page=15",
    "lastPageUrl": "http://localhost:8080/api/v1/payrolls?page=2&per_page=15",
    "nextPageUrl": "http://localhost:8080/api/v1/payrolls?page=2&per_page=15",
    "prevPageUrl": null
  }
}
```

---

### 6.3 Get Payroll by ID
**Endpoint:** `GET /payrolls/:id`

**Access:** All authenticated users

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Payroll fetched successfully",
  "data": {
    "id": "aa0e8400-e29b-41d4-a716-446655440006",
    "employeeId": "770e8400-e29b-41d4-a716-446655440002",
    "employeeName": "John Doe",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "John Doe",
    "periodStart": "2024-01-01",
    "periodEnd": "2024-01-31",
    "baseSalary": 15000000,
    "totalAllowance": 800000,
    "totalDeduction": 0,
    "netSalary": 15800000,
    "status": "DRAFT",
    "items": [
      {
        "id": "bb0e8400-e29b-41d4-a716-446655440007",
        "name": "Transport Allowance",
        "amount": 500000,
        "type": "EARNING"
      },
      {
        "id": "cc0e8400-e29b-41d4-a716-446655440008",
        "name": "Meal Allowance",
        "amount": 300000,
        "type": "EARNING"
      }
    ],
    "generatedAt": "2024-01-30T10:30:00Z",
    "createdAt": "2024-01-30T10:30:00Z",
    "updatedAt": "2024-01-30T10:30:00Z"
  }
}
```

---

### 6.4 Update Payroll Status
**Endpoint:** `PATCH /payrolls/:id/status`

**Access:** ADMIN, SUPER_USER

**Request Body:**
```json
{
  "status": "APPROVED"
}
```

**Field Validation:**
- `status`: Required, one of: `DRAFT`, `APPROVED`, `PAID`

**Status Transition Rules:**
- `DRAFT` → `APPROVED` ✓
- `DRAFT` → `DRAFT` ✓
- `APPROVED` → `PAID` ✓
- `APPROVED` → `APPROVED` ✓
- `PAID` → `PAID` ✓
- Other transitions: ✗ NOT allowed

**Success Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Status updated successfully",
  "data": null
}
```

---

### 6.5 Export Payroll to CSV
**Endpoint:** `GET /payrolls/export/csv`

**Access:** ADMIN, SUPER_USER

**Query Parameters:**
- `month` (required): Payroll month (1-12)
- `year` (required): Payroll year (2020-2100)

**Example Request:**
```
GET /payrolls/export/csv?month=1&year=2024
```

**Success Response (200 OK):**
```
Content-Type: text/csv
Content-Disposition: attachment; filename=payroll_export_20240130_103045.csv

Bank Name,Account Number,Account Holder,Amount,Description
BCA,1234567890,John Doe,15800000.00,Payroll 2024-01-01 - John Doe
Mandiri,0987654321,Jane Smith,16200000.00,Payroll 2024-01-01 - Jane Smith
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid/missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 422 | Validation Error - Invalid input data |
| 500 | Internal Server Error |

---

## 7. LEAVE MANAGEMENT MODULE

### 7.1 Create Leave Request
Submit a new leave request.

**Endpoint:** `POST /leave/requests`

**Authentication:** Required (All authenticated users)

**Request Body:**
```json
{
  "leaveTypeId": "uuid",
  "startDate": "2026-02-10",
  "endDate": "2026-02-12",
  "reason": "Family vacation",
  "attachmentUrl": "",
  "emergencyContact": "08123456789"
}
```

**Response:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Leave request created successfully",
  "data": {
    "id": "uuid",
    "employeeId": "uuid",
    "employeeName": "Test Employee",
    "employeePosition": "Software Engineer",
    "leaveType": {
      "id": "uuid",
      "name": "Cuti Tahunan",
      "code": "ANNUAL",
      "isPaid": true
    },
    "startDate": "2026-02-10",
    "endDate": "2026-02-12",
    "totalDays": 3,
    "reason": "Family vacation",
    "status": "PENDING"
  }
}
```

---

### 7.2 Get My Leave Requests
Get current user's leave requests with pagination.

**Endpoint:** `GET /leave/requests/my`

**Authentication:** Required

**Query Parameters:**
- `page` (integer, default: 1)
- `per_page` (integer, default: 15)

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave requests retrieved successfully",
  "data": [...],
  "pagination": {...}
}
```

---

### 7.3 Get Leave Request by ID
Get specific leave request details.

**Endpoint:** `GET /leave/requests/:id`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave request retrieved successfully",
  "data": {...}
}
```

---

### 7.4 Get My Leave Balances
Get current user's leave balances for a specific year.

**Endpoint:** `GET /leave/balances/my`

**Authentication:** Required

**Query Parameters:**
- `year` (integer, default: current year)

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave balances retrieved successfully",
  "data": {
    "employeeId": "uuid",
    "employeeName": "Test Employee",
    "year": 2026,
    "balances": [
      {
        "leaveTypeId": "uuid",
        "leaveTypeName": "Cuti Tahunan",
        "balance": 12,
        "used": 0,
        "pending": 0,
        "available": 12
      }
    ]
  }
}
```

---

### 7.5 Get Pending Leave Requests (Admin Only)
Get all pending leave requests.

**Endpoint:** `GET /leave/requests/pending`

**Authentication:** Required (Admin, Super User)

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Pending requests retrieved successfully",
  "data": [...]
}
```

---

### 7.6 Approve Leave Request (Admin Only)
Approve a pending leave request.

**Endpoint:** `PUT /leave/requests/:id/approve`

**Authentication:** Required (Admin, Super User)

**Request Body (Optional):**
```json
{
  "approvalNote": "Approved. Enjoy your vacation!"
}
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave request approved successfully",
  "data": {
    "status": "APPROVED",
    "approvedBy": "admin-uuid",
    "approvedAt": "2026-01-31T12:00:00Z"
  }
}
```

---

### 7.7 Reject Leave Request (Admin Only)
Reject a pending leave request.

**Endpoint:** `PUT /leave/requests/:id/reject`

**Authentication:** Required (Admin, Super User)

**Request Body:**
```json
{
  "rejectionReason": "Insufficient staff coverage for requested dates"
}
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave request rejected successfully",
  "data": {
    "status": "REJECTED",
    "rejectionReason": "Insufficient staff coverage"
  }
}
```

---

## 8. OVERTIME MANAGEMENT MODULE

### 8.1 Get Active Overtime Policies
Get list of all active overtime policies.

**Endpoint:** `GET /overtime/policies`

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Policies retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "Standard Overtime (1.5x)",
      "rateType": "MULTIPLIER",
      "rateMultiplier": 1.5,
      "maxOvertimeHoursPerDay": 4,
      "maxOvertimeHoursPerMonth": 40
    }
  ]
}
```

---

### 8.2 Create Overtime Request
Submit a new overtime request.

**Endpoint:** `POST /overtime/requests`

**Authentication:** Required

**Request Body:**
```json
{
  "overtimeDate": "2026-02-10",
  "startTime": "18:00",
  "endTime": "21:00",
  "reason": "Need to finish project milestone",
  "overtimePolicyId": "uuid"
}
```

**Response:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Overtime request created successfully",
  "data": {
    "id": "uuid",
    "employeeName": "Test Employee",
    "overtimeDate": "2026-02-10",
    "startTime": "18:00",
    "endTime": "21:00",
    "totalHours": 3,
    "status": "PENDING",
    "overtimePolicy": {
      "rateMultiplier": 1.5
    }
  }
}
```

---

### 8.3 Get My Overtime Requests
Get current user's overtime requests.

**Endpoint:** `GET /overtime/requests/my`

**Authentication:** Required

**Query Parameters:**
- `page` (integer, default: 1)
- `per_page` (integer, default: 15)

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime requests retrieved successfully",
  "data": [...],
  "pagination": {...}
}
```

---

### 8.4 Get Pending Overtime Requests (Admin Only)
Get all pending overtime requests.

**Endpoint:** `GET /overtime/requests/pending`

**Authentication:** Required (Admin, Super User)

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Pending requests retrieved successfully",
  "data": [...]
}
```

---

### 8.5 Approve Overtime Request (Admin Only)
Approve a pending overtime request.

**Endpoint:** `PUT /overtime/requests/:id/approve`

**Authentication:** Required (Admin, Super User)

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime request approved successfully",
  "data": {
    "status": "APPROVED",
    "approvedAt": "2026-01-31T12:00:00Z"
  }
}
```

---

### 8.6 Reject Overtime Request (Admin Only)
Reject a pending overtime request.

**Endpoint:** `PUT /overtime/requests/:id/reject`

**Authentication:** Required (Admin, Super User)

**Request Body:**
```json
{
  "rejectionReason": "Overtime not justified"
}
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime request rejected successfully",
  "data": {
    "status": "REJECTED"
  }
}
```

---

### 8.7 Clock In for Overtime
Clock in when starting overtime work.

**Endpoint:** `POST /overtime/requests/:id/clock-in`

**Authentication:** Required

**Request Body (Optional):**
```json
{
  "notes": "Starting overtime work"
}
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Clocked in successfully",
  "data": {
    "clockInTime": "2026-01-31T18:05:00Z",
    "clockOutTime": null
  }
}
```

---

### 8.8 Clock Out for Overtime
Clock out when finishing overtime work.

**Endpoint:** `POST /overtime/requests/:id/clock-out`

**Authentication:** Required

**Request Body (Optional):**
```json
{
  "notes": "Overtime work completed"
}
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Clocked out successfully",
  "data": {
    "clockInTime": "2026-01-31T18:05:00Z",
    "clockOutTime": "2026-01-31T21:00:00Z",
    "actualHours": 2.92
  }
}
```

---

### 8.9 Calculate Overtime Pay (Admin Only)
Calculate overtime pay for an employee within date range.

**Endpoint:** `GET /overtime/calculation/:employeeId`

**Authentication:** Required (Admin, Super User)

**Query Parameters:**
- `start_date` (string, required): YYYY-MM-DD format
- `end_date` (string, required): YYYY-MM-DD format

**Calculation:**
```
hourly_rate = salary_base / 173
overtime_pay = total_hours × hourly_rate × rate_multiplier
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime pay calculated successfully",
  "data": {
    "employeeId": "uuid",
    "employeeName": "Test Employee",
    "totalHours": 15,
    "rateMultiplier": 1.5,
    "hourlyRate": 46243,
    "overtimePay": 1040470
  }
}
```

---

## Common Headers

### Request Headers
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

### Response Headers
```
Content-Type: application/json
```

For file upload:
```
Content-Type: multipart/form-data
```

For CSV download:
```
Content-Type: text/csv
Content-Disposition: attachment; filename=payroll_export_YYYYMMDD_HHMMSS.csv
```
