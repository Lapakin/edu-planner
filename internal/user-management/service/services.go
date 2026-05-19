package service

import (
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/user-management/repository"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
)

func NewServices(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *Services {
	return &Services{
		UserSvc:         NewUserService(db, rm, l.WithField("service", "UserService")),
		AuthSvc:         NewAuthService(db, rm, l.WithField("service", "AuthService")),
		AcademicYearSvc: NewAcademicYearService(db, rm, l.WithField("service", "AcademicYearInfoService")),
	}
}
