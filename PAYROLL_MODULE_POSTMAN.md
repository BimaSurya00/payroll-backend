# Payroll Module API Documentation

## Base URL
```
http://localhost:8080/api/v1/payrolls
```

## Authentication
All endpoints require a valid Bearer token in the `Authorization` header:
```
Authorization: Bearer <your_access_token>
```

## Roles
- `USER`: Can view their own payroll details
- `ADMIN`: Can generate, view, and export payrolls
- `SUPER_USER`: Full access to all payroll operations

---

## Endpoints

### 1. Generate Bulk Payroll

Generate payroll for all active employees for a specific period.

**Endpoint:** `POST /payrolls/generate`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Request Body:**
```json
{
  "periodMonth": 1,
  "periodYear": 2024
}
```

**Field Descriptions:**
- `periodMonth` (required): Payroll month (1-12)
- `periodYear` (required): Payroll year (2020-2100)

**Behavior:**
- Fetches all active employees
- Calculates salary for each employee:
  - **Base Salary**: From employee record
  - **Allowances**: Transport (Rp 500.000) + Meal (Rp 300.000)
  - **Deductions**: Late days × Rp 50.000 (if any)
  - **Net Salary**: Base + Allowances - Deductions
- Creates payroll records with status `DRAFT`
- Returns summary of generated payrolls

**Response (201 Created):**
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

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "payroll for this period already exists",
  "error": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "no employees found",
  "error": null
}
```

---

### 2. Get All Payrolls

Retrieve paginated list of all payrolls.

**Endpoint:** `GET /payrolls`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `per_page` (optional, default: 15): Items per page (max: 100)

**Example Request:**
```
GET /payrolls?page=1&per_page=15
```

**Response (200 OK):**
```json
{
  "currentPage": 1,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "employeeName": "John Doe",
      "period": "2024-01 - 2024-01",
      "netSalary": 15800000,
      "status": "DRAFT",
      "generatedAt": "2024-01-28T10:30:00Z"
    }
  ],
  "firstPageUrl": "http://localhost:8080/api/v1/payrolls?page=1&per_page=15",
  "from": 1,
  "lastPage": 2,
  "lastPageUrl": "http://localhost:8080/api/v1/payrolls?page=2&per_page=15",
  "links": [
    {
      "url": null,
      "label": "pagination.previous",
      "active": false
    },
    {
      "url": "http://localhost:8080/api/v1/payrolls?page=1&per_page=15",
      "label": "1",
      "active": true
    }
  ],
  "nextPageUrl": "http://localhost:8080/api/v1/payrolls?page=2&per_page=15",
  "path": "http://localhost:8080/api/v1/payrolls",
  "perPage": 15,
  "prevPageUrl": null,
  "to": 15,
  "total": 30
}
```

---

### 3. Get Payroll by ID

Retrieve detailed payroll information including all items.

**Endpoint:** `GET /payrolls/:id`

**Roles Required:** All authenticated users (can view own payroll)

**URL Parameters:**
- `id` (required): Payroll ID (UUID)

**Example Request:**
```
GET /payrolls/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Payroll fetched successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "employeeId": "660e8400-e29b-41d4-a716-446655440000",
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
        "id": "770e8400-e29b-41d4-a716-446655440000",
        "name": "Transport Allowance",
        "amount": 500000,
        "type": "EARNING"
      },
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "name": "Meal Allowance",
        "amount": 300000,
        "type": "EARNING"
      }
    ],
    "generatedAt": "2024-01-28T10:30:00Z",
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
  "message": "payroll not found",
  "error": null
}
```

---

### 4. Update Payroll Status

Update payroll status (DRAFT → APPROVED → PAID).

**Endpoint:** `PATCH /payrolls/:id/status`

**Roles Required:** `ADMIN`, `SUPER_USER`

**URL Parameters:**
- `id` (required): Payroll ID (UUID)

**Request Body:**
```json
{
  "status": "APPROVED"
}
```

**Field Descriptions:**
- `status` (required): New status
  - `DRAFT`: Initial state after generation
  - `APPROVED`: Ready for payment
  - `PAID`: Payment completed

**Status Transition Rules:**
- `DRAFT` → `APPROVED` ✓
- `DRAFT` → `DRAFT` ✓
- `APPROVED` → `PAID` ✓
- `APPROVED` → `APPROVED` ✓
- `PAID` → `PAID` ✓
- Other transitions are **NOT allowed**

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Status updated successfully",
  "data": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "invalid status transition",
  "error": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "payroll not found",
  "error": null
}
```

---

### 5. Export Payroll to CSV

Export approved payrolls to CSV format for bank transfer.

**Endpoint:** `GET /payrolls/export/csv`

**Roles Required:** `ADMIN`, `SUPER_USER`

**Query Parameters:**
- `month` (required): Payroll month (1-12)
- `year` (required): Payroll year (2020-2100)

**Example Request:**
```
GET /payrolls/export/csv?month=1&year=2024
```

**Behavior:**
- Fetches all payrolls for the specified period
- Filters only `APPROVED` status payrolls
- Generates CSV with columns:
  - **Bank Name**: Employee's bank name
  - **Account Number**: Employee's bank account number
  - **Account Holder**: Employee's bank account holder name
  - **Amount**: Net salary to transfer
  - **Description**: Payroll period and employee name

**Response (200 OK):**
```
Content-Type: text/csv
Content-Disposition: attachment; filename=payroll_export_20240128_103045.csv

Bank Name,Account Number,Account Holder,Amount,Description
BCA,1234567890,John Doe,15800000.00,Payroll 2024-01-01 - John Doe
Mandiri,0987654321,Jane Smith,16200000.00,Payroll 2024-01-01 - Jane Smith
BNI,1122334456,Bob Johnson,15900000.00,Payroll 2024-01-01 - Bob Johnson
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "Invalid month parameter",
  "error": "Month must be between 1 and 12"
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "success": false,
  "statusCode": 500,
  "message": "Failed to export CSV",
  "error": "Database connection error"
}
```

**CSV Format Notes:**
- Filename format: `payroll_export_YYYYMMDD_HHMMSS.csv`
- Amount format: 2 decimal places
- Description format: `Payroll YYYY-MM-DD - Employee Name`
- Can be imported directly to internet banking systems

---

## Salary Calculation Logic

### Formula
```
Base Salary = From employee.salary_base
Allowances = Transport Allowance (Rp 500.000) + Meal Allowance (Rp 300.000)
Deductions = Late Days × Late Deduction Per Day (Rp 50.000)
Net Salary = Base Salary + Allowances - Deductions
```

### Example Calculation
```
Employee: John Doe
Base Salary: Rp 15.000.000
Late Days: 2 days

Transport Allowance: Rp 500.000
Meal Allowance: Rp 300.000
Late Deduction: 2 × Rp 50.000 = Rp 100.000

Net Salary = 15.000.000 + 500.000 + 300.000 - 100.000
Net Salary = Rp 15.700.000
```

### Allowance Constants
- **Transport Allowance**: Rp 500.000
- **Meal Allowance**: Rp 300.000
- **Late Deduction**: Rp 50.000 per day

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Payroll/Employee not found |
| 409 | Conflict - Payroll already exists for period |
| 422 | Validation Error - Invalid input data |
| 500 | Internal Server Error |

---

## Validation Rules

### Generate Bulk Payroll
- `periodMonth`: Required, integer between 1-12
- `periodYear`: Required, integer between 2020-2100

### Update Payroll Status
- `status`: Required, must be one of: DRAFT, APPROVED, PAID

### Export CSV
- `month`: Required, integer between 1-12
- `year`: Required, integer between 2020-2100

---

## Payroll Status Workflow

```
┌─────────┐     ┌──────────┐     ┌──────┐
│  DRAFT  │────▶│ APPROVED │────▶│ PAID │
└─────────┘     └──────────┘     └──────┘
     │               │
     └───────────────┘ (can stay in same status)
```

**Status Meanings:**
- **DRAFT**: Newly generated, not yet reviewed
- **APPROVED**: Reviewed and approved, ready for payment
- **PAID**: Payment completed

**Important Notes:**
- Can only export CSV for `APPROVED` payrolls
- Once `PAID`, status cannot be changed
- Status transitions are validated server-side

---

## Payroll Items Structure

Each payroll contains multiple items with two types:

### EARNING Type
- Transport Allowance
- Meal Allowance
- Other bonuses (if added in future)

### DEDUCTION Type
- Late Deduction
- Other deductions (if added in future)

**Example Items:**
```json
"items": [
  {
    "id": "uuid",
    "name": "Transport Allowance",
    "amount": 500000,
    "type": "EARNING"
  },
  {
    "id": "uuid",
    "name": "Meal Allowance",
    "amount": 300000,
    "type": "EARNING"
  },
  {
    "id": "uuid",
    "name": "Late Deduction",
    "amount": 100000,
    "type": "DEDUCTION"
  }
]
```

---

## Postman Collection JSON

Import this JSON into Postman:

```json
{
  "info": {
    "name": "Payroll Module API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api/v1/payrolls",
      "type": "string"
    },
    {
      "key": "accessToken",
      "value": "your_access_token_here",
      "type": "string"
    },
    {
      "key": "payrollId",
      "value": "payroll_id_here",
      "type": "string"
    }
  ],
  "item": [
    {
      "name": "1. Generate Bulk Payroll",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"periodMonth\": 1,\n  \"periodYear\": 2024\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/generate",
          "host": ["{{baseUrl}}"],
          "path": ["generate"]
        }
      }
    },
    {
      "name": "2. Get All Payrolls",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}?page=1&per_page=15",
          "host": ["{{baseUrl}}"],
          "query": [
            {
              "key": "page",
              "value": "1"
            },
            {
              "key": "per_page",
              "value": "15"
            }
          ]
        }
      }
    },
    {
      "name": "3. Get Payroll by ID",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/:id",
          "host": ["{{baseUrl}}"],
          "path": [":id"],
          "variable": [
            {
              "key": "id",
              "value": "{{payrollId}}"
            }
          ]
        }
      }
    },
    {
      "name": "4. Update Payroll Status to APPROVED",
      "request": {
        "method": "PATCH",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"status\": \"APPROVED\"\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/:id/status",
          "host": ["{{baseUrl}}"],
          "path": [":id", "status"],
          "variable": [
            {
              "key": "id",
              "value": "{{payrollId}}"
            }
          ]
        }
      }
    },
    {
      "name": "5. Update Payroll Status to PAID",
      "request": {
        "method": "PATCH",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"status\": \"PAID\"\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/:id/status",
          "host": ["{{baseUrl}}"],
          "path": [":id", "status"],
          "variable": [
            {
              "key": "id",
              "value": "{{payrollId}}"
            }
          ]
        }
      }
    },
    {
      "name": "6. Export Payroll to CSV",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/export/csv?month=1&year=2024",
          "host": ["{{baseUrl}}"],
          "path": ["export", "csv"],
          "query": [
            {
              "key": "month",
              "value": "1"
            },
            {
              "key": "year",
              "value": "2024"
            }
          ]
        }
      },
      "response": []
    }
  ]
}
```

---

## Test Flow in Postman

### Recommended Testing Sequence:

1. **Generate Payroll**
   - Call `POST /payrolls/generate` with periodMonth and periodYear
   - Verify all active employees get payroll records
   - Verify status is `DRAFT`
   - Save the payroll ID from response

2. **List All Payrolls**
   - Call `GET /payrolls` to see all generated payrolls
   - Verify pagination works correctly
   - Check employee names and net salaries

3. **View Payroll Details**
   - Call `GET /payrolls/:id` with saved payroll ID
   - Verify complete payroll information including:
     - Employee details (name, bank info)
     - Salary breakdown (base, allowances, deductions)
     - All payroll items

4. **Approve Payroll**
   - Call `PATCH /payrolls/:id/status` with status `APPROVED`
   - Verify status changed successfully
   - Attempt invalid status transition (should fail)

5. **Export to CSV**
   - Call `GET /payrolls/export/csv?month=X&year=YYYY`
   - Verify only APPROVED payrolls are included
   - Download and verify CSV format
   - Check bank information is correct

6. **Mark as Paid**
   - Call `PATCH /payrolls/:id/status` with status `PAID`
   - Verify status cannot be changed after PAID

7. **Test Duplicate Generation**
   - Try generating payroll for same period again
   - Verify 409 Conflict error

---

## Bank Transfer Integration

### CSV Format for Internet Banking

The exported CSV follows a generic format compatible with most Indonesian internet banking systems:

```csv
Bank Name,Account Number,Account Holder,Amount,Description
BCA,1234567890,John Doe,15800000.00,Payroll 2024-01-01 - John Doe
```

### Supported Banks

The system supports any bank. Common Indonesian banks:
- **BCA** (Bank Central Asia)
- **Mandiri**
- **BNI** (Bank Nasional Indonesia)
- **BRI** (Bank Rakyat Indonesia)
- **CIMB Niaga**
- **Danamon**
- **Permata**

### Import to Bank

Most internet banking systems support bulk transfer via CSV upload:

1. Login to your corporate internet banking
2. Navigate to "Bulk Transfer" or "Batch Transfer"
3. Download the CSV template from the bank
4. Match columns:
   - Bank Name → Bank Code/Name
   - Account Number → Beneficiary Account
   - Account Holder → Beneficiary Name
   - Amount → Transfer Amount
   - Description → Reference/Remark
5. Upload the exported CSV
6. Verify and confirm transfers

---

## Important Notes

1. **Period Uniqueness**: Payroll can only be generated once per period. Re-generating for the same period will fail with 409 Conflict.

2. **Employee Bank Info**: Ensure employees have complete bank information (bank_name, bank_account_number, bank_account_holder) before generating payroll for CSV export.

3. **Status Workflow**: Follow the proper status workflow: DRAFT → APPROVED → PAID. Only APPROVED payrolls can be exported to CSV.

4. **Calculation Constants**: Allowance and deduction amounts are currently hardcoded in the system:
   - Transport: Rp 500.000
   - Meal: Rp 300.000
   - Late deduction: Rp 50.000/day

5. **Transaction Safety**: Payroll generation uses database transactions to ensure data integrity. If any payroll creation fails, the entire batch fails.

6. **Late Days Integration**: Currently set to 0 (no late deductions). Future enhancement can integrate with attendance module to calculate actual late days.

7. **CSV Filename**: Generated filename includes timestamp: `payroll_export_YYYYMMDD_HHMMSS.csv`

8. **Currency**: All amounts are in Indonesian Rupiah (IDR).

---

## Future Enhancements

Potential improvements for the payroll module:

1. **Overtime Calculation**: Add overtime pay based on attendance records
2. **Custom Allowances**: Allow admin to configure custom allowances per employee
3. **Tax Calculation**: Integrate Indonesian tax calculation (PPh 21)
4. **BPJS Integration**: Add BPJS deductions (health, employment)
5. **Payslip PDF**: Generate individual payslip PDF for employees
6. **Notification**: Send email notification when payroll is generated/paid
7. **Audit Log**: Track who generated/approved/payroll and when
8. **Custom Periods**: Support weekly/bi-weekly payroll periods
9. **Mass Update**: Allow bulk status updates for multiple payrolls
10. **Reporting**: Generate payroll summary reports by department/period

---

## Troubleshooting

### Common Issues

**Issue**: Payroll generation returns "no employees found"
- **Solution**: Ensure at least one employee exists in the system with salary_base > 0

**Issue**: CSV export returns empty file
- **Solution**: Ensure payroll status is `APPROVED` before exporting

**Issue**: Status update fails with "invalid status transition"
- **Solution**: Follow the proper workflow: DRAFT → APPROVED → PAID

**Issue**: Bank info shows empty in CSV
- **Solution**: Update employee records with complete bank information

**Issue**: Cannot generate payroll for same period twice
- **Solution**: This is expected behavior. Delete existing payrolls first if needed (not yet implemented)

---

## Security Considerations

1. **Access Control**: Only admins can generate and export payrolls
2. **Data Privacy**: Payroll information is sensitive and should be accessed only by authorized personnel
3. **CSV Security**: Downloaded CSV files contain sensitive bank information - handle with care
4. **Audit Trail**: Consider adding audit logs for all payroll operations
5. **Encryption**: Consider encrypting salary data in database for additional security

---

## Performance Notes

1. **Bulk Generation**: Generates payroll for all employees in a single transaction
2. **Pagination**: List endpoint uses pagination to handle large datasets
3. **Database Indexes**: Ensure indexes exist on period_start, period_end columns for optimal query performance
4. **CSV Generation**: Efficient streaming for large datasets
