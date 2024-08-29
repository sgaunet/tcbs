package tcbs

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"context"
	"fmt"

	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const PostgresDockerImage = "postgres:16.2"
const DeadlineTimeout = 20 * time.Second

// TestDB is a struct that holds the postgresql container and the DSN
type PostgresContainer struct {
	postgresqlC    testcontainers.Container
	dataSourceName dsn.DSN
}

// NewTestDB creates a postgresql database in docker for tests
func NewPostgresContainer(postgresUser, postgresPassword, postgresDBName string) (*PostgresContainer, error) {
	var err error
	newpgDB := &PostgresContainer{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(DeadlineTimeout))
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16.2",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDBName,
		},
	}
	newpgDB.postgresqlC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("could not start postgres container: %v", err)
	}
	endpoint, err := newpgDB.postgresqlC.Endpoint(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("could not get postgres container endpoint: %v", err)
	}
	fmt.Println("Postgres container started on", endpoint)
	newpgDB.dataSourceName, err = dsn.New(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", postgresUser, postgresPassword, endpoint, postgresDBName))
	if err != nil {
		return nil, fmt.Errorf("could not create DSN: %v", err)
	}
	// Wait for the database to be ready
	err = waitForDB(ctx, newpgDB.dataSourceName.GetPostgresUri())
	if err != nil {
		return nil, fmt.Errorf("could not wait for postgres container to be ready: %v", err)
	}
	return newpgDB, nil
}

// waitForDB waits for the database to be ready
func waitForDB(ctx context.Context, pginfo string) error {
	chDBReady := make(chan struct{})
	go func() {
		for {
			db, err := sql.Open("postgres", pginfo)
			select {
			case <-ctx.Done():
				return
			default:
				if err == nil {
					err = db.Ping()
					defer db.Close()
					if err == nil {
						close(chDBReady)
						return
					}
					fmt.Println("Database not ready (not pingable)")
					time.Sleep(1 * time.Second)
				}
				// fmt.Println("Waiting for database to be ready...", pgdsn, err.Error())
			}
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-chDBReady:
		return nil
	}
}

// Terminate stops the postgresql container
func (pg *PostgresContainer) Terminate(ctx context.Context) error {
	return pg.postgresqlC.Terminate(ctx)
}

// GetDSN returns the DSN
func (pg *PostgresContainer) GetDSN() dsn.DSN {
	return pg.dataSourceName
}

// GetDSNString returns the DSN string
// The format is 'host=... port=... user=... password=... dbname=... sslmode=...'
func (pg *PostgresContainer) GetDSNString() string {
	return pg.dataSourceName.GetPostgresUri()
}

// GetDBUser returns the user of the database
func (pg *PostgresContainer) GetDBUser() string {
	return pg.dataSourceName.GetUser()
}

// GetDBPassword returns the password of the database
func (pg *PostgresContainer) GetDBPassword() string {
	return pg.dataSourceName.GetPassword()
}

// GetDBHost returns the host of the database
func (pg *PostgresContainer) GetDBHost() string {
	return pg.dataSourceName.GetHost()
}

// GetDBPort returns the port of the database
func (pg *PostgresContainer) GetDBPort() string {
	return pg.dataSourceName.GetPort("5432")
}

// GetDBPortInt returns the port of the database as an integer
func (pg *PostgresContainer) GetDBPortInt() int {
	return pg.dataSourceName.GetPortInt(5432)
}

// GetDBName returns the name of the database
func (pg *PostgresContainer) GetDBName() string {
	return pg.dataSourceName.GetDBName()
}
