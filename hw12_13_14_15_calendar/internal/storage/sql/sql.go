package sql

func (p *PgRepo) GetSQLs1() [][2]string {
	sqls := [][2]string{}

	sqls = append(
		sqls,
		[2]string{`
DROP TABLE IF EXISTS "events"`,
			"Failed to drop table: "},
		// Sequence structure for 'events_id_seq'
		[2]string{`
DROP SEQUENCE IF EXISTS "public"."events_id_seq"`,
			"Failed to drop sequence: "},
		[2]string{`
CREATE SEQUENCE "public"."events_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1`,
			"Failed to create sequence: "},
		[2]string{`
ALTER SEQUENCE "public"."events_id_seq" 
OWNER TO "user"`,
			"Failed to alter sequence owner: "},
		// Table structure for events
		[2]string{`
CREATE TABLE "public"."events" (
	"id" BIGINT NOT NULL DEFAULT nextval('events_id_seq'::regclass),
	"name" VARCHAR(255) NOT NULL,
	"date" TIMESTAMP NOT NULL,
	"expiry" TIMESTAMP NOT NULL,
	"description" TEXT,
	"user_id" BIGINT NOT NULL,
	"time_alarm" TIMESTAMP NOT NULL
)`,
			"Failed to create table: "},
	)
	return sqls
}

func (p *PgRepo) GetSQLs2() [][2]string {
	sqls := [][2]string{}
	sqls = append(
		sqls,
		// Table structure for events
		[2]string{`
ALTER TABLE "public"."events" 
OWNER TO "user"`,
			"Failed to alter table owner: "},
		// Alter sequences owned by
		[2]string{`
ALTER SEQUENCE "public"."events_id_seq"
OWNED BY "public"."events"."id"`,
			"Failed to alter sequence owned by: "},
		[2]string{`
SELECT setval('"public"."events_id_seq"', 10, true)`,
			"Failed to set sequence value: "},
		// Primary Key structure for table events
		[2]string{`
ALTER TABLE "public"."events" 
ADD CONSTRAINT "events_pkey" 
PRIMARY KEY ("id")`,
			"Failed to add primary key: "},
	)
	return sqls
}
