CREATE TABLE "_tables" (
  "id" BIGSERIAL PRIMARY KEY NOT NULL,
  "user_id" BIGINT NOT NULL,
  "name" varchar(60) UNIQUE NOT NULL,
  "columns" VARCHAR NOT NULL,
  "created" timestamptz NOT NULL DEFAULT (now()),
  "updated" timestamptz NOT NULL DEFAULT (now()),
  FOREIGN KEY (user_id) REFERENCES "public"."_users" (id) ON DELETE CASCADE
);

CREATE INDEX ON "public"."_tables" ("user_id");
COMMENT ON COLUMN "public"."_tables"."id" IS 'numeric,server-side auto generated id';
COMMENT ON COLUMN "public"."_tables"."user_id" IS 'required,numeric,shouldref=_users.id';
COMMENT ON COLUMN "public"."_tables"."name" IS 'required,alphanumeric,unique,min=3';
COMMENT ON COLUMN "public"."_tables"."columns" IS 'required,jsonstring,desc=columns of the tables';
