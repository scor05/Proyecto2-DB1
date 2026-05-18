package repositories

import (
	"context"
	"strings"

	"proyecto2/backend/models"
)

func (r *Repository) Categories(ctx context.Context) ([]models.Category, error) {
	categories := make([]models.Category, 0)
	err := r.db.WithContext(ctx).Order("nombre").Find(&categories).Error
	return categories, err
}

func (r *Repository) Category(ctx context.Context, id int) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).First(&category, "id_categoria = ?", id).Error; err != nil {
		return nil, notFound(err)
	}
	return &category, nil
}

func (r *Repository) CreateCategory(ctx context.Context, input models.CategoryWrite) (*models.Category, error) {
	if err := validateCategoryWrite(input); err != nil {
		return nil, err
	}

	category := models.Category{Nombre: strings.TrimSpace(input.Nombre)}
	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, id int, input models.CategoryWrite) (*models.Category, error) {
	if err := validateCategoryWrite(input); err != nil {
		return nil, err
	}

	result := r.db.WithContext(ctx).Model(&models.Category{}).
		Where("id_categoria = ?", id).
		Update("nombre", strings.TrimSpace(input.Nombre))
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return r.Category(ctx, id)
}

func (r *Repository) DestroyCategory(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Category{}, "id_categoria = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Providers(ctx context.Context) ([]models.Provider, error) {
	providers := make([]models.Provider, 0)
	err := r.db.WithContext(ctx).Order("nombre").Find(&providers).Error
	return providers, err
}

func (r *Repository) Provider(ctx context.Context, id int) (*models.Provider, error) {
	var provider models.Provider
	if err := r.db.WithContext(ctx).First(&provider, "id_proveedor = ?", id).Error; err != nil {
		return nil, notFound(err)
	}
	return &provider, nil
}

func (r *Repository) CreateProvider(ctx context.Context, input models.ProviderWrite) (*models.Provider, error) {
	if err := validateProviderWrite(input); err != nil {
		return nil, err
	}

	provider := models.Provider{
		Nombre:    strings.TrimSpace(input.Nombre),
		Telefono:  strings.TrimSpace(input.Telefono),
		Correo:    strings.TrimSpace(input.Correo),
		Direccion: strings.TrimSpace(input.Direccion),
	}
	if err := r.db.WithContext(ctx).Create(&provider).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *Repository) UpdateProvider(ctx context.Context, id int, input models.ProviderWrite) (*models.Provider, error) {
	if err := validateProviderWrite(input); err != nil {
		return nil, err
	}

	updates := map[string]any{
		"nombre":    strings.TrimSpace(input.Nombre),
		"telefono":  strings.TrimSpace(input.Telefono),
		"correo":    strings.TrimSpace(input.Correo),
		"direccion": strings.TrimSpace(input.Direccion),
	}
	result := r.db.WithContext(ctx).Model(&models.Provider{}).Where("id_proveedor = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return r.Provider(ctx, id)
}

func (r *Repository) DestroyProvider(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Provider{}, "id_proveedor = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Employees(ctx context.Context) ([]models.Employee, error) {
	employees := make([]models.Employee, 0)
	err := r.db.WithContext(ctx).Order("nombre").Find(&employees).Error
	return employees, err
}

func (r *Repository) Employee(ctx context.Context, id int) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.WithContext(ctx).First(&employee, "id_empleado = ?", id).Error; err != nil {
		return nil, notFound(err)
	}
	return &employee, nil
}

func (r *Repository) CreateEmployee(ctx context.Context, input models.EmployeeWrite) (*models.Employee, error) {
	if err := validateEmployeeWrite(input); err != nil {
		return nil, err
	}

	employee := models.Employee{
		Nombre: strings.TrimSpace(input.Nombre),
		Estado: strings.TrimSpace(input.Estado),
		Correo: strings.TrimSpace(input.Correo),
	}
	if err := r.db.WithContext(ctx).Create(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *Repository) UpdateEmployee(ctx context.Context, id int, input models.EmployeeWrite) (*models.Employee, error) {
	if err := validateEmployeeWrite(input); err != nil {
		return nil, err
	}

	updates := map[string]any{
		"nombre": strings.TrimSpace(input.Nombre),
		"estado": strings.TrimSpace(input.Estado),
		"correo": strings.TrimSpace(input.Correo),
	}
	result := r.db.WithContext(ctx).Model(&models.Employee{}).Where("id_empleado = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return r.Employee(ctx, id)
}

func (r *Repository) DestroyEmployee(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Employee{}, "id_empleado = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Clients(ctx context.Context) ([]models.Client, error) {
	clients := make([]models.Client, 0)
	err := r.db.WithContext(ctx).Order("nombre").Find(&clients).Error
	return clients, err
}

func (r *Repository) Client(ctx context.Context, id int) (*models.Client, error) {
	var client models.Client
	if err := r.db.WithContext(ctx).First(&client, "id_cliente = ?", id).Error; err != nil {
		return nil, notFound(err)
	}
	return &client, nil
}

func (r *Repository) CreateClient(ctx context.Context, input models.ClientWrite) (*models.Client, error) {
	if err := validateClientWrite(input); err != nil {
		return nil, err
	}

	client := models.Client{
		Nombre:   strings.TrimSpace(input.Nombre),
		Telefono: strings.TrimSpace(input.Telefono),
		Correo:   strings.TrimSpace(input.Correo),
	}
	if err := r.db.WithContext(ctx).Create(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *Repository) UpdateClient(ctx context.Context, id int, input models.ClientWrite) (*models.Client, error) {
	if err := validateClientWrite(input); err != nil {
		return nil, err
	}

	updates := map[string]any{
		"nombre":   strings.TrimSpace(input.Nombre),
		"telefono": strings.TrimSpace(input.Telefono),
		"correo":   strings.TrimSpace(input.Correo),
	}
	result := r.db.WithContext(ctx).Model(&models.Client{}).Where("id_cliente = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return r.Client(ctx, id)
}

func (r *Repository) DestroyClient(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Client{}, "id_cliente = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
