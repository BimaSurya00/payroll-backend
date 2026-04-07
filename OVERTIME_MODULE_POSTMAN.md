# Overtime Management API - Postman Collection

## Base URL
```
http://localhost:8080/api/v1/overtime
```

## Authentication
All endpoints require Bearer Token authentication. Include the token in the request header:
```
Authorization: Bearer {{access_token}}
```

---

## 📋 Overtime Management Endpoints

### 1. Get Active Overtime Policies
Get list of all active overtime policies.

**Endpoint:** `GET /overtime/policies`

**Query Parameters:** None

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Policies retrieved successfully",
  "data": [
    {
      "id": "8ed09c24-e2a5-4a2d-af24-3c3a38e0bf95",
      "name": "Standard Overtime (1.5x)",
      "description": "",
      "rateType": "MULTIPLIER",
      "rateMultiplier": 1.5,
      "fixedAmount": 0,
      "minOvertimeMinutes": 60,
      "maxOvertimeHoursPerDay": 4,
      "maxOvertimeHoursPerMonth": 40,
      "requiresApproval": true,
      "isActive": true,
      "createdAt": "2026-01-31T09:33:19Z",
      "updatedAt": "2026-01-31T09:33:19Z"
    }
  ]
}
```

**Field Descriptions:**
- `rateType`: Either "FIXED" (fixed amount per hour) or "MULTIPLIER" (multiplier of hourly rate)
- `rateMultiplier`: Multiplier value when rateType is MULTIPLIER (e.g., 1.5x)
- `fixedAmount`: Fixed amount per hour when rateType is FIXED
- `minOvertimeMinutes`: Minimum overtime duration in minutes
- `maxOvertimeHoursPerDay`: Maximum overtime hours allowed per day
- `maxOvertimeHoursPerMonth`: Maximum overtime hours allowed per month
- `requiresApproval`: Whether overtime requests need approval

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/overtime/policies \
  -H "Authorization: Bearer {{access_token}}"
```

---

### 2. Create Overtime Request
Submit a new overtime request.

**Endpoint:** `POST /overtime/requests`

**Request Body:**
```json
{
  "overtimeDate": "2026-02-10",
  "startTime": "18:00",
  "endTime": "21:00",
  "reason": "Need to finish project milestone",
  "overtimePolicyId": "8ed09c24-e2a5-4a2d-af24-3c3a38e0bf95"
}
```

**Field Descriptions:**
- `overtimeDate` (string, required): Date of overtime in format YYYY-MM-DD
- `startTime` (string, required): Start time in format HH:MM (24-hour format)
- `endTime` (string, required): End time in format HH:MM (24-hour format)
- `reason` (string, required): Reason for overtime (min 10 chars, max 1000 chars)
- `overtimePolicyId` (UUID, required): ID of the overtime policy to apply

**Validation Rules:**
- Overtime hours must not exceed `maxOvertimeHoursPerDay` from policy
- Cannot create duplicate overtime request for same date
- Total hours calculated automatically: `(endTime - startTime) in hours`

**Response Example:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Overtime request created successfully",
  "data": {
    "id": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
    "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
    "employeeName": "Test Employee",
    "employeePosition": "Software Engineer",
    "overtimeDate": "2026-02-10",
    "startTime": "18:00",
    "endTime": "21:00",
    "totalHours": 3,
    "reason": "Need to finish project milestone",
    "overtimePolicy": {
      "id": "8ed09c24-e2a5-4a2d-af24-3c3a38e0bf95",
      "name": "Standard Overtime (1.5x)",
      "rateType": "MULTIPLIER",
      "rateMultiplier": 1.5,
      "fixedAmount": 0
    },
    "status": "PENDING",
    "approvedBy": null,
    "approvedByName": null,
    "approvedAt": null,
    "rejectionReason": "",
    "createdAt": "2026-01-31T19:00:26+07:00",
    "updatedAt": "2026-01-31T19:00:26+07:00"
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/overtime/requests \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{access_token}}" \
  -d '{
    "overtimeDate": "2026-02-10",
    "startTime": "18:00",
    "endTime": "21:00",
    "reason": "Need to finish project milestone",
    "overtimePolicyId": "8ed09c24-e2a5-4a2d-af24-3c3a38e0bf95"
  }'
```

---

### 3. Get My Overtime Requests
Get list of current user's overtime requests with pagination.

**Endpoint:** `GET /overtime/requests/my`

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 15)

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime requests retrieved successfully",
  "data": [
    {
      "id": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
      "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
      "employeeName": "Test Employee",
      "employeePosition": "Software Engineer",
      "overtimeDate": "2026-02-10",
      "startTime": "18:00",
      "endTime": "21:00",
      "totalHours": 3,
      "reason": "Need to finish project milestone",
      "overtimePolicy": {
        "id": "8ed09c24-e2a5-4a2d-af24-3c3a38e0bf95",
        "name": "Standard Overtime (1.5x)",
        "rateType": "MULTIPLIER",
        "rateMultiplier": 1.5
      },
      "status": "PENDING"
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
curl -X GET "http://localhost:8080/api/v1/overtime/requests/my?page=1&per_page=15" \
  -H "Authorization: Bearer {{access_token}}"
```

---

### 4. Get Overtime Request by ID
Get specific overtime request details.

**Endpoint:** `GET /overtime/requests/:id`

**Path Parameters:**
- `id` (UUID, required): Overtime request ID

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime request retrieved successfully",
  "data": {
    "id": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
    "employeeName": "Test Employee",
    "overtimeDate": "2026-02-10",
    "totalHours": 3,
    "status": "APPROVED"
  }
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/overtime/requests/8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1 \
  -H "Authorization: Bearer {{access_token}}"
```

---

### 5. Get Pending Overtime Requests (Admin Only)
Get all pending overtime requests that need approval.

**Endpoint:** `GET /overtime/requests/pending`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Pending requests retrieved successfully",
  "data": [
    {
      "id": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
      "employeeName": "Test Employee",
      "employeePosition": "Software Engineer",
      "overtimeDate": "2026-02-10",
      "startTime": "18:00",
      "endTime": "21:00",
      "totalHours": 3,
      "reason": "Need to finish project milestone",
      "overtimePolicy": {
        "name": "Standard Overtime (1.5x)",
        "rateType": "MULTIPLIER",
        "rateMultiplier": 1.5
      },
      "status": "PENDING"
    }
  ]
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/overtime/requests/pending \
  -H "Authorization: Bearer {{admin_access_token}}"
```

---

### 6. Approve Overtime Request (Admin Only)
Approve a pending overtime request.

**Endpoint:** `PUT /overtime/requests/:id/approve`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Path Parameters:**
- `id` (UUID, required): Overtime request ID

**Request Body (Optional):**
```json
{
  "approvalNote": "Approved. Complete the task on time."
}
```

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime request approved successfully",
  "data": {
    "id": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
    "status": "APPROVED",
    "approvedBy": "admin-uuid",
    "approvedAt": "2026-01-31T12:00:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/overtime/requests/8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{admin_access_token}}" \
  -d '{"approvalNote": "Approved"}'
```

---

### 7. Reject Overtime Request (Admin Only)
Reject a pending overtime request.

**Endpoint:** `PUT /overtime/requests/:id/reject`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Path Parameters:**
- `id` (UUID, required): Overtime request ID

**Request Body:**
```json
{
  "rejectionReason": "Overtime not justified. Please provide more details."
}
```

**Field Descriptions:**
- `rejectionReason` (string, required): Reason for rejection (min 10 chars, max 500 chars)

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime request rejected successfully",
  "data": {
    "id": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
    "status": "REJECTED",
    "rejectionReason": "Overtime not justified. Please provide more details."
  }
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/overtime/requests/8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1/reject \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{admin_access_token}}" \
  -d '{
    "rejectionReason": "Overtime not justified. Please provide more details."
  }'
```

---

### 8. Clock In for Overtime (Employee Only)
Clock in when starting overtime work. Request must be approved first.

**Endpoint:** `POST /overtime/requests/:id/clock-in`

**Path Parameters:**
- `id` (UUID, required): Overtime request ID

**Request Body (Optional):**
```json
{
  "notes": "Starting overtime work"
}
```

**Field Descriptions:**
- `notes` (string, optional): Notes when clocking in (max 500 chars)

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Clocked in successfully",
  "data": {
    "id": "attendance-uuid",
    "overtimeRequestId": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
    "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
    "employeeName": "Test Employee",
    "clockInTime": "2026-01-31T18:05:00Z",
    "clockOutTime": null,
    "actualHours": 0,
    "notes": "Starting overtime work",
    "createdAt": "2026-01-31T18:05:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/overtime/requests/8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1/clock-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{access_token}}" \
  -d '{"notes": "Starting overtime work"}'
```

---

### 9. Clock Out for Overtime (Employee Only)
Clock out when finishing overtime work. Actual hours will be calculated automatically.

**Endpoint:** `POST /overtime/requests/:id/clock-out`

**Path Parameters:**
- `id` (UUID, required): Overtime request ID

**Request Body (Optional):**
```json
{
  "notes": "Overtime work completed"
}
```

**Field Descriptions:**
- `notes` (string, optional): Notes when clocking out (max 500 chars)

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Clocked out successfully",
  "data": {
    "id": "attendance-uuid",
    "overtimeRequestId": "8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1",
    "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
    "employeeName": "Test Employee",
    "clockInTime": "2026-01-31T18:05:00Z",
    "clockOutTime": "2026-01-31T21:00:00Z",
    "actualHours": 2.92,
    "notes": "Overtime work completed",
    "createdAt": "2026-01-31T18:05:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/overtime/requests/8ceb6aae-4add-4b9e-a1a4-0f0a2413efd1/clock-out \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {{access_token}}" \
  -d '{"notes": "Overtime work completed"}'
```

---

### 10. Calculate Overtime Pay (Admin/Payroll Only)
Calculate overtime pay for an employee within a date range.

**Endpoint:** `GET /overtime/calculation/:employeeId`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Path Parameters:**
- `employeeId` (UUID, required): Employee ID to calculate for

**Query Parameters:**
- `start_date` (string, required): Start date in format YYYY-MM-DD
- `end_date` (string, required): End date in format YYYY-MM-DD

**Calculation Formula:**
```
hourly_rate = salary_base / 173
overtime_pay = total_hours × hourly_rate × rate_multiplier
```

**Response Example:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Overtime pay calculated successfully",
  "data": {
    "employeeId": "d6d48f00-dc0e-4dfb-97db-3bd596ea41ec",
    "employeeName": "Test Employee",
    "totalHours": 15,
    "rateType": "MULTIPLIER",
    "rateMultiplier": 1.5,
    "hourlyRate": 46243,
    "overtimePay": 1040470
  }
}
```

**Field Descriptions:**
- `totalHours`: Total approved overtime hours in date range
- `rateType`: Rate type from policy (MULTIPLIER or FIXED)
- `rateMultiplier`: Multiplier value used
- `hourlyRate`: Calculated hourly rate from salary base (rounded)
- `overtimePay`: Total overtime pay in IDR (rounded)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/overtime/calculation/d6d48f00-dc0e-4dfb-97db-3bd596ea41ec?start_date=2026-02-01&end_date=2026-02-28" \
  -H "Authorization: Bearer {{admin_access_token}}"
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
  "overtime_policy_id": "YOUR_OVERTIME_POLICY_ID",
  "overtime_request_id": "YOUR_OVERTIME_REQUEST_ID"
}
```

---

## 📝 Common Use Cases

### Use Case 1: Employee Requests Overtime
1. **Login** to get access token
2. **GET /overtime/policies** to see available overtime policies
3. **POST /overtime/requests** with overtime date, times, and reason
4. **GET /overtime/requests/my** to see request status

### Use Case 2: Manager Approves Overtime Request
1. **Login** as admin/manager to get access token
2. **GET /overtime/requests/pending** to see pending requests
3. **PUT /overtime/requests/:id/approve** to approve or **PUT /overtime/requests/:id/reject** to reject
4. **GET /overtime/requests/:id** to verify updated status

### Use Case 3: Employee Performs Overtime Work
1. **Wait for approval** of overtime request
2. **POST /overtime/requests/:id/clock-in** when starting work
3. **POST /overtime/requests/:id/clock-out** when finishing work
4. Actual hours calculated automatically from clock-in to clock-out

### Use Case 4: Payroll Calculates Overtime Pay
1. **Login** as admin/payroll to get access token
2. **GET /overtime/calculation/:employeeId** with date range
3. Use calculated `overtimePay` for payroll processing

---

## ⚠️ Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "statusCode": 400,
  "message": "Invalid time range",
  "error": "error details"
}
```

Common causes:
- End time is before start time
- Overtime hours exceed maximum allowed
- Invalid date/time format

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
  "message": "Overtime request not found"
}
```

### 409 Conflict
```json
{
  "success": false,
  "statusCode": 409,
  "message": "Overtime request for this date already exists"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "statusCode": 500,
  "message": "Failed to create overtime request",
  "error": "detailed error message"
}
```

---

## 📚 Additional Notes

### Overtime Status Flow
```
PENDING → APPROVED → [Clock In] → [Clock Out] → COMPLETED
PENDING → REJECTED
```

### Overtime Pay Calculation
- **Hourly Rate**: Calculated as `salary_base / 173` (assuming 173 working hours/month)
- **Multiplier Type**: `overtime_pay = hours × hourly_rate × rate_multiplier`
- **Fixed Type**: `overtime_pay = hours × fixed_amount`

### Time Format
- All times in **24-hour format** (HH:MM)
- Example: 18:00 for 6 PM, 21:00 for 9 PM
- Timezone: Server timezone (WIB/UTC+7)

### Attendance Integration
- Clock in/out records actual hours worked
- If no attendance record, uses requested hours
- Actual hours calculated from clock-in to clock-out duration

### Policy Limits
- Cannot request overtime exceeding `maxOvertimeHoursPerDay`
- System checks total monthly overtime against `maxOvertimeHoursPerMonth`
- Minimum overtime duration: `minOvertimeMinutes`

---

## 🧪 Testing Checklist

- [ ] Create overtime request with valid data
- [ ] Create overtime request with invalid time range (should fail)
- [ ] Create duplicate overtime request for same date (should fail)
- [ ] Create overtime exceeding daily limit (should fail)
- [ ] Get my overtime requests with pagination
- [ ] Get active overtime policies
- [ ] Admin approve pending request
- [ ] Admin reject pending request
- [ ] Clock in for approved overtime request
- [ ] Clock out and verify actual hours calculated
- [ ] Calculate overtime pay for date range
- [ ] Verify payroll integration with overtime data

---

## 🔗 Integration with Other Modules

### Payroll Module
Overtime pay calculated can be included in payroll:
```json
{
  "overtimePay": 1040470,
  "description": "Overtime pay for Feb 2026 (15 hours × 1.5x)"
}
```

### Attendance Module
Overtime attendance records are separate but follow similar pattern:
- Regular attendance: 9 AM - 5 PM
- Overtime attendance: 6 PM - 9 PM (after regular hours)

### Employee Module
Employee salary base used for hourly rate calculation:
```json
{
  "salaryBase": 8000000,
  "hourlyRate": 46243
}
```

---

**Document Version:** 1.0
**Last Updated:** January 31, 2026
**API Version:** v1
