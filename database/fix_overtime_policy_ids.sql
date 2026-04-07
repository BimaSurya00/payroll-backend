-- ============================================
-- Fix Invalid Overtime Policy IDs
-- ============================================

-- 1. Check for invalid overtime_policy_id values
SELECT 
    id,
    employee_id,
    overtime_date,
    overtime_policy_id,
    status
FROM overtime_requests
WHERE overtime_policy_id IS NULL 
   OR overtime_policy_id = '' 
   OR LENGTH(overtime_policy_id) != 36
   OR overtime_policy_id ~ 'pending'
ORDER BY created_at DESC;

-- 2. Get the default/active overtime policy ID
SELECT id, name, is_active
FROM overtime_policies
WHERE is_active = true
ORDER BY created_at
LIMIT 1;

-- 3. Update invalid overtime_policy_id with the default policy
UPDATE overtime_requests
SET overtime_policy_id = (
    SELECT id 
    FROM overtime_policies 
    WHERE is_active = true 
    ORDER BY created_at 
    LIMIT 1
)
WHERE overtime_policy_id IS NULL 
   OR overtime_policy_id = '' 
   OR LENGTH(overtime_policy_id) != 36
   OR overtime_policy_id ~ 'pending';

-- 4. Verify the fix
SELECT 
    COUNT(*) as total_fixed,
    COUNT(CASE WHEN overtime_policy_id IS NOT NULL THEN 1 END) as with_policy,
    COUNT(CASE WHEN overtime_policy_id IS NULL THEN 1 END) as without_policy
FROM overtime_requests;

-- 5. Show sample of fixed records
SELECT 
    id,
    employee_id,
    overtime_date,
    overtime_policy_id,
    status
FROM overtime_requests
ORDER BY created_at DESC
LIMIT 5;
