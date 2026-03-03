DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'translation_input') THEN
        CREATE TYPE translation_input AS (
            company_id bigint,
            normalized_hash text,
            source_language text,
            target_language text,
            original_text text,
            translated_text text,
            confidence_score numeric,
            provider text
        );
    END IF;
END$$;
