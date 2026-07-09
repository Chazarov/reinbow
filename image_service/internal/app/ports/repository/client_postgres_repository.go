package repository

import (
	"context"
	"database/sql"
	"time"
	"traineesheep/imageservice/internal/app/models"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
)

// ClientPostgresRepository хранит учётные записи клиентов в PostgreSQL.
type ClientPostgresRepository struct {
	db     *sql.DB
	logger *logger.ContextLogger
}

// NewClientPostgreRepository создаёт репозиторий клиентов в PostgreSQL.
func NewClientPostgreRepository(db *sql.DB, logger *logger.ContextLogger) *ClientPostgresRepository {
	return &ClientPostgresRepository{
		db:     db,
		logger: logger,
	}
}

// CreateClient сохраняет нового клиента в базе данных.
func (r *ClientPostgresRepository) CreateClient(ctx context.Context, client models.CreateClientMto) (models.Client, errors.DomainError) {
	query := `
		INSERT INTO clients (pass_hash, name, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, pass_hash, name, created_at
	`

	var createdClient models.Client
	err := r.db.QueryRowContext(
		ctx,
		query,
		client.PassHash,
		client.Name,
		time.Now(),
	).Scan(
		&createdClient.Id,
		&createdClient.PassHash,
		&createdClient.Name,
		&createdClient.CreatedAt,
	)

	if err != nil {
		r.logger.Errorf(ctx, "Failed to create client: %v", err)
		derr := errors.NewInternalServiceError("failed to create client", err)
		return models.Client{}, derr
	}

	r.logger.Infof(ctx, "Client created with id: %d", createdClient.Id)
	return createdClient, nil
}

// GetClientByName возвращает клиента по логину.
func (r *ClientPostgresRepository) GetClientByName(ctx context.Context, name string) (models.Client, errors.DomainError) {
	query := `
		SELECT id, pass_hash, name, created_at 
		FROM clients 
		WHERE name = $1
	`

	var client models.Client
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&client.Id,
		&client.PassHash,
		&client.Name,
		&client.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Errorf(ctx, "Client with name %s not found", name)
			derr := errors.NewNotFoundError("client", name)
			return models.Client{}, derr
		}
		r.logger.Errorf(ctx, "Failed to get client by name %s: %v", name, err)
		derr := errors.NewInternalServiceError("failed to get client", err)
		return models.Client{}, derr
	}

	return client, nil
}

// UpdateClientToken обновляет токен клиента.
func (r *ClientPostgresRepository) UpdateClientToken(ctx context.Context, userId int, newToken string) errors.DomainError {
	query := `
		UPDATE clients 
		SET token = $1 
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, newToken, userId)
	if err != nil {
		r.logger.Errorf(ctx, "Failed to update client %d token: %v", userId, err)
		return errors.NewInternalServiceError("failed to update client token", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf(ctx, "Failed to get rows affected for client %d: %v", userId, err)
		return errors.NewInternalServiceError("failed to update client token", err)
	}

	if rowsAffected == 0 {
		r.logger.Errorf(ctx, "Client with id %d not found for token update", userId)
		return errors.NewNotFoundError("client", userId)
	}

	r.logger.Infof(ctx, "Client %d token updated successfully", userId)
	return nil
}

// GetClientByID возвращает клиента по идентификатору.
func (r *ClientPostgresRepository) GetClientByID(ctx context.Context, id string) (models.Client, errors.DomainError) {
	query := `
		SELECT id, pass_hash, name, created_at 
		FROM clients 
		WHERE id = $1
	`

	var client models.Client
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&client.Id,
		&client.PassHash,
		&client.Name,
		&client.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Errorf(ctx, "Client with id %s not found", id)
			derr := errors.NewNotFoundError("client", id)
			return models.Client{}, derr
		}
		r.logger.Errorf(ctx, "Failed to get client by name %s: %v", id, err)
		derr := errors.NewInternalServiceError("failed to get client", err)
		return models.Client{}, derr
	}

	return client, nil
}
