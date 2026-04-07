# Employee Module - New Fields Documentation

## Overview
This document describes the new fields added to the employee module as of February 10, 2026.

## New Fields Added

### 1. Employment Status (`employment_status`)
- **Type**: VARCHAR(20)
- **Default**: 'PROBATION'
- **Required**: Yes
- **Allowed Values**:
  - `PERMANENT` - Full-time permanent employee
  - `CONTRACT` - Contract-based employee
  - `PROBATION` - Employee under probation period
- **Description**: Indicates the employment contract type/status of the employee
- **Database Constraint**: CHECK constraint with enum values
- **Index**: Created on `employment_status` for faster filtering

### 2. Job Level (`job_level`)
- **Type**: VARCHAR(20)
- **Default**: 'STAFF'
- **Required**: Yes
- **Allowed Values**:
  - `CEO` - Chief Executive Officer
  - `MANAGER` - Managerial position
  - `SUPERVISOR` - Supervisor position
  - `STAFF` - Staff/Individual contributor
- **Description**: Defines the hierarchical level of the employee in the organization
- **Database Constraint**: CHECK constraint with enum values
- **Index**: Created on `job_level` for faster filtering

### 3. Gender (`gender`)
- **Type**: VARCHAR(10)
- **Default**: None (nullable)
- **Required**: No (but recommended)
- **Allowed Values**:
  - `MALE` - Male
  - `FEMALE` - Female
- **Description**: Gender identification of the employee
- **Database Constraint**: CHECK constraint with enum values
- **Index**: Created on `gender` for faster filtering

### 4. Division (`division`)
- **Type**: VARCHAR(100)
- **Default**: 'GENERAL'
- **Required**: Yes
- **Allowed Values**: Any division name (free text, but should reference divisions table)
- **Description**: Department or division where the employee belongs
- **Database Constraint**: None (free text)
- **Index**: Created on `division` for faster filtering
- **Reference Table**: `divisions` table contains standardized divisions

## Divisions Reference Table

A new `divisions` table has been created to standardize division names across the system.

### Divisions Table Structure
```sql
CREATE TABLE divisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Pre-seeded Divisions
1. **Information Technology** (IT) - IT and Software Development
2. **Human Resources** (HR) - HR and Recruitment
3. **Finance** (FIN) - Finance and Accounting
4. **Marketing** (MKT) - Marketing and Sales
5. **Operations** (OPS) - Operations and Logistics
6. **General** (GEN) - General Administration

## API Updates

### Create Employee Request
New required fields added:
```json
{
  "employmentStatus": "PERMANENT",  // Required: PERMANENT|CONTRACT|PROBATION
  "jobLevel": "STAFF",              // Required: CEO|MANAGER|SUPERVISOR|STAFF
  "gender": "MALE",                 // Required: MALE|FEMALE
  "division": "Information Technology" // Required: Division name
}
```

### Update Employee Request
New optional fields added:
```json
{
  "employmentStatus": "PERMANENT",  // Optional
  "jobLevel": "MANAGER",            // Optional
  "gender": "FEMALE",               // Optional
  "division": "Human Resources"     // Optional
}
```

### Employee Response
New fields included in response:
```json
{
  "employmentStatus": "PERMANENT",
  "jobLevel": "MANAGER",
  "gender": "MALE",
  "division": "Information Technology"
}
```

## Validation Rules

### Create Employee
- `employmentStatus`: Required, must be one of: PERMANENT, CONTRACT, PROBATION
- `jobLevel`: Required, must be one of: CEO, MANAGER, SUPERVISOR, STAFF
- `gender`: Required, must be one of: MALE, FEMALE
- `division`: Required, max 100 characters

### Update Employee
- `employmentStatus`: Optional, must be one of: PERMANENT, CONTRACT, PROBATION
- `jobLevel`: Optional, must be one of: CEO, MANAGER, SUPERVISOR, STAFF
- `gender`: Optional, must be one of: MALE, FEMALE
- `division`: Optional, max 100 characters

## Migration Details

### Migration File
- **File**: `database/migrations/000005_add_employee_fields.up.sql`
- **Date**: February 10, 2026
- **Rollback**: `database/migrations/000005_add_employee_fields.down.sql`

### Changes Made
1. Added `employment_status` column with CHECK constraint
2. Added `job_level` column with CHECK constraint
3. Added `gender` column with CHECK constraint
4. Added `division` column
5. Created indexes for all four new columns
6. Created `divisions` reference table
7. Seeded default divisions data
8. Updated existing employee records with default values

### Database Constraints
All enum fields use CHECK constraints to ensure data integrity:
```sql
CHECK (employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION'))
CHECK (job_level IN ('CEO', 'MANAGER', 'SUPERVISOR', 'STAFF'))
CHECK (gender IN ('MALE', 'FEMALE'))
```

## Usage Examples

### Creating a New Employee
```bash
POST /api/v1/employees
Authorization: Bearer <token>

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "Password123!",
  "position": "Senior Software Engineer",
  "phoneNumber": "+628123456789",
  "address": "Jl. Sudirman No. 1",
  "salaryBase": 18000000,
  "joinDate": "2024-01-15",
  "bankName": "BCA",
  "bankAccountNumber": "1234567890",
  "bankAccountHolder": "John Doe",
  "scheduleId": "uuid",
  "employmentStatus": "PERMANENT",
  "jobLevel": "STAFF",
  "gender": "MALE",
  "division": "Information Technology"
}
```

### Updating Employee Fields
```bash
PATCH /api/v1/employees/:id
Authorization: Bearer <token>

{
  "employmentStatus": "PERMANENT",
  "jobLevel": "MANAGER",
  "division": "Information Technology"
}
```

### Querying Employees by New Fields
```bash
# Get all permanent employees
GET /api/v1/employees?employment_status=PERMANENT

# Get all managers
GET /api/v1/employees?job_level=MANAGER

# Get all IT division employees
GET /api/v1/employees?division=Information%20Technology

# Filter by multiple fields
GET /api/v1/employees?job_level=STAFF&division=IT&gender=MALE
```

## Benefits

1. **Better Employee Categorization**: Employment status helps distinguish between permanent, contract, and probation employees
2. **Organizational Structure**: Job level provides clear hierarchy within the organization
3. **Diversity Tracking**: Gender field enables HR analytics and diversity tracking
4. **Department Management**: Division field allows better organization and reporting by departments
5. **Data Integrity**: CHECK constraints ensure only valid values are stored
6. **Query Performance**: Indexes on all new fields improve query performance
7. **Standardization**: Divisions reference table ensures consistent naming

## Backward Compatibility

- Existing employee records have been updated with default values:
  - `employment_status`: 'PERMANENT'
  - `job_level`: 'STAFF'
  - `division`: 'GENERAL'
  - `gender`: NULL (remains null until updated)
- All existing API endpoints continue to work
- New fields are optional in update operations

## Future Enhancements

Possible future improvements:
1. Add more job levels as the organization grows
2. Add foreign key constraint between employees.division and divisions.name
3. Add division head/manager assignment
4. Create department-based reporting features
5. Add organizational chart visualization
6. Implement division-based access control
