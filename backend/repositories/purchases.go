package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"proyecto2/backend/models"
)

const compraSelect = `
	SELECT c.id_compra,
	       c.id_empleado,
	       e.nombre AS nombre_empleado,
	       c.id_cliente,
	       cl.nombre AS nombre_cliente,
	       to_char(c.fecha_compra, 'YYYY-MM-DD') AS fecha_compra,
	       COALESCE(MIN(pc.id_producto), 0) AS id_producto,
	       COALESCE(string_agg(p.nombre, ', ' ORDER BY p.nombre), 'Sin productos') AS productos,
	       c.total_compra
	FROM compra c
	JOIN empleado e ON e.id_empleado = c.id_empleado
	JOIN cliente cl ON cl.id_cliente = c.id_cliente
	LEFT JOIN producto_compra pc ON pc.id_compra = c.id_compra
	LEFT JOIN producto p ON p.id_producto = pc.id_producto
`

const compraGroupOrder = `
	GROUP BY c.id_compra, c.id_empleado, e.nombre, c.id_cliente, cl.nombre, c.fecha_compra, c.total_compra
	ORDER BY c.id_compra
`

func (r *Repository) Compras(ctx context.Context) ([]models.Compra, error) {
	compras := make([]models.Compra, 0)
	err := r.db.WithContext(ctx).Raw(compraSelect + compraGroupOrder).Scan(&compras).Error
	return compras, err
}

func (r *Repository) Compra(ctx context.Context, id int) (*models.Compra, error) {
	return compraByID(ctx, r.db, id)
}

func (r *Repository) CreateCompra(ctx context.Context, input models.CompraWrite) (*models.Compra, error) {
	if err := validateCompraWrite(input); err != nil {
		return nil, err
	}

	var compra *models.Compra
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		purchase := models.Compra{
			IDEmpleado:  input.IDEmpleado,
			IDCliente:   input.IDCliente,
			FechaCompra: input.FechaCompra,
			TotalCompra: input.TotalCompra,
		}
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		item := models.ProductCompra{
			IDCompra:         purchase.IDCompra,
			IDProducto:       input.IDProducto,
			CantidadProducto: 1,
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}

		loaded, err := compraByID(ctx, tx, purchase.IDCompra)
		if err != nil {
			return err
		}
		compra = loaded
		return nil
	})
	if err != nil {
		return nil, err
	}
	return compra, nil
}

func (r *Repository) UpdateCompra(ctx context.Context, id int, input models.CompraWrite) (*models.Compra, error) {
	if err := validateCompraWrite(input); err != nil {
		return nil, err
	}

	var compra *models.Compra
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"id_empleado":  input.IDEmpleado,
			"id_cliente":   input.IDCliente,
			"fecha_compra": input.FechaCompra,
			"total_compra": input.TotalCompra,
		}
		result := tx.Model(&models.Compra{}).Where("id_compra = ?", id).Updates(updates)
		if result.Error != nil {
			return fmt.Errorf("transaction rollback after update error: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}

		if err := tx.Delete(&models.ProductCompra{}, "id_compra = ?", id).Error; err != nil {
			return fmt.Errorf("transaction rollback after detail delete error: %w", err)
		}

		item := models.ProductCompra{
			IDCompra:         id,
			IDProducto:       input.IDProducto,
			CantidadProducto: 1,
		}
		if err := tx.Create(&item).Error; err != nil {
			return fmt.Errorf("transaction rollback after detail insert error: %w", err)
		}

		loaded, err := compraByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("transaction rollback after reload error: %w", err)
		}
		compra = loaded
		return nil
	})
	if err != nil {
		return nil, err
	}
	return compra, nil
}

func (r *Repository) DestroyCompra(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.ProductCompra{}, "id_compra = ?", id).Error; err != nil {
			return err
		}

		result := tx.Delete(&models.Compra{}, "id_compra = ?", id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

func compraByID(ctx context.Context, db *gorm.DB, id int) (*models.Compra, error) {
	var compra models.Compra
	err := db.WithContext(ctx).Raw(compraSelect+`
		WHERE c.id_compra = ?
		GROUP BY c.id_compra, c.id_empleado, e.nombre, c.id_cliente, cl.nombre, c.fecha_compra, c.total_compra
	`, id).Scan(&compra).Error
	if err != nil {
		return nil, err
	}
	if compra.IDCompra == 0 {
		return nil, ErrNotFound
	}
	if compra.Productos == "" {
		compra.Productos = "Sin productos"
	}
	return &compra, nil
}
