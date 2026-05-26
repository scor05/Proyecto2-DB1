package repositories

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrNotFound           = errors.New("record not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func notFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

func storedProcedureNotFound(err error) error {
	if err == nil {
		return nil
	}
	message := err.Error()
	if strings.Contains(message, "producto_not_found") || strings.Contains(message, "compra_not_found") {
		return ErrNotFound
	}
	return err
}
