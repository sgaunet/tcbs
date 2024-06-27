package tcbs_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sgaunet/tcbs"
)

func TestNewPostgresContainer(t *testing.T) {
	newpgDB, err := tcbs.NewPostgresContainer("postgres", "password", "postgres")
	if err != nil {
		t.Fatalf("could not create postgres container: %v", err)
	}
	defer newpgDB.Terminate(context.Background())
	db, err := sql.Open("postgres", newpgDB.GetDSNString())
	if err != nil {
		t.Fatalf("could not open postgres connection: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatalf("could not ping postgres: %v", err)
	}
}
