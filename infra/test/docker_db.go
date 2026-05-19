package test

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
)

type DockerDB struct {
	DB       *pg.DB
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

func buildDBURL(resource *dockertest.Resource, serviceName string) string {
	return fmt.Sprintf("postgres://dev:12345@%v/%s?sslmode=disable", resource.GetHostPort("5432/tcp"), serviceName)
}

func getPostgresContainerConfig(serviceName string) *dockertest.RunOptions {
	return &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.4",
		Env: []string{
			"POSTGRES_USER=dev",
			"POSTGRES_PASSWORD=12345",
			"POSTGRES_DB=" + serviceName,
			"listen_addresses = '*'",
		},
	}
}

func startPostgresContainer(pool *dockertest.Pool, serviceName string) (*dockertest.Resource, error) {
	runOptions := getPostgresContainerConfig(serviceName)
	resource, err := pool.RunWithOptions(runOptions, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		if resource != nil {
			_ = pool.Purge(resource)
		}
		return nil, fmt.Errorf("could not start PostgreSQL container. err: %w", err)
	}
	return resource, nil
}

func connectToPostgres(pool *dockertest.Pool, dbURL string) (*pg.DB, error) {
	var (
		db  *pg.DB
		err error
	)

	retryFunc := func() error {
		db, err = pg.NewDB(dbURL, "ReadCommitted")
		if err != nil {
			return err
		}
		return db.Ping()
	}

	if err = pool.Retry(retryFunc); err != nil {
		return nil, fmt.Errorf("could not connect to PostgreSQL. err: %w", err)
	}

	return db, nil
}

func getMigrationsDir() (string, error) {
	_, err := os.Stat("../migration")
	if err == nil {
		return "../migration", nil
	}

	if os.IsNotExist(err) {
		_, err = os.Stat("../repository/postgres/migration")
		if err == nil {
			return "../repository/postgres/migration", nil
		}
		if os.IsNotExist(err) {
			return "", errors.New("migrations directory not found")
		}
	}

	return "", fmt.Errorf("could not check migrations directory. err: %w", err)
}

func applyMigrations(migrationsDir, dbURL string) error {
	m, err := migrate.New("file://"+migrationsDir, dbURL)
	if err != nil {
		return fmt.Errorf("could not open migration sources. err: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not apply migrations. err: %w", err)
	}

	return nil
}

func (d *DockerDB) Close() error {
	d.DB.Close()
	return d.Pool.Purge(d.Resource)
}

func CreateDockerDB(serviceName string) (*DockerDB, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker. err: %w", err)
	}

	resource, err := startPostgresContainer(pool, serviceName)
	if err != nil {
		return nil, err
	}

	testDBURL := buildDBURL(resource, serviceName)

	db, err := connectToPostgres(pool, testDBURL)
	if err != nil {
		_ = pool.Purge(resource)
		return nil, err
	}

	migrationsDir, err := getMigrationsDir()
	if err != nil {
		db.Close()
		_ = pool.Purge(resource)
		return nil, err
	}

	if err = applyMigrations(migrationsDir, testDBURL); err != nil {
		db.Close()
		_ = pool.Purge(resource)
		return nil, fmt.Errorf("could not apply migrations. err: %w", err)
	}

	return &DockerDB{
		DB:       db,
		Pool:     pool,
		Resource: resource,
	}, nil
}
