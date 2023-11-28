package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type StcDB struct {
	Conn *sqlx.DB
}

func NewInMemStcDB() *StcDB {
	return &StcDB{Conn: sqlx.MustConnect("sqlite3", ":memory:")}
}

func (s *StcDB) Migrate() {
	s.Conn.MustExec(`
	CREATE TABLE IF NOT EXISTS "buckets" (
		"id"	INTEGER NOT NULL,
		"name"	TEXT NOT NULL UNIQUE,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	`)

	s.Conn.MustExec(`
	CREATE TABLE IF NOT EXISTS "bucket_data" (
		"bucket_id"	INTEGER NOT NULL,
		"connector_run_id" INTEGER,
		"created_at" TEXT,
		"raw" BLOB,
		"type" TEXT,
		"meta" TEXT
	);
	`)

	s.Conn.MustExec(`
	CREATE TABLE IF NOT EXISTS "connectors" (
		"id"	INTEGER NOT NULL,
		"bucket_id"	INTEGER NOT NULL,
		"conn_type"	INTEGER NOT NULL,
		"trigger"	INTEGER NOT NULL,
		"params"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	`)

	s.Conn.MustExec(`
	CREATE TABLE IF NOT EXISTS "connector_runs" (
		"id"	INTEGER NOT NULL,
		"connector_id"	INTEGER NOT NULL,
		"error"	TEXT,
		"started_at"	TEXT,
		"finished_at"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	`)
}
