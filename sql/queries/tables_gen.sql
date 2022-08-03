-- ------------------------------ GET MULTIPLE _TABLES <== [_TABLES] ------------------------------

-- name: GetAllTables :many
SELECT * FROM "public"."_tables";

-- name: GetSomeTables :many
SELECT * FROM "public"."_tables" LIMIT @limit_::int OFFSET @offset_::int;


-- --------------------- GET MULTIPLE _TABLES OF _USERS.user_id <== [_TABLES] ---------------------

-- name: GetTablesWhereUser :many
SELECT * FROM "public"."_tables" WHERE user_id = $1;

-- name: GetSomeTablesWhereUser :many
SELECT * FROM "public"."_tables" WHERE user_id = $1 LIMIT @limit_::int OFFSET @offset_::int;


-- -------------------------- GET ONE _TABLES <- _TABLES --------------------------

-- name: GetTableWhereName :one
SELECT * FROM "public"."_tables" WHERE name = $1 LIMIT 1;

-- name: GetTableByUserIdAndTableName :one
SELECT * FROM "public"."_tables" WHERE user_id = $1 AND name = $2 LIMIT 1;


-- -------------------------- ADD _TABLES <-> _TABLES --------------------------

-- name: CreateTable :one
INSERT INTO "public"."_tables" (name,user_id,columns) VALUES ($1, $2, $3) RETURNING *;


-- -------------------------- REMOVE _TABLES <-> _TABLES --------------------------

-- name: DeleteTablesWhereUser :exec
DELETE FROM "public"."_tables" WHERE user_id = $1;

-- name: DeleteTableWhereName :exec
DELETE FROM "public"."_tables" WHERE name = $1;


-- -------------------------- UPDATE _TABLES <-> _TABLES --------------------------

-- name: UpdateTableColumns :one
UPDATE "public"."_tables" SET columns = $1 WHERE id = $2 RETURNING *;


