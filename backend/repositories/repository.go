package repositories

import (
	"errors"
	"fmt"
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

func storedProcedureError(err error) error {
	if err == nil {
		return nil
	}
	message := err.Error()
	if strings.Contains(message, "producto_not_found") || strings.Contains(message, "compra_not_found") {
		return ErrNotFound
	}
	if strings.Contains(message, "producto_referencia_invalida") ||
		strings.Contains(message, "producto_datos_invalidos") ||
		strings.Contains(message, "producto_duplicado") {
		return fmt.Errorf("%w: stored procedure rejected product data", ErrInvalidInput)
	}
	return err
}
