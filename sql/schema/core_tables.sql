CREATE TABLE "core_tables" (
  "tid" BIGSERIAL PRIMARY KEY NOT NULL,
  "uid" BIGINT NOT NULL,
  "tablename" VARCHAR(60) UNIQUE NOT NULL,
  "columns" VARCHAR NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  FOREIGN KEY (uid) REFERENCES "public"."core_users" (uid)
);

CREATE INDEX ON "public"."core_tables" ("uid");
COMMENT ON COLUMN "public"."core_tables"."uid" IS 'required,numeric,shouldref=core_users.uid';
COMMENT ON COLUMN "public"."core_tables"."columns" IS 'required,jsonstring,desc=columns of the tables';
