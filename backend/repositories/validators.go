package repositories

import (
	"fmt"
	"strings"
	"time"

	"proyecto2/backend/models"
)

func validateProductWrite(input models.ProductWrite) error {
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

func validateCategoryWrite(input models.CategoryWrite) error {
	if strings.TrimSpace(input.Nombre) == "" {
		return fmt.Errorf("%w: nombre is required", ErrInvalidInput)
	}
	return nil
}

func validateProviderWrite(input models.ProviderWrite) error {
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

func validateEmployeeWrite(input models.EmployeeWrite) error {
	estado := strings.TrimSpace(input.Estado)
	switch {
	case strings.TrimSpace(input.Nombre) == "":
		return fmt.Errorf("%w: nombre is required", ErrInvalidInput)
	case estado != "activo" && estado != "inactivo":
		return fmt.Errorf("%w: estado must be activo or inactivo", ErrInvalidInput)
	case strings.TrimSpace(input.Correo) == "":
		return fmt.Errorf("%w: correo is required", ErrInvalidInput)
	case !strings.Contains(input.Correo, "@"):
		return fmt.Errorf("%w: correo must be a valid email", ErrInvalidInput)
	default:
		return nil
	}
}

func validateClientWrite(input models.ClientWrite) error {
	switch {
	case strings.TrimSpace(input.Nombre) == "":
		return fmt.Errorf("%w: nombre is required", ErrInvalidInput)
	case strings.TrimSpace(input.Telefono) == "":
		return fmt.Errorf("%w: telefono is required", ErrInvalidInput)
	case strings.TrimSpace(input.Correo) == "":
		return fmt.Errorf("%w: correo is required", ErrInvalidInput)
	case !strings.Contains(input.Correo, "@"):
		return fmt.Errorf("%w: correo must be a valid email", ErrInvalidInput)
	default:
		return nil
	}
}

func validateProductPatch(input models.ProductPatch) error {
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

func validateCompraWrite(input models.CompraWrite) error {
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
