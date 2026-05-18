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
		var id int
		if err := tx.Raw(`
			INSERT INTO compra (id_empleado, id_cliente, fecha_compra, total_compra)
			VALUES (?, ?, ?, ?)
			RETURNING id_compra
		`, input.IDEmpleado, input.IDCliente, input.FechaCompra, input.TotalCompra).Scan(&id).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			INSERT INTO producto_compra (id_compra, id_producto, cantidad_producto)
			VALUES (?, ?, 1)
		`, id, input.IDProducto).Error; err != nil {
			return err
		}

		loaded, err := compraByID(ctx, tx, id)
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
		result := tx.Exec(`
			UPDATE compra
			SET id_empleado = ?,
			    id_cliente = ?,
			    fecha_compra = ?,
			    total_compra = ?
			WHERE id_compra = ?
		`, input.IDEmpleado, input.IDCliente, input.FechaCompra, input.TotalCompra, id)
		if result.Error != nil {
			return fmt.Errorf("transaction rollback after update error: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}

		if err := tx.Exec("DELETE FROM producto_compra WHERE id_compra = ?", id).Error; err != nil {
			return fmt.Errorf("transaction rollback after detail delete error: %w", err)
		}
		if err := tx.Exec(`
			INSERT INTO producto_compra (id_compra, id_producto, cantidad_producto)
			VALUES (?, ?, 1)
		`, id, input.IDProducto).Error; err != nil {
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
		if err := tx.Exec("DELETE FROM producto_compra WHERE id_compra = ?", id).Error; err != nil {
			return err
		}

		result := tx.Exec("DELETE FROM compra WHERE id_compra = ?", id)
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
