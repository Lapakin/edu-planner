package integration_test

import (
	"log"
	"os"
	"testing"

	"github.com/Lapakin/edu-planner/infra/test"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
)

const serviceName = "user-management"

var db *pg.DB

func TestMain(m *testing.M) {
	exitCode := 0
	if test.GetTestingType() == test.UnitTesting {
		os.Exit(exitCode)
	}

	dockerDB, err := test.CreateDockerDB(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	db = dockerDB.DB

	defer func() {
		if err = dockerDB.Close(); err != nil {
			log.Printf("failed to close docker DB: %v", err)
		}
		os.Exit(exitCode)
	}()

	exitCode = m.Run()
}
