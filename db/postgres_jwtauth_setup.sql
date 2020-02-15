CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE SCHEMA IF NOT EXISTS basic_auth;

CREATE TABLE IF NOT EXISTS basic_auth.users (
    id       SERIAL PRIMARY KEY,
    email    TEXT NOT NULL CHECK (email ~* '[A-Z0-9._%-]+@[A-Z0-9._%-]+\.[A-Z]{2,4}'),
    pass     TEXT NOT NULL CHECK (LENGTH(pass) < 512),
    role     NAME NOT NULL CHECK (LENGTH(role) < 512),
    verified BOOLEAN NOT NULL default false
);
-- function to check if inserted user role exists
CREATE OR REPLACE FUNCTION basic_auth.check_role_exists() RETURNS trigger AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles AS r WHERE r.rolname = NEW.role) THEN 
        RAISE FOREIGN_KEY_VIOLATION USING MESSAGE = 'unknown database role: ' || NEW.role;
        RETURN NULL;
    END IF;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;
-- trigger to call role check function
DROP TRIGGER IF EXISTS ensure_user_role_exists ON basic_auth.users;
CREATE TRIGGER ensure_user_role_exists
    AFTER INSERT OR UPDATE ON basic_auth.users
    FOR EACH ROW EXECUTE PROCEDURE basic_auth.check_role_exists();

-- function to encrypt inserted or updated passwords
CREATE OR REPLACE FUNCTION basic_auth.encrypt_pass() RETURNS TRIGGER AS $$
BEGIN
    IF tg_op = 'INSERT' OR NEW.pass <> OLD.pass THEN
        NEW.pass = crypt(NEW.pass, gen_salt('bf'));
    END IF;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;
-- trigger to call password encryption function
DROP TRIGGER IF EXISTS encrypt_pass on basic_auth.users;
CREATE TRIGGER encrypt_pass
    BEFORE INSERT OR UPDATE ON basic_auth.users
    FOR EACH ROW EXECUTE PROCEDURE basic_auth.encrypt_pass();

-- function to call from application to verify submitted username/password
CREATE OR REPLACE FUNCTION basic_auth.user_role(email text, pass text) RETURNS name AS $$
BEGIN
    RETURN (
        SELECT  role 
        FROM    basic_auth.users
        WHERE   users.email = user_role.email
        AND     users.pass = crypt(user_role.pass, users.pass)
        );
END
$$ LANGUAGE plpgsql;

-- **DANGEROUS** - command to truncate all tables in a schema
DO
$$
BEGIN
   RAISE NOTICE '%', 
   -- EXECUTE
   (SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
    FROM   pg_class
    WHERE  relkind = 'r' --only tables
    AND    relnamespace = 'schema_name'::regnamespace
   );
END
$$;