CREATE TABLE "core_sessions" (
  "id" UUID PRIMARY KEY NOT NULL,
  "user_id" BIGINT NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "blocked" BOOLEAN NOT NULL DEFAULT FALSE,  
  "expires" TIMESTAMPTZ NOT NULL,
  "created" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  FOREIGN KEY (user_id) REFERENCES "public"."core_users" (id)
);

CREATE INDEX ON "public"."core_sessions" ("user_id");
COMMENT ON COLUMN "core_sessions"."id" IS 'auto generated session id';
COMMENT ON COLUMN "public"."core_sessions"."user_id" IS 'required,numeric,shouldref=core_users.id';
COMMENT ON COLUMN "core_sessions"."client_ip" IS 'required,ipaddress';
COMMENT ON COLUMN "core_sessions"."user_agent" IS 'required,user_agent';
COMMENT ON COLUMN "core_sessions"."refresh_token" IS 'required,token,desc=refresh token';
COMMENT ON COLUMN "core_sessions"."blocked" IS 'required,boolean,desc=refresh token blocked or not';
COMMENT ON COLUMN "core_sessions"."expires" IS 'timstamp when the refresh token gets expired';
COMMENT ON COLUMN "core_sessions"."created" IS 'timestamp when the refresh token was created';


