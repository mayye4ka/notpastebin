package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/mayye4ka/notpastebin/internal/errs"
)

const (
	maxNoteLen   = 65535
	hashLen      = 32
	hashAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type service struct {
	db Repository
}

type Repository interface {
	CreateNote(ctx context.Context, text, adminHash, readerHash string) error
	GetNote(ctx context.Context, hash string) (string, error)
	UpdateNote(ctx context.Context, hash, text string) error
	DeleteNote(ctx context.Context, hash string) error
}

func New(db Repository) *service {
	return &service{
		db: db,
	}
}

func (s *service) CreateNote(ctx context.Context, text string) (string, string, error) {
	err := validateNote(text)
	if err != nil {
		return "", "", err
	}
	for {
		adminHash, readerHash := createHash(), createHash()
		if adminHash == readerHash {
			continue
		}
		err = s.db.CreateNote(ctx, text, adminHash, readerHash)
		if err == nil {
			return adminHash, readerHash, nil
		}
		if errors.Is(err, errs.ErrCollision) {
			continue
		}
		return "", "", fmt.Errorf("create note: %w", err)
	}
}

func (s *service) GetNote(ctx context.Context, hash string) (string, error) {
	return s.db.GetNote(ctx, hash)
}

func (s *service) UpdateNote(ctx context.Context, hash string, text string) error {
	err := validateNote(text)
	if err != nil {
		return err
	}
	return s.db.UpdateNote(ctx, hash, text)
}

func (s *service) DeleteNote(ctx context.Context, hash string) error {
	return s.db.DeleteNote(ctx, hash)
}

func createHash() string {
	b := make([]byte, hashLen)
	for i := range b {
		b[i] = hashAlphabet[rand.Intn(len(hashAlphabet))]
	}
	return string(b)
}

func validateNote(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("%w: empty text", errs.ErrInvalidInput)
	}
	if len(text) > maxNoteLen {
		return fmt.Errorf("%w: text longer than %d symbols", errs.ErrInvalidInput, maxNoteLen)
	}
	return nil
}
