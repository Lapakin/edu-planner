package integration_test

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Lapakin/edu-planner/infra/test"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/user-management/repository/postgres"
	"github.com/Lapakin/edu-planner/internal/user-management/router"
	"github.com/Lapakin/edu-planner/internal/user-management/service"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
	jwtAdapter "github.com/Lapakin/edu-planner/internal/adapter/jwt"
	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

const serviceName = "user-management"

var (
	db           *pg.DB
	ts           *httptest.Server
	adminToken   string
	deanToken    string
	teacherToken string
)

func TestMain(m *testing.M) {
	exitCode := 0
	if test.GetTestingType() == test.UnitTesting {
		os.Exit(exitCode)
	}

	jwtAdapter.Init(&jwtAdapter.Config{Secret: "test-secret", ExpireMinutes: 60})

	var err error
	adminTokenRaw, err := jwtAdapter.GenerateToken(1, "admin@test.com", "admin")
	if err != nil {
		log.Fatalf("failed to generate admin token: %v", err)
	}
	adminToken = "Bearer " + adminTokenRaw

	deanTokenRaw, err := jwtAdapter.GenerateToken(2, "dean@test.com", "dean")
	if err != nil {
		log.Fatalf("failed to generate dean token: %v", err)
	}
	deanToken = "Bearer " + deanTokenRaw

	teacherTokenRaw, err := jwtAdapter.GenerateToken(3, "teacher@test.com", "teacher")
	if err != nil {
		log.Fatalf("failed to generate teacher token: %v", err)
	}
	teacherToken = "Bearer " + teacherTokenRaw

	dockerDB, err := test.CreateDockerDB(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	db = dockerDB.DB

	svc := service.NewServices(
		db,
		postgres.NewRepoManager(),
		logging.NewLogger(logging.TraceLevel, logging.PrettyFormatter),
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
		log.Printf("failed to initialize test data: %v", err)
		exitCode = 1
		return
	}

	exitCode = m.Run()
}

func initTestData() error {
	ctx := context.Background()
	ayRepo := postgres.NewAcademicYearRepository(db)
	if err := ayRepo.UpsertAcademicYear(ctx, ta.AcademicYear1); err != nil {
		return err
	}
	if err := ayRepo.UpsertAcademicYear(ctx, ta.AcademicYear2); err != nil {
		return err
	}
	return nil
}
