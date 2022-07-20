-- -------------------------- ADD ONE TO -> CORE_SESSIONS --------------------------

-- name: AddSession :one
INSERT INTO "public"."core_sessions" 
(sid,uid,client_ip,user_agent,refresh_token,is_blocked,expires_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;


-- -------------------------- GET ONE FROM <- CORE_SESSIONS --------------------------

-- name: GetSession :one
SELECT * FROM "public"."core_sessions" WHERE sid = $1 LIMIT 1;
