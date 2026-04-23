UPDATE users SET company_id = NULL WHERE role = 'SUPER_USER' AND company_id IS NOT NULL;
