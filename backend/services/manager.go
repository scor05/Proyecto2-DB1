package services

import (
	"context"

	"proyecto2/backend/db"
	"proyecto2/backend/models"
	"proyecto2/backend/repositories"
)

var (
	ErrNotFound     = repositories.ErrNotFound
	ErrInvalidInput = repositories.ErrInvalidInput
)

type Product = models.Product
type ProductWrite = models.ProductWrite
type ProductPatch = models.ProductPatch
type Employee = models.Employee
type EmployeeWrite = models.EmployeeWrite
type Category = models.Category
type CategoryWrite = models.CategoryWrite
type Provider = models.Provider
type ProviderWrite = models.ProviderWrite
type Client = models.Client
type ClientWrite = models.ClientWrite
type Compra = models.Compra
type CompraWrite = models.CompraWrite

type Manager struct {
	conn *db.Connection
	repo *repositories.Repository
}

func NewManager() (*Manager, error) {
	conn, err := db.NewConnection()
	if err != nil {
		return nil, err
	}
	return &Manager{
		conn: conn,
		repo: repositories.New(conn.DB),
	}, nil
}

func (m *Manager) Close() error {
	return m.conn.Close()
}

func (m *Manager) Ping(ctx context.Context) error {
	return m.conn.Ping(ctx)
}

func (m *Manager) Index(ctx context.Context) ([]models.Product, error) {
	return m.repo.IndexProducts(ctx)
}

func (m *Manager) Show(ctx context.Context, id int) (*models.Product, error) {
	return m.repo.ShowProduct(ctx, id)
}

func (m *Manager) Create(ctx context.Context, input models.ProductWrite) (*models.Product, error) {
	return m.repo.CreateProduct(ctx, input)
}

func (m *Manager) Update(ctx context.Context, id int, input models.ProductWrite) (*models.Product, error) {
	return m.repo.UpdateProduct(ctx, id, input)
}

func (m *Manager) Patch(ctx context.Context, id int, input models.ProductPatch) (*models.Product, error) {
	return m.repo.PatchProduct(ctx, id, input)
}

func (m *Manager) Destroy(ctx context.Context, id int) error {
	return m.repo.DestroyProduct(ctx, id)
}

func (m *Manager) LoginEmployee(ctx context.Context, correo string) (*models.Employee, error) {
	return m.repo.LoginEmployee(ctx, correo)
}

func (m *Manager) Categories(ctx context.Context) ([]models.Category, error) {
	return m.repo.Categories(ctx)
}

func (m *Manager) Category(ctx context.Context, id int) (*models.Category, error) {
	return m.repo.Category(ctx, id)
}

func (m *Manager) CreateCategory(ctx context.Context, input models.CategoryWrite) (*models.Category, error) {
	return m.repo.CreateCategory(ctx, input)
}

func (m *Manager) UpdateCategory(ctx context.Context, id int, input models.CategoryWrite) (*models.Category, error) {
	return m.repo.UpdateCategory(ctx, id, input)
}

func (m *Manager) DestroyCategory(ctx context.Context, id int) error {
	return m.repo.DestroyCategory(ctx, id)
}

func (m *Manager) Providers(ctx context.Context) ([]models.Provider, error) {
	return m.repo.Providers(ctx)
}

func (m *Manager) Provider(ctx context.Context, id int) (*models.Provider, error) {
	return m.repo.Provider(ctx, id)
}

func (m *Manager) CreateProvider(ctx context.Context, input models.ProviderWrite) (*models.Provider, error) {
	return m.repo.CreateProvider(ctx, input)
}

func (m *Manager) UpdateProvider(ctx context.Context, id int, input models.ProviderWrite) (*models.Provider, error) {
	return m.repo.UpdateProvider(ctx, id, input)
}

func (m *Manager) DestroyProvider(ctx context.Context, id int) error {
	return m.repo.DestroyProvider(ctx, id)
}

func (m *Manager) Employees(ctx context.Context) ([]models.Employee, error) {
	return m.repo.Employees(ctx)
}

func (m *Manager) Employee(ctx context.Context, id int) (*models.Employee, error) {
	return m.repo.Employee(ctx, id)
}

func (m *Manager) CreateEmployee(ctx context.Context, input models.EmployeeWrite) (*models.Employee, error) {
	return m.repo.CreateEmployee(ctx, input)
}

func (m *Manager) UpdateEmployee(ctx context.Context, id int, input models.EmployeeWrite) (*models.Employee, error) {
	return m.repo.UpdateEmployee(ctx, id, input)
}

func (m *Manager) DestroyEmployee(ctx context.Context, id int) error {
	return m.repo.DestroyEmployee(ctx, id)
}

func (m *Manager) Clients(ctx context.Context) ([]models.Client, error) {
	return m.repo.Clients(ctx)
}

func (m *Manager) Client(ctx context.Context, id int) (*models.Client, error) {
	return m.repo.Client(ctx, id)
}

func (m *Manager) CreateClient(ctx context.Context, input models.ClientWrite) (*models.Client, error) {
	return m.repo.CreateClient(ctx, input)
}

func (m *Manager) UpdateClient(ctx context.Context, id int, input models.ClientWrite) (*models.Client, error) {
	return m.repo.UpdateClient(ctx, id, input)
}

func (m *Manager) DestroyClient(ctx context.Context, id int) error {
	return m.repo.DestroyClient(ctx, id)
}

func (m *Manager) Compras(ctx context.Context) ([]models.Compra, error) {
	return m.repo.Compras(ctx)
}

func (m *Manager) Compra(ctx context.Context, id int) (*models.Compra, error) {
	return m.repo.Compra(ctx, id)
}

func (m *Manager) CreateCompra(ctx context.Context, input models.CompraWrite) (*models.Compra, error) {
	return m.repo.CreateCompra(ctx, input)
}

func (m *Manager) UpdateCompra(ctx context.Context, id int, input models.CompraWrite) (*models.Compra, error) {
	return m.repo.UpdateCompra(ctx, id, input)
}

func (m *Manager) DestroyCompra(ctx context.Context, id int) error {
	return m.repo.DestroyCompra(ctx, id)
}
