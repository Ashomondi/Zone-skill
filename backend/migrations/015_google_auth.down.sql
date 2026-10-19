-- Manual rollback for Google sign-in support.
-- password_hash is intentionally left nullable: Google-only users may exist
-- and re-adding NOT NULL could fail.

DROP INDEX IF EXISTS users_google_id_unique;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS auth_provider;
ALTER TABLE users DROP COLUMN IF EXISTS google_id;
