-- Ensure testuser exists and has permissions
DO $$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'testuser') THEN
            CREATE ROLE testuser WITH LOGIN PASSWORD 'testpass';
            ALTER ROLE testuser CREATEDB;
        END IF;
    END $$;

-- Ensure testuser has access to testdb
GRANT ALL PRIVILEGES ON DATABASE testdb TO testuser;
