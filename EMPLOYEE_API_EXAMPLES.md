# Employee API Examples - With New Fields

## Create Employee Examples

### Example 1: Permanent IT Staff
```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ahmad Rizky",
    "email": "ahmad.rizky@example.com",
    "password": "Password123!",
    "position": "Software Engineer",
    "phoneNumber": "+6281234567801",
    "address": "Jl. Sudirman No. 1, Jakarta Selatan",
    "salaryBase": 15000000,
    "joinDate": "2024-01-15",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "Ahmad Rizky",
    "scheduleId": "schedule-uuid-here",
    "employmentStatus": "PERMANENT",
    "jobLevel": "STAFF",
    "gender": "MALE",
    "division": "Information Technology"
  }'
```

### Example 2: Contract HR Manager
```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Siti Rahayu",
    "email": "siti.rahayu@example.com",
    "password": "Password123!",
    "position": "HR Manager",
    "phoneNumber": "+6281234567802",
    "address": "Jl. Thamrin No. 2, Jakarta Pusat",
    "salaryBase": 20000000,
    "joinDate": "2024-02-01",
    "bankName": "Mandiri",
    "bankAccountNumber": "0987654321",
    "bankAccountHolder": "Siti Rahayu",
    "scheduleId": "schedule-uuid-here",
    "employmentStatus": "CONTRACT",
    "jobLevel": "MANAGER",
    "gender": "FEMALE",
    "division": "Human Resources"
  }'
```

### Example 3: Probation Finance Staff
```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Budi Santoso",
    "email": "budi.santoso@example.com",
    "password": "Password123!",
    "position": "Junior Accountant",
    "phoneNumber": "+6281234567803",
    "address": "Jl. Gatot Subroto No. 3, Jakarta Selatan",
    "salaryBase": 8000000,
    "joinDate": "2024-02-10",
    "bankName": "BNI",
    "bankAccountNumber": "1111222233",
    "bankAccountHolder": "Budi Santoso",
    "scheduleId": "schedule-uuid-here",
    "employmentStatus": "PROBATION",
    "jobLevel": "STAFF",
    "gender": "MALE",
    "division": "Finance"
  }'
```

### Example 4: CEO (Permanent)
```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hendra Wijaya",
    "email": "hendra.wijaya@example.com",
    "password": "Password123!",
    "position": "Chief Executive Officer",
    "phoneNumber": "+6281234567804",
    "address": "Jl. Sudirman Kav 50, Jakarta Pusat",
    "salaryBase": 50000000,
    "joinDate": "2024-01-01",
    "bankName": "BCA",
    "bankAccountNumber": "999988887777",
    "bankAccountHolder": "Hendra Wijaya",
    "scheduleId": "schedule-uuid-here",
    "employmentStatus": "PERMANENT",
    "jobLevel": "CEO",
    "gender": "MALE",
    "division": "General"
  }'
```

## Update Employee Examples

### Update Job Level and Division
```bash
curl -X PATCH http://localhost:8080/api/v1/employees/<employee_id> \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "jobLevel": "MANAGER",
    "division": "Information Technology",
    "employmentStatus": "PERMANENT"
  }'
```

### Update Only Employment Status
```bash
curl -X PATCH http://localhost:8080/api/v1/employees/<employee_id> \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "employmentStatus": "PERMANENT"
  }'
```

## Query Examples

### Get All Permanent Employees
```bash
curl -X GET "http://localhost:8080/api/v1/employees?employment_status=PERMANENT" \
  -H "Authorization: Bearer <your_token>"
```

### Get All Managers
```bash
curl -X GET "http://localhost:8080/api/v1/employees?job_level=MANAGER" \
  -H "Authorization: Bearer <your_token>"
```

### Get All IT Division Employees
```bash
curl -X GET "http://localhost:8080/api/v1/employees?division=Information%20Technology" \
  -H "Authorization: Bearer <your_token>"
```

### Get All Female Employees
```bash
curl -X GET "http://localhost:8080/api/v1/employees?gender=FEMALE" \
  -H "Authorization: Bearer <your_token>"
```

### Combined Filters
```bash
curl -X GET "http://localhost:8080/api/v1/employees?job_level=STAFF&division=IT&gender=MALE" \
  -H "Authorization: Bearer <your_token>"
```

## Expected Response Format

### Success Response (201 Created)
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Employee created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "userId": "660e8400-e29b-41d4-a716-446655440000",
    "userName": "Ahmad Rizky",
    "userEmail": "ahmad.rizky@example.com",
    "position": "Software Engineer",
    "phoneNumber": "+6281234567801",
    "address": "Jl. Sudirman No. 1, Jakarta Selatan",
    "salaryBase": 15000000,
    "joinDate": "2024-01-15",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "Ahmad Rizky",
    "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
    "schedule": {
      "id": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
      "name": "Regular Office Hour",
      "timeIn": "09:00",
      "timeOut": "17:00",
      "allowedLateMinutes": 15
    },
    "employmentStatus": "PERMANENT",
    "jobLevel": "STAFF",
    "gender": "MALE",
    "division": "Information Technology",
    "createdAt": "2026-02-10T10:30:00Z",
    "updatedAt": "2026-02-10T10:30:00Z"
  }
}
```

### Update Response (200 OK)
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Employee updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "userId": "660e8400-e29b-41d4-a716-446655440000",
    "userName": "Ahmad Rizky",
    "userEmail": "ahmad.rizky@example.com",
    "position": "Software Engineer",
    "phoneNumber": "+6281234567801",
    "address": "Jl. Sudirman No. 1, Jakarta Selatan",
    "salaryBase": 15000000,
    "joinDate": "2024-01-15",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "Ahmad Rizky",
    "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
    "employmentStatus": "PERMANENT",
    "jobLevel": "MANAGER",
    "gender": "MALE",
    "division": "Information Technology",
    "createdAt": "2026-02-10T10:30:00Z",
    "updatedAt": "2026-02-10T11:00:00Z"
  }
}
```

### List Response (200 OK)
```json
{
  "currentPage": 1,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "userId": "660e8400-e29b-41d4-a716-446655440000",
      "userName": "Ahmad Rizky",
      "userEmail": "ahmad.rizky@example.com",
      "position": "Software Engineer",
      "phoneNumber": "+6281234567801",
      "address": "Jl. Sudirman No. 1, Jakarta Selatan",
      "salaryBase": 15000000,
      "joinDate": "2024-01-15",
      "bankName": "BCA",
      "bankAccountNumber": "1234567890",
      "bankAccountHolder": "Ahmad Rizky",
      "scheduleId": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
      "schedule": {
        "id": "8918b101-f313-4ef8-bbcb-c72ebfae3527",
        "name": "Regular Office Hour",
        "timeIn": "09:00",
        "timeOut": "17:00",
        "allowedLateMinutes": 15
      },
      "employmentStatus": "PERMANENT",
      "jobLevel": "STAFF",
      "gender": "MALE",
      "division": "Information Technology",
      "createdAt": "2026-02-10T10:30:00Z",
      "updatedAt": "2026-02-10T10:30:00Z"
    }
  ],
  "firstPageUrl": "http://localhost:8080/api/v1/employees?page=1&per_page=15",
  "from": 1,
  "lastPage": 5,
  "lastPageUrl": "http://localhost:8080/api/v1/employees?page=5&per_page=15",
  "links": [...],
  "nextPageUrl": "http://localhost:8080/api/v1/employees?page=2&per_page=15",
  "path": "http://localhost:8080/api/v1/employees",
  "perPage": 15,
  "prevPageUrl": null,
  "to": 15,
  "total": 67
}
```

## Error Responses

### Validation Error (422)
```json
{
  "success": false,
  "statusCode": 422,
  "message": "Validation failed",
  "error": [
    {
      "field": "employmentStatus",
      "message": "employment status must be one of: PERMANENT, CONTRACT, PROBATION"
    },
    {
      "field": "jobLevel",
      "message": "job level must be one of: CEO, MANAGER, SUPERVISOR, STAFF"
    }
  ]
}
```

## Testing Checklist

- [ ] Create employee with all required fields
- [ ] Create employee with each employment status (PERMANENT, CONTRACT, PROBATION)
- [ ] Create employee with each job level (CEO, MANAGER, SUPERVISOR, STAFF)
- [ ] Create employee with each gender (MALE, FEMALE)
- [ ] Create employee with different divisions
- [ ] Update each new field individually
- [ ] Update multiple new fields at once
- [ ] Query employees by employment status
- [ ] Query employees by job level
- [ ] Query employees by gender
- [ ] Query employees by division
- [ ] Test validation for invalid enum values
- [ ] Test backward compatibility with existing employees
