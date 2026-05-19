package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/user-management/repository"
)

type RepoManager struct{}

func NewRepoManager() *RepoManager {
	return &RepoManager{}
}

func (rm *RepoManager) NewUserRepo(db sqlx.ExtContext) repository.UserRepository {
	return NewUserRepository(db)
}

func (rm *RepoManager) NewMessageRepo(db sqlx.ExtContext) repository.MessageRepository {
	return NewMessageWriter(db)
}

func (rm *RepoManager) NewAuthRepo(db sqlx.ExtContext) repository.AuthRepository {
	return NewAuthRepository(db)
}

func (rm *RepoManager) NewAcademicYearRepo(db sqlx.ExtContext) repository.AcademicYearRepository {
	return NewAcademicYearRepository(db)
}
