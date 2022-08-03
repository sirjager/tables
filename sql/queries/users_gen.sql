-- ------------------------------ GET MULTIPLE _USERS <== [_USERS] ------------------------------

-- name: GetAllUsers :many
SELECT * FROM "public"."_users";

-- name: GetSomeUsers :many
SELECT * FROM "public"."_users" LIMIT @limit_::int OFFSET @offset_::int;


-- ------------------------------ GET ONE _USERS <== _USER  ------------------------------

-- name: GetUser :one
SELECT * FROM "public"."_users" WHERE id = $1 LIMIT 1;

-- name: GetUserWhereEmail :one
SELECT * FROM "public"."_users" WHERE email = $1 LIMIT 1;

-- name: GetUserWhereUsername :one
SELECT * FROM "public"."_users" WHERE username = $1 LIMIT 1;


-- ------------------------------ ADD ONE _USERS <-> _USER  ------------------------------

-- name: CreateUser :one  
INSERT INTO "public"."_users" (email,username,password,fullname) VALUES ($1, $2, $3, $4) RETURNING *;


-- ------------------------------ REMOVE ONE _USERS -> nil  ------------------------------

-- name: DeleteUser :exec
DELETE FROM "public"."_users" WHERE id = $1;


-- ------------------------------ UPDATE ONE _USERS <-> _USERS  ------------------------------

-- name: UpdateUserFullName :one
UPDATE "public"."_users" SET fullname = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserUsername :one
UPDATE "public"."_users" SET username = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserPassword :one
UPDATE "public"."_users" SET password = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserVerified :one 
UPDATE "public"."_users" SET verified = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserBlocked :one
UPDATE "public"."_users" SET blocked = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserPublic :one
UPDATE "public"."_users" SET public = $1 WHERE id = $2 RETURNING *;

