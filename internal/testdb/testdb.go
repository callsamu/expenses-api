package testdb

import (
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const migrations = "file://../../migrations"
const dsn = "postgres://pfapi_test:postgres@localhost:5432/pfapi_test?sslmode=disable"

type TestDB struct {
	DB *sql.DB
	m  *migrate.Migrate
}

func Open(t *testing.T) *TestDB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open test db: %s", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrations, "postgres", driver)
	if err != nil {
		t.Fatalf("failed to finds migrations: %s", err.Error())
	}
	m.Up()

	return &TestDB{
		DB: db,
		m:  m,
	}
}

func (tdb *TestDB) Close() {
	tdb.m.Drop()
	tdb.DB.Close()
}
