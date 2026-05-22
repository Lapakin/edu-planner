package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
)

type AuthRepository struct {
	db sqlx.ExtContext
}

func NewAuthRepository(db sqlx.ExtContext) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserCredentialByEmail(ctx context.Context, email string) (*domain.UserCredential, error) {
	credQuery, credArgs, err := q.Select(ctx).
		Columns("uc.user_id, uc.password_hash, uc.updated_at").
		From(`"user_credential" uc`).
		Join(`"user" u ON u.id = uc.user_id`).
		Where("u.email = ?", email).
		ToSQL()
	if err != nil {
		return nil, err
	}

	cred := &domain.UserCredential{}
	if err = sqlx.GetContext(ctx, r.db, cred, credQuery, credArgs...); err != nil {
		return nil, err
	}

	userQuery, userArgs, err := q.Select(ctx).
		Columns("id, email, role, is_active, is_deleted, created_at, modified_at").
		From(`"user"`).
		Where("id = ?", cred.UserID).
		ToSQL()
	if err != nil {
		return nil, err
	}

	user := &domain.User{}
	if err = sqlx.GetContext(ctx, r.db, user, userQuery, userArgs...); err != nil {
		return nil, err
	}

	cred.User = user
	return cred, nil
}

func (r *AuthRepository) CreateUserCredential(ctx context.Context, cred *domain.UserCredential) error {
	query, args, err := q.Insert(ctx).
		Into(`"user_credential"`).
		Columns(`
			user_id,
			password_hash,
			updated_at
		`).
		Values(
			cred.UserID,
			cred.PasswordHash,
			cred.UpdatedAt,
		).
		ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	query, args, err := q.Insert(ctx).
		Into(`"refresh_token"`).
		Columns(`
			user_id,
			token_hash,
			expires_at,
			is_revoked,
			created_at
		`).
		Values(
			token.UserID,
			token.TokenHash,
			token.ExpiresAt,
			token.IsRevoked,
			token.CreatedAt,
		).
		Returning("id").
		ToSQL()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&token.ID); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query, args, err := q.Update(ctx).
		Table(`"refresh_token"`).
		Set("is_revoked", true).
		Where("token_hash = ?", tokenHash).
		ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			user_id,
			token_hash,
			expires_at,
			is_revoked,
			created_at
		`).
		From(`"refresh_token"`).
		Where("token_hash = ?", tokenHash).
		ToSQL()
	if err != nil {
		return nil, err
	}

	token := &domain.RefreshToken{}
	if err = sqlx.GetContext(ctx, r.db, token, query, args...); err != nil {
		return nil, err
	}

	return token, nil
}

func (r *AuthRepository) UpsertUserCredential(ctx context.Context, cred *domain.UserCredential) error {
	query, args, err := q.Insert(ctx).
		Into(`"user_credential"`).
		Columns("user_id", "password_hash", "updated_at").
		Values(cred.UserID, cred.PasswordHash, cred.UpdatedAt).
		OnConflictDoUpdate("user_id", "password_hash", "updated_at").
		ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) CreateInviteToken(ctx context.Context, token *domain.InviteToken) error {
	query, args, err := q.Insert(ctx).
		Into(`"invite_token"`).
		Columns(`
			user_id,
			token,
			expires_at,
			created_at
		`).
		Values(
			token.UserID,
			token.Token,
			token.ExpiresAt,
			token.CreatedAt,
		).
		Returning("id").
		ToSQL()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&token.ID); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetInviteToken(ctx context.Context, token string) (*domain.InviteToken, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			user_id,
			token,
			expires_at,
			used_at,
			created_at
		`).
		From(`"invite_token"`).
		Where("token = ?", token).
		ToSQL()
	if err != nil {
		return nil, err
	}

	it := &domain.InviteToken{}
	if err = sqlx.GetContext(ctx, r.db, it, query, args...); err != nil {
		return nil, err
	}

	return it, nil
}

func (r *AuthRepository) MarkInviteTokenUsed(ctx context.Context, token string, usedAt time.Time) error {
	query, args, err := q.Update(ctx).
		Table(`"invite_token"`).
		Set("used_at", usedAt).
		Where("token = ?", token).
		ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) CreateResetPasswordToken(ctx context.Context, token *domain.ResetPasswordToken) error {
	query, args, err := q.Insert(ctx).
		Into(`"reset_password_token"`).
		Columns(`
			user_id,
			token,
			expires_at,
			created_at
		`).
		Values(
			token.UserID,
			token.Token,
			token.ExpiresAt,
			token.CreatedAt,
		).
		Returning("id").
		ToSQL()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&token.ID); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetResetPasswordToken(ctx context.Context, token string) (*domain.ResetPasswordToken, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			user_id,
			token,
			expires_at,
			used_at,
			created_at
		`).
		From(`"reset_password_token"`).
		Where("token = ?", token).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rt := &domain.ResetPasswordToken{}
	if err = sqlx.GetContext(ctx, r.db, rt, query, args...); err != nil {
		return nil, err
	}

	return rt, nil
}

func (r *AuthRepository) MarkResetPasswordTokenUsed(ctx context.Context, token string, usedAt time.Time) error {
	query, args, err := q.Update(ctx).
		Table(`"reset_password_token"`).
		Set("used_at", usedAt).
		Where("token = ?", token).
		ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
