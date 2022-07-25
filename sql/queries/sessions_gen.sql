-- -------------------------- ADD ONE TO -> CORE_SESSIONS --------------------------

-- name: CreateSession :one
INSERT INTO "public"."core_sessions" 
(id,user_id,client_ip,user_agent,refresh_token,blocked,expires) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;


-- -------------------------- GET ONE FROM <- CORE_SESSIONS --------------------------

-- name: GetSession :one
SELECT * FROM "public"."core_sessions" WHERE id = $1 LIMIT 1;
