package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mayye4ka/notpastebin/internal/errs"
	"github.com/mayye4ka/notpastebin/internal/service"
	"github.com/rs/zerolog"
)

type db struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

func New(pool *pgxpool.Pool, logger *zerolog.Logger) *db {
	return &db{
		pool:   pool,
		logger: logger,
	}
}

func (d *db) CreateNote(ctx context.Context, text, adminHash, hash string) error {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	countQuery := `select count(1) from notes
		where  admin_hash = @admin_hash
		or admin_hash = @reader_hash
		or reader_hash = @admin_hash
		or reader_hash = @reader_hash;`
	row, err := conn.Query(ctx, countQuery, pgx.NamedArgs{
		"admin_hash":  adminHash,
		"reader_hash": hash,
	})
	if err != nil {
		d.logger.Error().Err(err).Msg("check hash error")
		return fmt.Errorf("check hash is unique: %w", errs.ErrInternalError)
	}
	// TODO: maybee do row.Next
	var hits int
	err = row.Scan(&hits)
	if err != nil {
		d.logger.Error().Err(err).Msg("check hash error")
		return fmt.Errorf("check hash is unique: %w", errs.ErrInternalError)
	} else if hits > 0 {
		return fmt.Errorf("hash already used: %w", errs.ErrCollision)
	}

	createQuery := `INSERT INTO notes VALUES(@admin_hash, @reader_hash, @note)`
	_, err = conn.Exec(ctx, createQuery, pgx.NamedArgs{
		"admin_hash":  adminHash,
		"reader_hash": hash,
		"note":        text,
	})
	if err != nil {
		d.logger.Error().Err(err).Msg("create note error")
		return fmt.Errorf("can't create note: %w", errs.ErrInternalError)
	}
	return nil
}

func (d *db) GetNote(ctx context.Context, hash string) (service.GetNoteResponse, error) {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return service.GetNoteResponse{}, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	getQuery := `SELECT note, admin_hash, reader_hash FROM notes WHERE admin_hash = @hash or reader_hash = @hash`
	row := conn.QueryRow(ctx, getQuery, pgx.NamedArgs{
		"hash": hash,
	})
	var note, adminHash, readerHash string
	// TODO: maybee do row.Next
	// TODO: check only one found
	err = row.Scan(&note, &adminHash, &readerHash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return service.GetNoteResponse{}, fmt.Errorf("note not found: %w", errs.ErrNotFound)
	} else if err != nil {
		d.logger.Error().Err(err).Msg("get note error")
		return service.GetNoteResponse{}, fmt.Errorf("can't get note: %w", errs.ErrInternalError)
	}
	return service.GetNoteResponse{
		Note:       note,
		ReaderHash: readerHash,
		IsAdmin:    hash == adminHash,
	}, nil
}

func (d *db) UpdateNote(ctx context.Context, hash string, text string) error {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	createQuery := `UPDATE notes SET note = @note WHERE admin_hash = @admin_hash`
	row, err := conn.Exec(ctx, createQuery, pgx.NamedArgs{
		"admin_hash": hash,
		"note":       text,
	})
	if err != nil {
		d.logger.Error().Err(err).Msg("update note error")
		return fmt.Errorf("can't update note: %w", errs.ErrInternalError)
	}
	affected := row.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("note not found: %w", errs.ErrInvalidInput)
	}
	return nil
}

func (d *db) DeleteNote(ctx context.Context, hash string) error {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	createQuery := `DELETE FROM notes WHERE admin_hash = @admin_hash`
	row, err := conn.Exec(ctx, createQuery, pgx.NamedArgs{
		"admin_hash": hash,
	})
	if err != nil {
		d.logger.Error().Err(err).Msg("delete note error")
		return fmt.Errorf("can't delete note %w", errs.ErrInternalError)
	}
	affected := row.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("can't delete note: %w", errs.ErrInvalidInput)
	}
	return nil
}
