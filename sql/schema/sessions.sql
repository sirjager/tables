CREATE TABLE "_sessions" (
  "id" UUID PRIMARY KEY NOT NULL,
  "user_id" BIGINT NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "blocked" BOOLEAN NOT NULL DEFAULT FALSE,  
  "expires" TIMESTAMPTZ NOT NULL,
  "created" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  FOREIGN KEY (user_id) REFERENCES "public"."_users" (id) ON DELETE CASCADE
);

CREATE INDEX ON "public"."_sessions" ("user_id");
COMMENT ON COLUMN "public"."_sessions"."id" IS 'auto generated session id';
COMMENT ON COLUMN "public"."_sessions"."user_id" IS 'required,numeric,shouldref=_users.id';
COMMENT ON COLUMN "public"."_sessions"."client_ip" IS 'required,ipaddress';
COMMENT ON COLUMN "public"."_sessions"."user_agent" IS 'required,user_agent';
COMMENT ON COLUMN "public"."_sessions"."refresh_token" IS 'required,token,desc=refresh token';
COMMENT ON COLUMN "public"."_sessions"."blocked" IS 'required,boolean,desc=refresh token blocked or not';
COMMENT ON COLUMN "public"."_sessions"."expires" IS 'timstamp when the refresh token gets expired';
COMMENT ON COLUMN "public"."_sessions"."created" IS 'timestamp when the refresh token was created';


