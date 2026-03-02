-- name: GetTranslationByHash :one
SELECT * FROM translations
WHERE id = $1 LIMIT 1;
