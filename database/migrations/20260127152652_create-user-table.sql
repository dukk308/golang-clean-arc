-- +goose Up
-- create "users" table
CREATE TABLE "public"."users" (
  "id" text NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp NULL,
  "created_by" text NULL,
  "updated_by" text NULL,
  "deleted_by" text NULL,
  "auth_provider" character varying(255) NULL DEFAULT 'local',
  "auth_provider_id" character varying(255) NULL,
  "username" character varying(255) NULL,
  "email" character varying(255) NULL,
  "password" character varying(255) NULL,
  "role" character varying(255) NULL DEFAULT 'viewer',
  PRIMARY KEY ("id")
);
-- create index "idx_provider_id" to table: "users"
CREATE INDEX "idx_provider_id" ON "public"."users" ("auth_provider", "auth_provider_id");
-- create index "idx_provider_username" to table: "users"
CREATE INDEX "idx_provider_username" ON "public"."users" ("auth_provider", "username");

-- +goose Down
-- reverse: create index "idx_provider_username" to table: "users"
DROP INDEX "public"."idx_provider_username";
-- reverse: create index "idx_provider_id" to table: "users"
DROP INDEX "public"."idx_provider_id";
-- reverse: create "users" table
DROP TABLE "public"."users";
