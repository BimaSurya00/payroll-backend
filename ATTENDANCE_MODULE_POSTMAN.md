# Attendance Module API Documentation

## Base URL
```
http://localhost:8080/api/v1/attendance
```

## Authentication
All endpoints require a valid Bearer token in the `Authorization` header:
```
Authorization: Bearer <your_access_token>
```

## Roles
- `USER`, `ADMIN`, `SUPER_USER`: All authenticated users can clock in/out and view their own history

## Features
- **GPS Validation**: Clock in/out requires being within the allowed radius of the office
- **Late Detection**: Automatically marks attendance as LATE if clocked in after allowed time
- **Daily Records**: One attendance record per employee per day
- **Distance Tracking**: Records distance from office for each clock in/out

---

## Endpoints

### 1. Clock In

Record attendance clock-in with GPS location validation.

**Endpoint:** `POST /attendance/clock-in`

**Roles Required:** `USER`, `ADMIN`, `SUPER_USER`

**Request Body:**
```json
{
  "lat": -6.200100,
  "long": 106.816700
}
```

**Field Descriptions:**
- `lat` (required): Current latitude coordinate (-90 to 90)
- `long` (required): Current longitude coordinate (-180 to 180)

**How It Works:**
1. System retrieves the employee's assigned schedule
2. Calculates distance between current location and office location
3. Validates if within allowed radius (defined in schedule)
4. Checks if already clocked in today
5. Determines status (PRESENT or LATE) based on schedule time
6. Creates attendance record

**Status Determination Logic:**
- `PRESENT`: Clocked in before or at (schedule.timeIn + allowedLateMinutes)
- `LATE`: Clocked in after (schedule.timeIn + allowedLateMinutes)

**Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Clock in successful",
  "data": {
    "attendanceId": "770e8400-e29b-41d4-a716-446655440000",
    "employeeId": "550e8400-e29b-41d4-a716-446655440000",
    "clockInTime": "2024-01-29T01:30:00Z",
    "status": "PRESENT",
    "distance": 45.2,
    "scheduleName": "Regular Office Hour"
  }
}
```

**Response Examples by Status:**

**On Time:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Clock in successful",
  "data": {
    "attendanceId": "770e8400-e29b-41d4-a716-446655440000",
    "employeeId": "550e8400-e29b-41d4-a716-446655440000",
    "clockInTime": "2024-01-29T01:00:00Z",
    "status": "PRESENT",
    "distance": 12.5,
    "scheduleName": "Regular Office Hour"
  }
}
```

**Late:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Clock in successful",
  "data": {
    "attendanceId": "770e8400-e29b-41d4-a716-446655440000",
    "employeeId": "550e8400-e29b-41d4-a716-446655440000",
    "clockInTime": "2024-01-29T01:20:00Z",
    "status": "LATE",
    "distance": 25.8,
    "scheduleName": "Regular Office Hour"
  }
}
```

**Error Responses:**

**Already Clocked In (409 Conflict):**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "already clocked in today",
  "error": null
}
```

**Out of Office Range (400 Bad Request):**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "out of office range",
  "error": "You are 150 meters away from the office. Maximum allowed distance is 50 meters."
}
```

**Employee Not Found (404 Not Found):**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "employee not found",
  "error": null
}
```

**Schedule Not Found (404 Not Found):**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "schedule not found",
  "error": "Employee is not assigned to any schedule"
}
```

**Validation Error (422):**
```json
{
  "success": false,
  "statusCode": 422,
  "message": "Validation failed",
  "error": [
    {
      "field": "lat",
      "message": "lat must be greater than or equal to -90"
    }
  ]
}
```

---

### 2. Clock Out

Record attendance clock-out with GPS location.

**Endpoint:** `POST /attendance/clock-out`

**Roles Required:** `USER`, `ADMIN`, `SUPER_USER`

**Request Body:**
```json
{
  "lat": -6.200100,
  "long": 106.816700
}
```

**Field Descriptions:**
- `lat` (required): Current latitude coordinate (-90 to 90)
- `long` (required): Current longitude coordinate (-180 to 180)

**How It Works:**
1. System retrieves today's attendance record
2. Validates that employee has already clocked in
3. Updates attendance with clock-out time and location

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Clock out successful",
  "data": {
    "attendanceId": "770e8400-e29b-41d4-a716-446655440000",
    "clockOutTime": "2024-01-29T09:30:00Z",
    "distance": 18.3
  }
}
```

**Error Responses:**

**Not Clocked In (400 Bad Request):**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "not clocked in yet",
  "error": "You must clock in before clocking out"
}
```

**Employee Not Found (404 Not Found):**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "employee not found",
  "error": null
}
```

---

### 3. Get Attendance History

Retrieve attendance history for the authenticated user.

**Endpoint:** `GET /attendance/history`

**Roles Required:** `USER`, `ADMIN`, `SUPER_USER`

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 15, max: 100): Items per page

**Example Request:**
```
GET /attendance/history?page=1&per_page=15
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "History retrieved successfully",
  "data": {
    "currentPage": 1,
    "data": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440000",
        "date": "2024-01-29T00:00:00Z",
        "clockInTime": "2024-01-29T01:00:00Z",
        "clockOutTime": "2024-01-29T09:30:00Z",
        "status": "PRESENT",
        "notes": "",
        "scheduleName": "Regular Office Hour"
      },
      {
        "id": "770e8400-e29b-41d4-a716-446655440001",
        "date": "2024-01-28T00:00:00Z",
        "clockInTime": "2024-01-28T01:05:00Z",
        "clockOutTime": "2024-01-28T09:25:00Z",
        "status": "LATE",
        "notes": "",
        "scheduleName": "Regular Office Hour"
      },
      {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "date": "2024-01-27T00:00:00Z",
        "clockInTime": "2024-01-27T00:58:00Z",
        "clockOutTime": null,
        "status": "PRESENT",
        "notes": "",
        "scheduleName": "Regular Office Hour"
      }
    ],
    "firstPageUrl": "http://localhost:8080/api/v1/attendance/history?page=1&per_page=15",
    "from": 1,
    "lastPage": 1,
    "lastPageUrl": "http://localhost:8080/api/v1/attendance/history?page=1&per_page=15",
    "links": [
      {
        "url": null,
        "label": "&laquo; Previous",
        "active": false
      },
      {
        "url": "http://localhost:8080/api/v1/attendance/history?page=1&per_page=15",
        "label": "1",
        "active": true
      },
      {
        "url": null,
        "label": "Next &raquo;",
        "active": false
      }
    ],
    "nextPageUrl": null,
    "path": "http://localhost:8080/api/v1/attendance/history",
    "perPage": 15,
    "prevPageUrl": null,
    "to": 3,
    "total": 3
  }
}
```

**Field Descriptions:**
- `id`: Attendance record ID
- `date`: Date of attendance (UTC)
- `clockInTime`: Clock-in timestamp (null if not clocked in)
- `clockOutTime`: Clock-out timestamp (null if not clocked out)
- `status`: Attendance status (`PRESENT`, `LATE`, `ABSENT`, `LEAVE`)
- `notes`: Additional notes about the attendance
- `scheduleName`: Name of the schedule used for this attendance

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

## GPS Calculation Details

### Distance Calculation
The system uses the **Haversine Formula** to calculate the great-circle distance between two points on Earth:

```
Distance = 6,371,000 meters × acos(sin(lat1) × sin(lat2) + cos(lat1) × cos(lat2) × cos(long2 - long1))
```

### Example Distances
- **1 degree latitude**: ~111 km (69 miles)
- **1 degree longitude**: ~111 km (69 miles) at equator, varies by latitude
- **Office radius validation**: User must be within `allowedRadiusMeters` of office coordinates

### Jakarta Office Example
```json
{
  "officeLat": -6.200000,
  "officeLong": 106.816666,
  "allowedRadiusMeters": 50
}
```

**Valid Clock In (within 50m):**
```json
{
  "lat": -6.200050,
  "long": 106.816700,
  "distance": 8.5  // meters
}
```

**Invalid Clock In (outside 50m):**
```json
{
  "lat": -6.201000,
  "long": 106.817000,
  "distance": 125.3  // meters - REJECTED
}
```

---

## Attendance Workflow

### Complete Daily Attendance Flow

```
1. Employee arrives at office
   ↓
2. Clock In POST /attendance/clock-in
   - GPS validation (within 50m of office)
   - Time validation (before 09:15 for Regular Office Hour)
   - Status: PRESENT or LATE
   ↓
3. Works throughout the day
   ↓
4. Clock Out POST /attendance/clock-out
   - Records clock-out time and location
   ↓
5. View History GET /attendance/history
   - See all attendance records
```

### Example Timeline (Regular Office Hour: 09:00-17:00, Late Tolerance: 15 min)

| Time | Action | Status | Notes |
|------|--------|--------|-------|
| 08:55 | Clock In | PRESENT | 5 min early |
| 09:00 | Clock In | PRESENT | On time |
| 09:14 | Clock In | PRESENT | Within tolerance |
| 09:15 | Clock In | PRESENT | Exactly at tolerance limit |
| 09:16 | Clock In | LATE | 1 minute late |
| 17:00 | Clock Out | - | End of shift |
| 17:30 | Clock Out | - | 30 min after shift |

---

## Common Error Responses

### 401 Unauthorized
```json
{
  "success": false,
  "statusCode": 401,
  "message": "Missing or malformed JWT",
  "error": "Unauthorized"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "statusCode": 403,
  "message": "Insufficient permissions",
  "error": null
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "statusCode": 500,
  "message": "Failed to process request",
  "error": "Database connection error"
}
```

---

## Usage Examples with cURL

### Clock In
```bash
curl -X POST http://localhost:8080/api/v1/attendance/clock-in \
  -H "Authorization: Bearer <your_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "lat": -6.200100,
    "long": 106.816700
  }'
```

### Clock Out
```bash
curl -X POST http://localhost:8080/api/v1/attendance/clock-out \
  -H "Authorization: Bearer <your_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "lat": -6.200100,
    "long": 106.816700
  }'
```

### Get History
```bash
curl -X GET "http://localhost:8080/api/v1/attendance/history?page=1&per_page=15" \
  -H "Authorization: Bearer <your_access_token>"
```

---

## Testing Scenarios

### Scenario 1: Successful Clock In (On Time)
**Setup:** Schedule: 09:00-17:00, Late tolerance: 15 min

**Request (08:55 UTC):**
```bash
curl -X POST http://localhost:8080/api/v1/attendance/clock-in \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"lat": -6.200000, "long": 106.816666}'
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Clock in successful",
  "data": {
    "attendanceId": "...",
    "employeeId": "...",
    "clockInTime": "2024-01-29T01:55:00Z",
    "status": "PRESENT",
    "distance": 0,
    "scheduleName": "Regular Office Hour"
  }
}
```

### Scenario 2: Late Clock In
**Request (09:20 UTC - 20 min late):**
```bash
curl -X POST http://localhost:8080/api/v1/attendance/clock-in \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"lat": -6.200000, "long": 106.816666}'
```

**Expected Response:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Clock in successful",
  "data": {
    "attendanceId": "...",
    "employeeId": "...",
    "clockInTime": "2024-01-29T02:20:00Z",
    "status": "LATE",
    "distance": 0,
    "scheduleName": "Regular Office Hour"
  }
}
```

### Scenario 3: Out of Office Range
**Request (from home, 5km away):**
```bash
curl -X POST http://localhost:8080/api/v1/attendance/clock-in \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"lat": -6.250000, "long": 106.866666}'
```

**Expected Response:**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "out of office range",
  "error": null
}
```

### Scenario 4: Already Clocked In
**First Request (09:00 UTC):** → Success
**Second Request (09:05 UTC):**
```bash
curl -X POST http://localhost:8080/api/v1/attendance/clock-in \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"lat": -6.200000, "long": 106.816666}'
```

**Expected Response:**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "already clocked in today",
  "error": null
}
```

---

## Notes

1. **Timezone Handling**: All times are stored and compared in UTC

2. **Date Calculation**: "Today" is determined based on UTC date, not local timezone

3. **One Record Per Day**: Each employee can only have one attendance record per day

4. **GPS Accuracy**: 
   - GPS coordinates are accurate to approximately 5-10 meters
   - Consider this when setting `allowedRadiusMeters`
   - Recommended minimum: 50 meters for office environments

5. **Distance Calculation**: Uses Haversine formula for accurate great-circle distance

6. **Status Persistence**: Once status is set (PRESENT/LATE), it cannot be changed for that day

7. **Employee-Schedule Relationship**: 
   - Employee must be assigned to a schedule to clock in/out
   - If no schedule assigned, clock in will fail with "schedule not found"

8. **Clock Out Without Clock In**: Attempting to clock out without clocking in first will fail

9. **Pagination**: History endpoint supports pagination for viewing large attendance records

10. **Distance Tracking**: Both clock in and clock out locations are recorded for audit purposes

---

## Related Modules
- **Schedule Module**: Defines work hours and office locations for GPS validation
- **Employee Module**: Employees are assigned to schedules and have attendance records
- **Auth Module**: Provides JWT tokens for authentication
