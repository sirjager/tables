CREATE TABLE "core_sessions" (
  "sid" UUID PRIMARY KEY NOT NULL,
  "uid" BIGINT NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "is_blocked" BOOLEAN NOT NULL DEFAULT FALSE,  
  "expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  FOREIGN KEY (uid) REFERENCES "public"."core_users" (uid)
);

CREATE INDEX ON "public"."core_sessions" ("uid");
COMMENT ON COLUMN "core_sessions"."sid" IS 'auto generated session id';
COMMENT ON COLUMN "public"."core_sessions"."uid" IS 'required,numeric,shouldref=core_users.uid';
COMMENT ON COLUMN "core_sessions"."client_ip" IS 'required,ipaddress';
COMMENT ON COLUMN "core_sessions"."user_agent" IS 'required,user_agent';
COMMENT ON COLUMN "core_sessions"."refresh_token" IS 'required,token,desc=refresh token';
COMMENT ON COLUMN "core_sessions"."is_blocked" IS 'required,boolean,desc=refresh token blocked or not';
COMMENT ON COLUMN "core_sessions"."expires_at" IS 'timstamp when the refresh token gets expired';
COMMENT ON COLUMN "core_sessions"."created_at" IS 'timestamp when the refresh token was created';


