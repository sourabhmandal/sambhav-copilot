-- Create "source_strings" table
CREATE TABLE "source_strings" (
  "id" serial NOT NULL,
  "normalized_hash" character(64) NOT NULL,
  "original_text" text NOT NULL,
  "created_at" timestamp NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "source_strings_normalized_hash_key" UNIQUE ("normalized_hash")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "bio" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
-- Create "translations" table
CREATE TABLE "translations" (
  "id" serial NOT NULL,
  "source_id" integer NOT NULL,
  "language_code" character varying(10) NOT NULL,
  "translated_text" text NOT NULL,
  "confidence_score" numeric(4,3) NULL,
  "provider" character varying(50) NULL,
  "created_at" timestamp NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "translations_source_id_language_code_key" UNIQUE ("source_id", "language_code"),
  CONSTRAINT "translations_source_id_fkey" FOREIGN KEY ("source_id") REFERENCES "source_strings" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_translations_lookup" to table: "translations"
CREATE INDEX "idx_translations_lookup" ON "translations" ("source_id", "language_code");
