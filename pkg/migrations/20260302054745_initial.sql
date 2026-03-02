-- Drop index "idx_translations_lookup" from table: "translations"
DROP INDEX "idx_translations_lookup";
-- Modify "translations" table
ALTER TABLE "translations" DROP COLUMN "source_id", ADD COLUMN "normalized_hash" character(64) NOT NULL, ADD COLUMN "original_text" text NOT NULL, ADD CONSTRAINT "translations_normalized_hash_key" UNIQUE ("normalized_hash"), ADD CONSTRAINT "translations_normalized_hash_language_code_key" UNIQUE ("normalized_hash", "language_code");
-- Create index "idx_translations_lookup" to table: "translations"
CREATE INDEX "idx_translations_lookup" ON "translations" ("normalized_hash", "language_code");
-- Drop "source_strings" table
DROP TABLE "source_strings";
