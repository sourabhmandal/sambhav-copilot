-- Create "companies" table
CREATE TABLE "companies" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "created_at" timestamp NULL DEFAULT now(),
  PRIMARY KEY ("id")
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
