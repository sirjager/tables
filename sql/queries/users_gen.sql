-- ------------------------------ GET MULTIPLE CORE_USERS <== [CORE_USERS] ------------------------------

-- name: GetAllUsers :many
SELECT * FROM "public"."core_users";

-- name: GetSomeUsers :many
SELECT * FROM "public"."core_users" LIMIT @limit_::int OFFSET @offset_::int;


-- ------------------------------ GET ONE CORE_USERS <== CORE_USER  ------------------------------

-- name: GetUser :one
SELECT * FROM "public"."core_users" WHERE id = $1 LIMIT 1;

-- name: GetUserWhereEmail :one
SELECT * FROM "public"."core_users" WHERE email = $1 LIMIT 1;

-- name: GetUserWhereUsername :one
SELECT * FROM "public"."core_users" WHERE username = $1 LIMIT 1;


-- ------------------------------ ADD ONE CORE_USERS <-> CORE_USER  ------------------------------

-- name: CreateUser :one  
INSERT INTO "public"."core_users" (email,username,password,fullname) VALUES ($1, $2, $3, $4) RETURNING *;


-- ------------------------------ REMOVE ONE CORE_USERS -> nil  ------------------------------

-- name: DeleteUser :exec
DELETE FROM "public"."core_users" WHERE id = $1;


-- ------------------------------ UPDATE ONE CORE_USERS <-> CORE_USERS  ------------------------------

-- name: UpdateUserFullName :one
UPDATE "public"."core_users" SET fullname = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserUsername :one
UPDATE "public"."core_users" SET username = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserPassword :one
UPDATE "public"."core_users" SET password = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserVerified :one 
UPDATE "public"."core_users" SET verified = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserBlocked :one
UPDATE "public"."core_users" SET blocked = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserPublic :one
UPDATE "public"."core_users" SET public = $1 WHERE id = $2 RETURNING *;

