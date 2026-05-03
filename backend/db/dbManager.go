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
	Categoria   string  `json:"categoria"`
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

type Employee struct {
	IDEmpleado int    `json:"id_empleado"`
	Nombre     string `json:"nombre"`
	Estado     string `json:"estado"`
	Correo     string `json:"correo"`
}

type Category struct {
	IDCategoria int    `json:"id_categoria"`
	Nombre      string `json:"nombre"`
}

type CategoryWrite struct {
	Nombre string `json:"nombre"`
}

type Provider struct {
	IDProveedor int    `json:"id_proveedor"`
	Nombre      string `json:"nombre"`
	Telefono    string `json:"telefono"`
	Correo      string `json:"correo"`
	Direccion   string `json:"direccion"`
}

type ProviderWrite struct {
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
	Direccion string `json:"direccion"`
}

type Client struct {
	IDCliente int    `json:"id_cliente"`
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
}

type Compra struct {
	IDCompra       int     `json:"id_compra"`
	IDEmpleado     int     `json:"id_empleado"`
	NombreEmpleado string  `json:"nombre_empleado"`
	IDCliente      int     `json:"id_cliente"`
	NombreCliente  string  `json:"nombre_cliente"`
	FechaCompra    string  `json:"fecha_compra"`
	IDProducto     int     `json:"id_producto"`
	Productos      string  `json:"productos"`
	TotalCompra    float64 `json:"total_compra"`
}

type CompraWrite struct {
	IDEmpleado  int     `json:"id_empleado"`
	IDCliente   int     `json:"id_cliente"`
	FechaCompra string  `json:"fecha_compra"`
	IDProducto  int     `json:"id_producto"`
	TotalCompra float64 `json:"total_compra"`
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
		SELECT p.id_producto, p.id_categoria, p.id_proveedor, c.nombre AS categoria,
		       p.precio, p.stock, p.nombre, p.imagen, p.descripcion
		FROM producto p
		JOIN categoria c ON c.id_categoria = p.id_categoria
		ORDER BY p.id_producto
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
		SELECT p.id_producto, p.id_categoria, p.id_proveedor, c.nombre AS categoria,
		       p.precio, p.stock, p.nombre, p.imagen, p.descripcion
		FROM producto p
		JOIN categoria c ON c.id_categoria = p.id_categoria
		WHERE p.id_producto = $1
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

func (m *Manager) Create(ctx context.Context, input ProductWrite) (*Product, error) {
	if err := validateProductWrite(input); err != nil {
		return nil, err
	}

	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO producto (id_categoria, id_proveedor, precio, stock, nombre, imagen, descripcion)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_producto, id_categoria, id_proveedor,
		          (SELECT nombre FROM categoria WHERE id_categoria = producto.id_categoria) AS categoria,
		          precio, stock, nombre, imagen, descripcion
	`, input.IDCategoria, input.IDProveedor, input.Precio, input.Stock, input.Nombre, optionalString(input.Imagen), input.Descripcion)

	product, err := scanProduct(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
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
		RETURNING id_producto, id_categoria, id_proveedor,
		          (SELECT nombre FROM categoria WHERE id_categoria = producto.id_categoria) AS categoria,
		          precio, stock, nombre, imagen, descripcion
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
		RETURNING id_producto, id_categoria, id_proveedor,
		          (SELECT nombre FROM categoria WHERE id_categoria = producto.id_categoria) AS categoria,
		          precio, stock, nombre, imagen, descripcion
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

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM producto_compra
		WHERE id_producto = $1
	`, id); err != nil {
		return err
	}

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

func (m *Manager) LoginEmployee(ctx context.Context, correo string) (*Employee, error) {
	row := m.conn.QueryRowContext(ctx, `
		SELECT id_empleado, nombre, estado, correo
		FROM empleado
		WHERE lower(correo) = lower($1)
	`, strings.TrimSpace(correo))

	var employee Employee
	if err := row.Scan(&employee.IDEmpleado, &employee.Nombre, &employee.Estado, &employee.Correo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (m *Manager) Categories(ctx context.Context) ([]Category, error) {
	rows, err := m.conn.QueryContext(ctx, `
		SELECT id_categoria, nombre
		FROM categoria
		ORDER BY nombre
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.IDCategoria, &category.Nombre); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *Manager) Category(ctx context.Context, id int) (*Category, error) {
	row := m.conn.QueryRowContext(ctx, `
		SELECT id_categoria, nombre
		FROM categoria
		WHERE id_categoria = $1
	`, id)

	var category Category
	if err := row.Scan(&category.IDCategoria, &category.Nombre); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (m *Manager) CreateCategory(ctx context.Context, input CategoryWrite) (*Category, error) {
	if err := validateCategoryWrite(input); err != nil {
		return nil, err
	}

	row := m.conn.QueryRowContext(ctx, `
		INSERT INTO categoria (nombre)
		VALUES ($1)
		RETURNING id_categoria, nombre
	`, strings.TrimSpace(input.Nombre))

	var category Category
	if err := row.Scan(&category.IDCategoria, &category.Nombre); err != nil {
		return nil, err
	}
	return &category, nil
}

func (m *Manager) UpdateCategory(ctx context.Context, id int, input CategoryWrite) (*Category, error) {
	if err := validateCategoryWrite(input); err != nil {
		return nil, err
	}

	row := m.conn.QueryRowContext(ctx, `
		UPDATE categoria
		SET nombre = $1
		WHERE id_categoria = $2
		RETURNING id_categoria, nombre
	`, strings.TrimSpace(input.Nombre), id)

	var category Category
	if err := row.Scan(&category.IDCategoria, &category.Nombre); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (m *Manager) DestroyCategory(ctx context.Context, id int) error {
	result, err := m.conn.ExecContext(ctx, `
		DELETE FROM categoria
		WHERE id_categoria = $1
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
	return nil
}

func (m *Manager) Providers(ctx context.Context) ([]Provider, error) {
	rows, err := m.conn.QueryContext(ctx, `
		SELECT id_proveedor, nombre, telefono, correo, direccion
		FROM proveedor
		ORDER BY nombre
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]Provider, 0)
	for rows.Next() {
		var provider Provider
		if err := rows.Scan(&provider.IDProveedor, &provider.Nombre, &provider.Telefono, &provider.Correo, &provider.Direccion); err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return providers, nil
}

func (m *Manager) Provider(ctx context.Context, id int) (*Provider, error) {
	row := m.conn.QueryRowContext(ctx, `
		SELECT id_proveedor, nombre, telefono, correo, direccion
		FROM proveedor
		WHERE id_proveedor = $1
	`, id)

	var provider Provider
	if err := row.Scan(&provider.IDProveedor, &provider.Nombre, &provider.Telefono, &provider.Correo, &provider.Direccion); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &provider, nil
}

func (m *Manager) CreateProvider(ctx context.Context, input ProviderWrite) (*Provider, error) {
	if err := validateProviderWrite(input); err != nil {
		return nil, err
	}

	row := m.conn.QueryRowContext(ctx, `
		INSERT INTO proveedor (nombre, telefono, correo, direccion)
		VALUES ($1, $2, $3, $4)
		RETURNING id_proveedor, nombre, telefono, correo, direccion
	`, strings.TrimSpace(input.Nombre), strings.TrimSpace(input.Telefono), strings.TrimSpace(input.Correo), strings.TrimSpace(input.Direccion))

	var provider Provider
	if err := row.Scan(&provider.IDProveedor, &provider.Nombre, &provider.Telefono, &provider.Correo, &provider.Direccion); err != nil {
		return nil, err
	}
	return &provider, nil
}

func (m *Manager) UpdateProvider(ctx context.Context, id int, input ProviderWrite) (*Provider, error) {
	if err := validateProviderWrite(input); err != nil {
		return nil, err
	}

	row := m.conn.QueryRowContext(ctx, `
		UPDATE proveedor
		SET nombre = $1,
		    telefono = $2,
		    correo = $3,
		    direccion = $4
		WHERE id_proveedor = $5
		RETURNING id_proveedor, nombre, telefono, correo, direccion
	`, strings.TrimSpace(input.Nombre), strings.TrimSpace(input.Telefono), strings.TrimSpace(input.Correo), strings.TrimSpace(input.Direccion), id)

	var provider Provider
	if err := row.Scan(&provider.IDProveedor, &provider.Nombre, &provider.Telefono, &provider.Correo, &provider.Direccion); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &provider, nil
}

func (m *Manager) DestroyProvider(ctx context.Context, id int) error {
	result, err := m.conn.ExecContext(ctx, `
		DELETE FROM proveedor
		WHERE id_proveedor = $1
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
	return nil
}

func (m *Manager) Employees(ctx context.Context) ([]Employee, error) {
	rows, err := m.conn.QueryContext(ctx, `
		SELECT id_empleado, nombre, estado, correo
		FROM empleado
		ORDER BY nombre
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]Employee, 0)
	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.IDEmpleado, &employee.Nombre, &employee.Estado, &employee.Correo); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return employees, nil
}

func (m *Manager) Clients(ctx context.Context) ([]Client, error) {
	rows, err := m.conn.QueryContext(ctx, `
		SELECT id_cliente, nombre, telefono, correo
		FROM cliente
		ORDER BY nombre
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]Client, 0)
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.IDCliente, &client.Nombre, &client.Telefono, &client.Correo); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clients, nil
}

func (m *Manager) Compras(ctx context.Context) ([]Compra, error) {
	rows, err := m.conn.QueryContext(ctx, `
		SELECT c.id_compra,
		       c.id_empleado,
		       e.nombre AS nombre_empleado,
		       c.id_cliente,
		       cl.nombre AS nombre_cliente,
		       c.fecha_compra,
		       COALESCE(MIN(pc.id_producto), 0) AS id_producto,
		       COALESCE(string_agg(p.nombre, ', ' ORDER BY p.nombre), 'Sin productos') AS productos,
		       c.total_compra
		FROM compra c
		JOIN empleado e ON e.id_empleado = c.id_empleado
		JOIN cliente cl ON cl.id_cliente = c.id_cliente
		LEFT JOIN producto_compra pc ON pc.id_compra = c.id_compra
		LEFT JOIN producto p ON p.id_producto = pc.id_producto
		GROUP BY c.id_compra, c.id_empleado, e.nombre, c.id_cliente, cl.nombre, c.fecha_compra, c.total_compra
		ORDER BY c.id_compra
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	compras := make([]Compra, 0)
	for rows.Next() {
		compra, err := scanCompra(rows)
		if err != nil {
			return nil, err
		}
		compras = append(compras, compra)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return compras, nil
}

func (m *Manager) Compra(ctx context.Context, id int) (*Compra, error) {
	return m.compraByID(ctx, m.conn, id)
}

func (m *Manager) CreateCompra(ctx context.Context, input CompraWrite) (*Compra, error) {
	if err := validateCompraWrite(input); err != nil {
		return nil, err
	}

	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO compra (id_empleado, id_cliente, fecha_compra, total_compra)
		VALUES ($1, $2, $3, $4)
		RETURNING id_compra
	`, input.IDEmpleado, input.IDCliente, input.FechaCompra, input.TotalCompra)

	var id int
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO producto_compra (id_compra, id_producto, cantidad_producto)
		VALUES ($1, $2, 1)
	`, id, input.IDProducto); err != nil {
		return nil, err
	}

	compra, err := m.compraByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return compra, nil
}

func (m *Manager) UpdateCompra(ctx context.Context, id int, input CompraWrite) (*Compra, error) {
	if err := validateCompraWrite(input); err != nil {
		return nil, err
	}

	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(ctx, `
		UPDATE compra
		SET id_empleado = $1,
		    id_cliente = $2,
		    fecha_compra = $3,
		    total_compra = $4
		WHERE id_compra = $5
	`, input.IDEmpleado, input.IDCliente, input.FechaCompra, input.TotalCompra, id)
	if err != nil {
		return nil, fmt.Errorf("transaction rollback after update error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("transaction rollback after rows affected error: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM producto_compra
		WHERE id_compra = $1
	`, id); err != nil {
		return nil, fmt.Errorf("transaction rollback after detail delete error: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO producto_compra (id_compra, id_producto, cantidad_producto)
		VALUES ($1, $2, 1)
	`, id, input.IDProducto); err != nil {
		return nil, fmt.Errorf("transaction rollback after detail insert error: %w", err)
	}

	compra, err := m.compraByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("transaction rollback after reload error: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}
	committed = true
	return compra, nil
}

func (m *Manager) DestroyCompra(ctx context.Context, id int) error {
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM producto_compra
		WHERE id_compra = $1
	`, id); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		DELETE FROM compra
		WHERE id_compra = $1
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

type compraScanner interface {
	Scan(dest ...any) error
}

type compraQueryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func scanProduct(scanner productScanner) (Product, error) {
	var product Product
	var imagen sql.NullString

	err := scanner.Scan(
		&product.IDProducto,
		&product.IDCategoria,
		&product.IDProveedor,
		&product.Categoria,
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

func scanCompra(scanner compraScanner) (Compra, error) {
	var compra Compra
	var fecha time.Time
	var productos sql.NullString

	err := scanner.Scan(
		&compra.IDCompra,
		&compra.IDEmpleado,
		&compra.NombreEmpleado,
		&compra.IDCliente,
		&compra.NombreCliente,
		&fecha,
		&compra.IDProducto,
		&productos,
		&compra.TotalCompra,
	)
	if err != nil {
		return Compra{}, err
	}

	compra.FechaCompra = fecha.Format("2006-01-02")
	if productos.Valid && productos.String != "" {
		compra.Productos = productos.String
	} else {
		compra.Productos = "Sin productos"
	}

	return compra, nil
}

func (m *Manager) compraByID(ctx context.Context, queryer compraQueryer, id int) (*Compra, error) {
	row := queryer.QueryRowContext(ctx, `
		SELECT c.id_compra,
		       c.id_empleado,
		       e.nombre AS nombre_empleado,
		       c.id_cliente,
		       cl.nombre AS nombre_cliente,
		       c.fecha_compra,
		       COALESCE(MIN(pc.id_producto), 0) AS id_producto,
		       COALESCE(string_agg(p.nombre, ', ' ORDER BY p.nombre), 'Sin productos') AS productos,
		       c.total_compra
		FROM compra c
		JOIN empleado e ON e.id_empleado = c.id_empleado
		JOIN cliente cl ON cl.id_cliente = c.id_cliente
		LEFT JOIN producto_compra pc ON pc.id_compra = c.id_compra
		LEFT JOIN producto p ON p.id_producto = pc.id_producto
		WHERE c.id_compra = $1
		GROUP BY c.id_compra, c.id_empleado, e.nombre, c.id_cliente, cl.nombre, c.fecha_compra, c.total_compra
	`, id)

	compra, err := scanCompra(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &compra, nil
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

func validateCategoryWrite(input CategoryWrite) error {
	if strings.TrimSpace(input.Nombre) == "" {
		return fmt.Errorf("%w: nombre is required", ErrInvalidInput)
	}
	return nil
}

func validateProviderWrite(input ProviderWrite) error {
	switch {
	case strings.TrimSpace(input.Nombre) == "":
		return fmt.Errorf("%w: nombre is required", ErrInvalidInput)
	case strings.TrimSpace(input.Telefono) == "":
		return fmt.Errorf("%w: telefono is required", ErrInvalidInput)
	case strings.TrimSpace(input.Correo) == "":
		return fmt.Errorf("%w: correo is required", ErrInvalidInput)
	case !strings.Contains(input.Correo, "@"):
		return fmt.Errorf("%w: correo must be a valid email", ErrInvalidInput)
	case strings.TrimSpace(input.Direccion) == "":
		return fmt.Errorf("%w: direccion is required", ErrInvalidInput)
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

func validateCompraWrite(input CompraWrite) error {
	switch {
	case input.IDEmpleado <= 0:
		return fmt.Errorf("%w: id_empleado must be greater than 0", ErrInvalidInput)
	case input.IDCliente <= 0:
		return fmt.Errorf("%w: id_cliente must be greater than 0", ErrInvalidInput)
	case strings.TrimSpace(input.FechaCompra) == "":
		return fmt.Errorf("%w: fecha_compra is required", ErrInvalidInput)
	case input.IDProducto <= 0:
		return fmt.Errorf("%w: id_producto must be greater than 0", ErrInvalidInput)
	case input.TotalCompra < 0:
		return fmt.Errorf("%w: total_compra cannot be negative", ErrInvalidInput)
	}

	if _, err := time.Parse("2006-01-02", input.FechaCompra); err != nil {
		return fmt.Errorf("%w: fecha_compra must use YYYY-MM-DD format", ErrInvalidInput)
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
