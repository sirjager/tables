CREATE TABLE "core_tables" (
  "id" BIGSERIAL PRIMARY KEY NOT NULL,
  "user_id" BIGINT NOT NULL,
  "name" varchar(60) UNIQUE NOT NULL,
  "columns" VARCHAR NOT NULL,
  "created" timestamptz NOT NULL DEFAULT (now()),
  "updated" timestamptz NOT NULL DEFAULT (now()),
  FOREIGN KEY (user_id) REFERENCES "public"."core_users" (id)
);

CREATE INDEX ON "public"."core_tables" ("user_id");
COMMENT ON COLUMN "public"."core_tables"."id" IS 'numeric,server-side auto generated id';
COMMENT ON COLUMN "public"."core_tables"."user_id" IS 'required,numeric,shouldref=core_users.id';
COMMENT ON COLUMN "public"."core_tables"."name" IS 'required,alphanumeric,unique,min=3';
COMMENT ON COLUMN "public"."core_tables"."columns" IS 'required,jsonstring,desc=columns of the tables';
