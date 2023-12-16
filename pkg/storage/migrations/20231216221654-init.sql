-- +migrate Up
CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS app.resource (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	description text NOT NULL,
	reference text NOT NULL,
	level varchar(50) NOT NULL,
	type varchar(50) NOT NULL,
	CONSTRAINT resource_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_resource_level ON app.resource USING btree (level);
CREATE INDEX IF NOT EXISTS idx_resource_type ON app.resource USING btree (type);

CREATE TABLE IF NOT EXISTS app.tag (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	name text NOT NULL,
	CONSTRAINT tag_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_tag_name ON app.tag USING btree (name);


CREATE TABLE IF NOT EXISTS app.resources_tags (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	resource_id uuid NOT NULL,
	tag_id uuid NOT NULL,

	CONSTRAINT resources_tags_pkey PRIMARY KEY (id),
	CONSTRAINT fk_resources_tags_resource FOREIGN KEY (resource_id) REFERENCES app.resource(id),
	CONSTRAINT fk_resources_tags_tag FOREIGN KEY (tag_id) REFERENCES app.tag(id)
);

-- +migrate Down
DROP SCHEMA IF EXISTS app cascade;
