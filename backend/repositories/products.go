package repositories

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"proyecto2/backend/models"
)

func (r *Repository) IndexProducts(ctx context.Context) ([]models.Product, error) {
	products := make([]models.Product, 0)
	err := r.db.WithContext(ctx).
		Table("producto p").
		Select("p.id_producto, p.id_categoria, p.id_proveedor, c.nombre AS categoria, p.precio, p.stock, p.nombre, p.imagen, p.descripcion").
		Joins("JOIN categoria c ON c.id_categoria = p.id_categoria").
		Order("p.id_producto").
		Scan(&products).Error
	return products, err
}

func (r *Repository) ShowProduct(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Table("producto p").
		Select("p.id_producto, p.id_categoria, p.id_proveedor, c.nombre AS categoria, p.precio, p.stock, p.nombre, p.imagen, p.descripcion").
		Joins("JOIN categoria c ON c.id_categoria = p.id_categoria").
		Where("p.id_producto = ?", id).
		Take(&product).Error
	if err != nil {
		return nil, notFound(err)
	}
	return &product, nil
}

func (r *Repository) CreateProduct(ctx context.Context, input models.ProductWrite) (*models.Product, error) {
	if err := validateProductWrite(input); err != nil {
		return nil, err
	}

	product := models.Product{
		IDCategoria: input.IDCategoria,
		IDProveedor: input.IDProveedor,
		Precio:      input.Precio,
		Stock:       input.Stock,
		Nombre:      strings.TrimSpace(input.Nombre),
		Imagen:      input.Imagen,
		Descripcion: strings.TrimSpace(input.Descripcion),
	}

	if err := r.db.WithContext(ctx).Create(&product).Error; err != nil {
		return nil, err
	}
	return r.ShowProduct(ctx, product.IDProducto)
}

func (r *Repository) UpdateProduct(ctx context.Context, id int, input models.ProductWrite) (*models.Product, error) {
	if err := validateProductWrite(input); err != nil {
		return nil, err
	}

	updates := map[string]any{
		"id_categoria": input.IDCategoria,
		"id_proveedor": input.IDProveedor,
		"precio":       input.Precio,
		"stock":        input.Stock,
		"nombre":       strings.TrimSpace(input.Nombre),
		"imagen":       input.Imagen,
		"descripcion":  strings.TrimSpace(input.Descripcion),
	}

	result := r.db.WithContext(ctx).Model(&models.Product{}).Where("id_producto = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return r.ShowProduct(ctx, id)
}

func (r *Repository) PatchProduct(ctx context.Context, id int, input models.ProductPatch) (*models.Product, error) {
	if err := validateProductPatch(input); err != nil {
		return nil, err
	}

	updates := make(map[string]any)
	if input.IDCategoria != nil {
		updates["id_categoria"] = *input.IDCategoria
	}
	if input.IDProveedor != nil {
		updates["id_proveedor"] = *input.IDProveedor
	}
	if input.Precio != nil {
		updates["precio"] = *input.Precio
	}
	if input.Stock != nil {
		updates["stock"] = *input.Stock
	}
	if input.Nombre != nil {
		updates["nombre"] = strings.TrimSpace(*input.Nombre)
	}
	if input.Imagen != nil {
		updates["imagen"] = *input.Imagen
	}
	if input.Descripcion != nil {
		updates["descripcion"] = strings.TrimSpace(*input.Descripcion)
	}

	result := r.db.WithContext(ctx).Model(&models.Product{}).Where("id_producto = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return r.ShowProduct(ctx, id)
}

func (r *Repository) DestroyProduct(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.ProductCompra{}, "id_producto = ?", id).Error; err != nil {
			return err
		}

		result := tx.Delete(&models.Product{}, "id_producto = ?", id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

func (r *Repository) LoginEmployee(ctx context.Context, correo string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).
		Where("lower(correo) = lower(?)", strings.TrimSpace(correo)).
		First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &employee, nil
}
