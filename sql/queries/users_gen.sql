-- ------------------------------ GET MULTIPLE CORE_USERS <== [CORE_USERS] ------------------------------

-- name: ListCoreUsers :many
SELECT * FROM "public"."core_users";

-- name: ListCoreUsersWithLimit :many
SELECT * FROM "public"."core_users" LIMIT @limit_::int;

-- name: ListCoreUsersWithLimitOffset :many
SELECT * FROM "public"."core_users" LIMIT @limit_::int OFFSET @offset_::int;


-- ------------------------------ GET ONE CORE_USERS <== CORE_USER  ------------------------------

-- name: GetCoreUserWithUid :one
SELECT * FROM "public"."core_users" WHERE id = $1 LIMIT 1;

-- name: GetCoreUserWithEmail :one
SELECT * FROM "public"."core_users" WHERE email = $1 LIMIT 1;

-- name: GetCoreUserWithUsername :one
SELECT * FROM "public"."core_users" WHERE username = $1 LIMIT 1;


-- ------------------------------ ADD ONE CORE_USERS <-> CORE_USER  ------------------------------

-- name: AddCoreUser :one  
INSERT INTO "public"."core_users" (email,username,password,fullname) VALUES ($1, $2, $3, $4) RETURNING *;


-- ------------------------------ REMOVE ONE CORE_USERS -> nil  ------------------------------

-- name: RemoveCoreUserWithUid :exec
DELETE FROM "public"."core_users" WHERE id = $1;

-- name: RemoveCoreUserWithEmail :exec
DELETE FROM "public"."core_users" WHERE email = $1;

-- name: RemoveCoreUserWithUsername :exec
DELETE FROM "public"."core_users" WHERE username = $1;


-- ------------------------------ UPDATE ONE CORE_USERS <-> CORE_USERS  ------------------------------

-- name: UpdateCoreUserName :one
UPDATE "public"."core_users" SET fullname = $1 WHERE id = $2 RETURNING *;

-- name: UpdateCoreUserUsername :one
UPDATE "public"."core_users" SET username = $1 WHERE id = $2 RETURNING *;

-- name: UpdateCoreUserPassword :one
UPDATE "public"."core_users" SET password = $1 WHERE id = $2 RETURNING *;

-- name: UpdateCoreUserVerified :one 
UPDATE "public"."core_users" SET verified = $1 WHERE id = $2 RETURNING *;

-- name: UpdateCoreUserBlocked :one
UPDATE "public"."core_users" SET blocked = $1 WHERE id = $2 RETURNING *;

-- name: UpdateCoreUserPublic :one
UPDATE "public"."core_users" SET public = $1 WHERE id = $2 RETURNING *;

