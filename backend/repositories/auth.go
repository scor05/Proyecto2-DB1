package repositories

import (
	"context"
	"crypto/subtle"
	"strings"

	"gorm.io/gorm"
	"proyecto2/backend/models"
)

type credentialRow struct {
	ID       int    `gorm:"column:id"`
	Nombre   string `gorm:"column:nombre"`
	Correo   string `gorm:"column:correo"`
	Password string `gorm:"column:password"`
}

type authSource struct {
	table  string
	idCol  string
	role   string
	active string
}

func (r *Repository) Authenticate(ctx context.Context, correo string, password string) (*models.AuthUser, error) {
	email := strings.TrimSpace(correo)
	pass := strings.TrimSpace(password)
	if email == "" || pass == "" {
		return nil, ErrInvalidCredentials
	}

	sources := []authSource{
		{table: "cliente", idCol: "id_cliente", role: "cliente"},
		{table: "proveedor", idCol: "id_proveedor", role: "proveedor"},
		{table: "empleado", idCol: "id_empleado", role: "empleado"},
		{table: "gerente", idCol: "id_gerente", role: "gerente"},
		{table: "superadmin", idCol: "id_superadmin", role: "superadmin"},
	}

	for _, source := range sources {
		user, err := r.authenticateFromSource(ctx, source, email, pass)
		if err == nil {
			return user, nil
		}
		if err != ErrNotFound && err != ErrInvalidCredentials {
			return nil, err
		}
	}

	return nil, ErrInvalidCredentials
}

func (r *Repository) authenticateFromSource(ctx context.Context, source authSource, correo string, password string) (*models.AuthUser, error) {
	var row credentialRow
	query := r.db.WithContext(ctx).
		Table(source.table).
		Select(source.idCol+" AS id, nombre, correo, password").
		Where("lower(correo) = lower(?)", correo)

	if source.active != "" {
		query = query.Where(source.active)
	}

	if err := query.Take(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if subtle.ConstantTimeCompare([]byte(row.Password), []byte(password)) != 1 {
		return nil, ErrInvalidCredentials
	}

	return &models.AuthUser{
		ID:     row.ID,
		Nombre: row.Nombre,
		Correo: row.Correo,
		Rol:    source.role,
	}, nil
}
