-- ------------------------------ GET MULTIPLE CORE_TABLES <== [CORE_TABLES] ------------------------------

-- name: GetAllTables :many
SELECT * FROM "public"."core_tables";

-- name: GetSomeTables :many
SELECT * FROM "public"."core_tables" LIMIT @limit_::int OFFSET @offset_::int;


-- --------------------- GET MULTIPLE CORE_TABLES OF CORE_USERS.user_id <== [CORE_TABLES] ---------------------

-- name: GetTablesWhereUser :many
SELECT * FROM "public"."core_tables" WHERE user_id = $1;

-- name: GetSomeTablesWhereUser :many
SELECT * FROM "public"."core_tables" WHERE user_id = $1 LIMIT @limit_::int OFFSET @offset_::int;


-- -------------------------- GET ONE CORE_TABLES <- CORE_TABLES --------------------------

-- name: GetTable :one
SELECT * FROM "public"."core_tables" WHERE id = $1 LIMIT 1;

-- name: GetTableWhereName :one
SELECT * FROM "public"."core_tables" WHERE name = $1 LIMIT 1;

-- name: GetTableWhereIDAndUser :one
SELECT * FROM "public"."core_tables" WHERE id = $1 AND user_id = $2 LIMIT 1;


-- -------------------------- ADD CORE_TABLES <-> CORE_TABLES --------------------------

-- name: CreateTable :one
INSERT INTO "public"."core_tables" (name,user_id,columns) VALUES ($1, $2, $3) RETURNING *;


-- -------------------------- REMOVE CORE_TABLES <-> CORE_TABLES --------------------------

-- name: DeleteTable :exec
DELETE FROM "public"."core_tables" WHERE id = $1;

-- name: DeleteTablesWhereUser :exec
DELETE FROM "public"."core_tables" WHERE user_id = $1;

-- name: DeleteTableWhereUserAndName :exec
DELETE FROM "public"."core_tables" WHERE user_id = $1 AND name = $2;


-- -------------------------- UPDATE CORE_TABLES <-> CORE_TABLES --------------------------

-- name: UpdateTableColumns :one
UPDATE "public"."core_tables" SET columns = $1 WHERE id = $2 RETURNING *;


