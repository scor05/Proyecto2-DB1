package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Manager struct {
	conn *sql.DB
}

type Product struct {
	IDProducto  int     `json:"id_producto"`
	IDCategoria int     `json:"id_categoria"`
	IDProveedor int     `json:"id_proveedor"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	Nombre      string  `json:"nombre"`
	Imagen      *string `json:"imagen"`
	Descripcion string  `json:"descripcion"`
}

type ProductWrite struct {
	IDCategoria int     `json:"id_categoria"`
	IDProveedor int     `json:"id_proveedor"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	Nombre      string  `json:"nombre"`
	Imagen      *string `json:"imagen"`
	Descripcion string  `json:"descripcion"`
}

type ProductPatch struct {
	IDCategoria *int     `json:"id_categoria"`
	IDProveedor *int     `json:"id_proveedor"`
	Precio      *float64 `json:"precio"`
	Stock       *int     `json:"stock"`
	Nombre      *string  `json:"nombre"`
	Imagen      *string  `json:"imagen"`
	Descripcion *string  `json:"descripcion"`
}

func NewManager() (*Manager, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getenv("DB_USER", "proy2"),
		getenv("DB_PASSWORD", "secret"),
		getenv("DB_HOST", "db"),
		getenv("DB_PORT", "5432"),
		getenv("DB_NAME", "proyecto2"),
	)

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour)

	var lastErr error
	for range 15 {
		if err := conn.Ping(); err == nil {
			return &Manager{conn: conn}, nil
		} else {
			lastErr = err
		}
		time.Sleep(time.Second)
	}

	conn.Close()
	return nil, lastErr
}

func (m *Manager) Close() error {
	return m.conn.Close()
}

func (m *Manager) Ping(ctx context.Context) error {
	return m.conn.PingContext(ctx)
}

func (m *Manager) Index(ctx context.Context) ([]Product, error) {
	rows, err := m.conn.QueryContext(ctx, `
		SELECT id_producto, id_categoria, id_proveedor, precio, stock, nombre, imagen, descripcion
		FROM producto
		ORDER BY id_producto
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (m *Manager) Show(ctx context.Context, id int) (*Product, error) {
	row := m.conn.QueryRowContext(ctx, `
		SELECT id_producto, id_categoria, id_proveedor, precio, stock, nombre, imagen, descripcion
		FROM producto
		WHERE id_producto = $1
	`, id)

	product, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *Manager) Update(ctx context.Context, id int, input ProductWrite) (*Product, error) {
	if err := validateProductWrite(input); err != nil {
		return nil, err
	}

	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		UPDATE producto
		SET id_categoria = $1,
		    id_proveedor = $2,
		    precio = $3,
		    stock = $4,
		    nombre = $5,
		    imagen = $6,
		    descripcion = $7
		WHERE id_producto = $8
		RETURNING id_producto, id_categoria, id_proveedor, precio, stock, nombre, imagen, descripcion
	`, input.IDCategoria, input.IDProveedor, input.Precio, input.Stock, input.Nombre, optionalString(input.Imagen), input.Descripcion, id)

	product, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *Manager) Patch(ctx context.Context, id int, input ProductPatch) (*Product, error) {
	if err := validateProductPatch(input); err != nil {
		return nil, err
	}

	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		UPDATE producto
		SET id_categoria = COALESCE($1, id_categoria),
		    id_proveedor = COALESCE($2, id_proveedor),
		    precio = COALESCE($3, precio),
		    stock = COALESCE($4, stock),
		    nombre = COALESCE($5, nombre),
		    imagen = COALESCE($6, imagen),
		    descripcion = COALESCE($7, descripcion)
		WHERE id_producto = $8
		RETURNING id_producto, id_categoria, id_proveedor, precio, stock, nombre, imagen, descripcion
	`, optionalInt(input.IDCategoria), optionalInt(input.IDProveedor), optionalFloat64(input.Precio), optionalInt(input.Stock), optionalString(input.Nombre), optionalString(input.Imagen), optionalString(input.Descripcion), id)

	product, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *Manager) Destroy(ctx context.Context, id int) error {
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		DELETE FROM producto
		WHERE id_producto = $1
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

type productScanner interface {
	Scan(dest ...any) error
}

func scanProduct(scanner productScanner) (Product, error) {
	var product Product
	var imagen sql.NullString

	err := scanner.Scan(
		&product.IDProducto,
		&product.IDCategoria,
		&product.IDProveedor,
		&product.Precio,
		&product.Stock,
		&product.Nombre,
		&imagen,
		&product.Descripcion,
	)
	if err != nil {
		return Product{}, err
	}

	if imagen.Valid {
		product.Imagen = &imagen.String
	}

	return product, nil
}

func validateProductWrite(input ProductWrite) error {
	switch {
	case input.IDCategoria <= 0:
		return fmt.Errorf("%w: id_categoria must be greater than 0", ErrInvalidInput)
	case input.IDProveedor <= 0:
		return fmt.Errorf("%w: id_proveedor must be greater than 0", ErrInvalidInput)
	case input.Precio <= 0:
		return fmt.Errorf("%w: precio must be greater than 0", ErrInvalidInput)
	case input.Stock < 0:
		return fmt.Errorf("%w: stock cannot be negative", ErrInvalidInput)
	case strings.TrimSpace(input.Nombre) == "":
		return fmt.Errorf("%w: nombre is required", ErrInvalidInput)
	case strings.TrimSpace(input.Descripcion) == "":
		return fmt.Errorf("%w: descripcion is required", ErrInvalidInput)
	default:
		return nil
	}
}

func validateProductPatch(input ProductPatch) error {
	if input.IDCategoria == nil &&
		input.IDProveedor == nil &&
		input.Precio == nil &&
		input.Stock == nil &&
		input.Nombre == nil &&
		input.Imagen == nil &&
		input.Descripcion == nil {
		return fmt.Errorf("%w: at least one field is required", ErrInvalidInput)
	}

	if input.IDCategoria != nil && *input.IDCategoria <= 0 {
		return fmt.Errorf("%w: id_categoria must be greater than 0", ErrInvalidInput)
	}
	if input.IDProveedor != nil && *input.IDProveedor <= 0 {
		return fmt.Errorf("%w: id_proveedor must be greater than 0", ErrInvalidInput)
	}
	if input.Precio != nil && *input.Precio <= 0 {
		return fmt.Errorf("%w: precio must be greater than 0", ErrInvalidInput)
	}
	if input.Stock != nil && *input.Stock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", ErrInvalidInput)
	}
	if input.Nombre != nil && strings.TrimSpace(*input.Nombre) == "" {
		return fmt.Errorf("%w: nombre cannot be empty", ErrInvalidInput)
	}
	if input.Descripcion != nil && strings.TrimSpace(*input.Descripcion) == "" {
		return fmt.Errorf("%w: descripcion cannot be empty", ErrInvalidInput)
	}
	return nil
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func optionalInt(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func optionalFloat64(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func optionalString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}
