# Leave Management API - Postman Collection

## Base URL
```
http://localhost:8080/api/v1/leave
```

## Authentication
All endpoints require Bearer Token authentication. Include the token in the request header:
```
Authorization: Bearer {{access_token}}
```

---

## 📋 Leave Management Endpoints

### 1. Get All Leave Types
Get list of all available leave types.

**Endpoint:** `GET /leave/types` (Note: This endpoint may need to be created if not exists)

**Query Parameters:** None

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave types retrieved successfully",
  "data": [
    {
      "id": "e887251e-8a3c-4bb6-a314-8f90544f4e2f",
      "name": "Cuti Tahunan",
      "code": "ANNUAL",
      "description": "",
      "maxDaysPerYear": 12,
      "isPaid": true,
      "requiresApproval": true,
      "isActive": true
    },
    {
      "id": "f68b3c49-5123-46ce-a5d9-c655a1626a5d",
      "name": "Cuti Sakit",
      "code": "SICK",
      "maxDaysPerYear": 14,
      "isPaid": true,
      "requiresApproval": true,
      "isActive": true
    }
  ]
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/leave/types \
  -H "Authorization: Bearer {{access_token}}"
```

---

### 2. Create Leave Request
Submit a new leave request.

**Endpoint:** `POST /leave/requests`

**Request Body:**
```json
{
  "leaveTypeId": "e887251e-8a3c-4bb6-a314-8f90544f4e2f",
  "startDate": "2026-02-10",
  "endDate": "2026-02-12",
  "reason": "Family vacation",
  "attachmentUrl": "",
  "emergencyContact": "08123456789"
}
```

**Field Descriptions:**
- `leaveTypeId` (UUID, required): ID of the leave type from leave types
- `startDate` (string, required): Start date in format YYYY-MM-DD
- `endDate` (string, required): End date in format YYYY-MM-DD
- `reason` (string, required): Reason for leave (min 10 chars, max 1000 chars)
- `attachmentUrl` (string, optional): URL to supporting document
- `emergencyContact` (string, optional): Emergency contact phone number

**Response Example:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Leave request created successfully",
  "data": {
    "id": "e92b9b24-a867-49c3-bfca-3e38bf24511b",
    "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
    "employeeName": "Test Employee",
    "employeePosition": "Software Engineer",
    "leaveType": {
      "id": "e887251e-8a3c-4bb6-a314-8f90544f4e2f",
      "name": "Cuti Tahunan",
      "code": "ANNUAL",
      "isPaid": true
    },
    "startDate": "2026-02-10",
    "endDate": "2026-02-12",
    "totalDays": 3,
    "reason": "Family vacation",
    "attachmentUrl": "",
    "emergencyContact": "08123456789",
    "status": "PENDING",
    "approvedBy": null,
    "approvedByName": null,
    "approvedAt": null,
    "rejectionReason": "",
    "createdAt": "2026-01-31T11:47:57Z",
    "updatedAt": "2026-01-31T11:47:57Z"
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/leave/requests \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{access_token}}" \
  -d '{
    "leaveTypeId": "e887251e-8a3c-4bb6-a314-8f90544f4e2f",
    "startDate": "2026-02-10",
    "endDate": "2026-02-12",
    "reason": "Family vacation",
    "emergencyContact": "08123456789"
  }'
```

---

### 3. Get My Leave Requests
Get list of current user's leave requests with pagination.

**Endpoint:** `GET /leave/requests/my`

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 15)

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave requests retrieved successfully",
  "data": [
    {
      "id": "e92b9b24-a867-49c3-bfca-3e38bf24511b",
      "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
      "employeeName": "Test Employee",
      "employeePosition": "Software Engineer",
      "leaveType": {
        "id": "e887251e-8a3c-4bb6-a314-8f90544f4e2f",
        "name": "Cuti Tahunan",
        "code": "ANNUAL",
        "isPaid": true
      },
      "startDate": "2026-02-10",
      "endDate": "2026-02-12",
      "totalDays": 3,
      "reason": "Family vacation",
      "status": "PENDING",
      "approvedBy": null,
      "approvedAt": null
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 15,
    "total": 1,
    "lastPage": 1
  }
}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/leave/requests/my?page=1&per_page=15" \
  -H "Authorization: Bearer {{access_token}}"
```

---

### 4. Get Leave Request by ID
Get specific leave request details.

**Endpoint:** `GET /leave/requests/:id`

**Path Parameters:**
- `id` (UUID, required): Leave request ID

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave request retrieved successfully",
  "data": {
    "id": "e92b9b24-a867-49c3-bfca-3e38bf24511b",
    "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
    "employeeName": "Test Employee",
    "leaveType": {
      "id": "e887251e-8a3c-4bb6-a314-8f90544f4e2f",
      "name": "Cuti Tahunan",
      "code": "ANNUAL"
    },
    "status": "PENDING"
  }
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/leave/requests/e92b9b24-a867-49c3-bfca-3e38bf24511b \
  -H "Authorization: Bearer {{access_token}}"
```

---

### 5. Get My Leave Balances
Get current user's leave balances for a specific year.

**Endpoint:** `GET /leave/balances/my`

**Query Parameters:**
- `year` (integer, optional): Year to get balances for (default: current year)

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave balances retrieved successfully",
  "data": {
    "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
    "employeeName": "Test Employee",
    "year": 2026,
    "balances": [
      {
        "leaveTypeId": "e887251e-8a3c-4bb6-a314-8f90544f4e2f",
        "leaveTypeName": "Cuti Tahunan",
        "balance": 12,
        "used": 0,
        "pending": 0,
        "available": 12
      },
      {
        "leaveTypeId": "f68b3c49-5123-46ce-a5d9-c655a1626a5d",
        "leaveTypeName": "Cuti Sakit",
        "balance": 12,
        "used": 0,
        "pending": 0,
        "available": 12
      }
    ]
  }
}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/leave/balances/my?year=2026" \
  -H "Authorization: Bearer {{access_token}}"
```

---

### 6. Get Pending Leave Requests (Admin Only)
Get all pending leave requests that need approval.

**Endpoint:** `GET /leave/requests/pending`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Pending requests retrieved successfully",
  "data": [
    {
      "id": "e92b9b24-a867-49c3-bfca-3e38bf24511b",
      "employeeName": "Test Employee",
      "employeePosition": "Software Engineer",
      "leaveType": {
        "name": "Cuti Tahunan",
        "code": "ANNUAL"
      },
      "startDate": "2026-02-10",
      "endDate": "2026-02-12",
      "totalDays": 3,
      "reason": "Family vacation",
      "status": "PENDING"
    }
  ]
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/leave/requests/pending \
  -H "Authorization: Bearer {{admin_access_token}}"
```

---

### 7. Approve Leave Request (Admin Only)
Approve a pending leave request.

**Endpoint:** `PUT /leave/requests/:id/approve`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Path Parameters:**
- `id` (UUID, required): Leave request ID

**Request Body (Optional):**
```json
{
  "approvalNote": "Approved. Enjoy your vacation!"
}
```

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave request approved successfully",
  "data": {
    "id": "e92b9b24-a867-49c3-bfca-3e38bf24511b",
    "status": "APPROVED",
    "approvedBy": "admin-uuid",
    "approvedAt": "2026-01-31T12:00:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/leave/requests/e92b9b24-a867-49c3-bfca-3e38bf24511b/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{admin_access_token}}" \
  -d '{"approvalNote": "Approved"}'
```

---

### 8. Reject Leave Request (Admin Only)
Reject a pending leave request.

**Endpoint:** `PUT /leave/requests/:id/reject`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Path Parameters:**
- `id` (UUID, required): Leave request ID

**Request Body:**
```json
{
  "rejectionReason": "Insufficient staff coverage for requested dates"
}
```

**Field Descriptions:**
- `rejectionReason` (string, required): Reason for rejection (min 10 chars, max 500 chars)

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Leave request rejected successfully",
  "data": {
    "id": "e92b9b24-a867-49c3-bfca-3e38bf24511b",
    "status": "REJECTED",
    "rejectionReason": "Insufficient staff coverage for requested dates"
  }
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/leave/requests/e92b9b24-a867-49c3-bfca-3e38bf24511b/reject \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{admin_access_token}}" \
  -d '{
    "rejectionReason": "Insufficient staff coverage for requested dates"
  }'
```

---

## 🔧 Postman Environment Variables

Create these variables in your Postman environment:

```json
{
  "base_url": "http://localhost:8080",
  "api_version": "api/v1",
  "access_token": "YOUR_JWT_ACCESS_TOKEN",
  "admin_access_token": "YOUR_ADMIN_JWT_ACCESS_TOKEN",
  "employee_id": "YOUR_EMPLOYEE_ID",
  "leave_type_id": "YOUR_LEAVE_TYPE_ID"
}
```

---

## 📝 Common Use Cases

### Use Case 1: Employee Requests Annual Leave
1. **Login** to get access token
2. **GET /leave/balances/my** to check available balance
3. **POST /leave/requests** with annual leave type ID
4. **GET /leave/requests/my** to see request status

### Use Case 2: Manager Approves Leave Request
1. **Login** as admin/manager to get access token
2. **GET /leave/requests/pending** to see pending requests
3. **PUT /leave/requests/:id/approve** to approve or **PUT /leave/requests/:id/reject** to reject
4. **GET /leave/requests/:id** to verify updated status

### Use Case 3: Employee Checks Leave History
1. **Login** to get access token
2. **GET /leave/requests/my?page=1&per_page=20** to see all requests
3. **GET /leave/balances/my?year=2026** to see current balances

---

## ⚠️ Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "statusCode": 400,
  "message": "Invalid date range",
  "error": "error details"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "statusCode": 401,
  "message": "Missing authorization header"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "statusCode": 403,
  "message": "You don't have permission to access this resource"
}
```

### 404 Not Found
```json
{
  "success": false,
  "statusCode": 404,
  "message": "Leave request not found"
}
```

### 409 Conflict
```json
{
  "success": false,
  "statusCode": 409,
  "message": "Overlapping leave request exists"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "statusCode": 500,
  "message": "Failed to create leave request",
  "error": "detailed error message"
}
```

---

## 📚 Additional Notes

1. **Leave Status Flow:** PENDING → APPROVED/REJECTED/CANCELLED
2. **Balance Calculation:** Available = Balance - Used - Pending
3. **Automatic Attendance:** Approved leaves create LEAVE status in attendance system
4. **Date Validation:** Start date must be before or equal to end date
5. **Overlap Check:** Cannot have overlapping leave requests for same employee
6. **Working Days Calculation:** Total days includes weekends (can be enhanced)

---

## 🧪 Testing Checklist

- [ ] Create leave request with valid data
- [ ] Create leave request with invalid date range (should fail)
- [ ] Create leave request with insufficient balance (should fail)
- [ ] Create overlapping leave requests (should fail)
- [ ] Get my leave requests with pagination
- [ ] Get leave balances for current year
- [ ] Admin approve pending request
- [ ] Admin reject pending request
- [ ] Verify attendance record created after approval
- [ ] Verify leave balance updated after approval

---

**Document Version:** 1.0
**Last Updated:** January 31, 2026
**API Version:** v1
