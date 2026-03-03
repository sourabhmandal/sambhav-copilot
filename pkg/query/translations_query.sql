-- name: GetAllTranslationsByHashes :many
SELECT * FROM translations
WHERE company_id = $1 AND normalized_hash = ANY($2::text[]);

-- name: BulkInsertTranslations :many
INSERT INTO translations (
    company_id,
    normalized_hash,
    source_language,
    target_language,
    original_text,
    translated_text,
    confidence_score,
    provider
)
    SELECT
        unnest($1::bigint[]),
        unnest($2::text[]),
        unnest($3::text[]),
        unnest($4::text[]),
        unnest($5::text[]),
        unnest($6::text[]),
        unnest($7::numeric[]),
        unnest($8::text[])
ON CONFLICT (
    company_id,
    normalized_hash,
    source_language,
    target_language
)
DO NOTHING
RETURNING *;