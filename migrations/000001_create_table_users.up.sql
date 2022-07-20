CREATE TABLE "core_users" (
  "uid" BIGSERIAL PRIMARY KEY NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(60) UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "fullname" VARCHAR(255) NOT NULL,
  "is_public" BOOLEAN NOT NULL DEFAULT FALSE,
  "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
  "is_blocked" BOOLEAN NOT NULL DEFAULT FALSE,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);


COMMENT ON COLUMN "core_users"."uid" IS 'server-side auto generated id';
COMMENT ON COLUMN "core_users"."email" IS 'required,max=255';
COMMENT ON COLUMN "core_users"."username" IS 'required,alphanumeric,min=3,max=60';
COMMENT ON COLUMN "core_users"."password" IS 'required,hashed';
COMMENT ON COLUMN "core_users"."fullname" IS 'optional,max=255';
COMMENT ON COLUMN "public"."core_users"."is_verified" IS 'optinal,default=false,desc=email is verified or not';
COMMENT ON COLUMN "public"."core_users"."is_public" IS 'optinal,default=false,desc=profile visible by others or not';
COMMENT ON COLUMN "public"."core_users"."is_blocked" IS 'optinal,default=false,desc=profile is accessible or not';
COMMENT ON COLUMN "core_users"."updated_at" IS 'server-side auto generated timestamp with time zone';
COMMENT ON COLUMN "core_users"."created_at" IS 'server-side auto generated timestamp with time zone';

