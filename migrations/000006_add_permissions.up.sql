CREATE TABLE IF NOT EXISTS permissions(
	id bigserial PRIMARY KEY,
	code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions(
	user_id bigint NOT NULL REFERENCES users on DELETE CASCADE,
	permissions_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
	PRIMARY KEY (user_id, permissions_id)
);

-- adding two permission into table
INSERT INTO permissions (code)
VALUES
 ('movies:read'),
 ('movies:write')
