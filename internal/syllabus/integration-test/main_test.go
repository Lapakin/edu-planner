package integration_test

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Lapakin/edu-planner/infra/test"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository/postgres"
	"github.com/Lapakin/edu-planner/internal/syllabus/router"
	"github.com/Lapakin/edu-planner/internal/syllabus/service"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
	jwtAdapter "github.com/Lapakin/edu-planner/internal/adapter/jwt"
	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

const serviceName = "syllabus"

var (
	db         *pg.DB
	ts         *httptest.Server
	adminToken string
)

func TestMain(m *testing.M) {
	exitCode := 0
	if test.GetTestingType() == test.UnitTesting {
		os.Exit(exitCode)
	}

	jwtAdapter.Init(&jwtAdapter.Config{Secret: "test-secret", ExpireMinutes: 60})

	adminTokenRaw, err := jwtAdapter.GenerateToken(1, "admin@test.com", "admin")
	if err != nil {
		log.Fatalf("failed to generate admin token: %v", err)
	}
	adminToken = "Bearer " + adminTokenRaw

	dockerDB, err := test.CreateDockerDB(serviceName)
	if err != nil {
		panic(err)
	}
	db = dockerDB.DB

	svc := service.NewServices(
		db,
		postgres.NewRepoManager(),
		logging.NewLogger(logging.TraceLevel, logging.PrettyFormatter),
		domain.DefaultGenerationConfig(),
	)

	r := router.NewRouter(svc)
	ts = httptest.NewServer(r)

	defer func() {
		ts.Close()
		if err = dockerDB.Close(); err != nil {
			log.Printf("failed to close docker DB: %v", err)
		}
		os.Exit(exitCode)
	}()

	if err = initTestData(); err != nil {
		log.Printf("failed to init test data: %v", err)
	}

	exitCode = m.Run()
}

func initTestData() error {
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(db)
	if err := userRepo.UpsertUser(ctx, ta.User1); err != nil {
		return err
	}
	if err := userRepo.UpsertUser(ctx, ta.User2); err != nil {
		return err
	}
	return nil
}
