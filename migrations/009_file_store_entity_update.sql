-- Migration: 009_file_store_entity_update
-- Generated: 2026-03-07T13:48:25+05:30

-- ====================================
-- UP Migration
-- ====================================

-- public."FileStoreMetadata" definition


CREATE TABLE public."FileStoreMetadata" (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	is_active bool NULL DEFAULT true,
	created_by_id int8 NULL,
	updated_by_id int8 NULL,
	channel_id int8 NULL,
	file_name text NULL,
	file_size int8 NULL,
	file_type text NULL,
	file_url text NULL,
	file_key text NULL,
	file_bucket text NULL,
	entity_id int8 NULL,
	entity_type text NULL,
	CONSTRAINT "FileStoreMetadata_pkey" PRIMARY KEY (id)
);

-- public."EntityDocuments" definition

-- Drop table

-- DROP TABLE public."EntityDocuments";

CREATE TABLE public."EntityDocuments" (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	is_active bool NULL DEFAULT true,
	created_by_id int8 NULL,
	updated_by_id int8 NULL,
	channel_id int8 NULL,
	"type" text NULL,
	document_type text NULL,
	entity_name text NULL,
	entity_id int8 NULL,
	description text NULL,
	CONSTRAINT "EntityDocuments_pkey" PRIMARY KEY (id)
);


-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- TODO: Add rollback statements manually
