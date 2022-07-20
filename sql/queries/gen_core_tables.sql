-- ------------------------------ GET MULTIPLE CORE_TABLES <== [CORE_TABLES] ------------------------------

-- name: ListCoreTables :many
SELECT * FROM "public"."core_tables";

-- name: ListCoreTablesWithLimit :many
SELECT * FROM "public"."core_tables" LIMIT @limit_::int;

-- name: ListCoreTablesWithLimitOffset :many
SELECT * FROM "public"."core_tables" LIMIT @limit_::int OFFSET @offset_::int;


-- --------------------- GET MULTIPLE CORE_TABLES OF CORE_USERS.UID <== [CORE_TABLES] ---------------------

-- name: ListCoreTablesWithUid :many
SELECT * FROM "public"."core_tables" WHERE uid = $1;

-- name: ListCoreTablesWithUidWithLimit :many
SELECT * FROM "public"."core_tables" WHERE uid = $1 LIMIT @limit_::int;

-- name: ListCoreTablesWithUidWithLimitOffset :many
SELECT * FROM "public"."core_tables" WHERE uid = $1 LIMIT @limit_::int OFFSET @offset_::int;



-- -------------------------- GET ONE CORE_TABLES <- CORE_TABLES --------------------------

-- name: GetCoreTableWithTid :one
SELECT * FROM "public"."core_tables" WHERE tid = $1 LIMIT 1;

-- name: GetCoreTableWithName :one
SELECT * FROM "public"."core_tables" WHERE tablename = $1 LIMIT 1;

-- name: GetCoreTableWithTidAndUid :one
SELECT * FROM "public"."core_tables" WHERE tid = $1 AND uid = $2 LIMIT 1;


-- -------------------------- ADD CORE_TABLES <-> CORE_TABLES --------------------------

-- name: AddCoreTable :one
INSERT INTO "public"."core_tables" (tablename,uid,columns) VALUES ($1, $2, $3) RETURNING *;


-- -------------------------- REMOVE CORE_TABLES <-> CORE_TABLES --------------------------

-- name: RemoveCoreTableWithUidAndTid :exec
DELETE FROM "public"."core_tables" WHERE uid = $1 AND tid = $2;

-- name: RemoveCoreTableWithUidAndName :exec
DELETE FROM "public"."core_tables" WHERE uid = $1 AND tablename = $2;
