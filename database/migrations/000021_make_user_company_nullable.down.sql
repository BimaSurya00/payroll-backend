ALTER TABLE users DROP CONSTRAINT IF EXISTS users_company_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_unique;
ALTER TABLE users ALTER COLUMN company_id SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_company_id_fkey FOREIGN KEY (company_id) REFERENCES companies(id);
ALTER TABLE users ADD CONSTRAINT users_company_id_email_key UNIQUE (company_id, email);
UPDATE users SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
