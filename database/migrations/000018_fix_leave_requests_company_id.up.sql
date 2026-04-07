-- Migration: Fix existing leave requests without company_id
-- Set company_id based on employee's company

UPDATE leave_requests lr
SET company_id = e.company_id
FROM employees e
WHERE lr.employee_id = e.id
  AND lr.company_id = '00000000-0000-0000-0000-000000000001'::uuid;
