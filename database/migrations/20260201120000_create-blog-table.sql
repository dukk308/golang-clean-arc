-- +goose Up
CREATE TABLE "public"."blogs" (
  "id" text NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp NULL,
  "created_by" text NULL,
  "updated_by" text NULL,
  "deleted_by" text NULL,
  "title" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "content" text NULL,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "uni_blogs_slug" ON "public"."blogs" ("slug");

-- +goose Down
DROP INDEX "public"."uni_blogs_slug";
DROP TABLE "public"."blogs";
