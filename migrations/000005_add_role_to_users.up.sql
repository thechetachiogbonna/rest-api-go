BEGIN;

CREATE TYPE enum_role AS ENUM (
	'admin',
	'user'
);

ALTER TABLE users ADD COLUMN "role" enum_role NOT NULL DEFAULT 'user';

COMMIT;