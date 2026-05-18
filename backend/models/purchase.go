package models

type Compra struct {
	IDCompra       int     `gorm:"column:id_compra;primaryKey" json:"id_compra"`
	IDEmpleado     int     `gorm:"column:id_empleado" json:"id_empleado"`
	NombreEmpleado string  `gorm:"column:nombre_empleado;->" json:"nombre_empleado"`
	IDCliente      int     `gorm:"column:id_cliente" json:"id_cliente"`
	NombreCliente  string  `gorm:"column:nombre_cliente;->" json:"nombre_cliente"`
	FechaCompra    string  `gorm:"column:fecha_compra" json:"fecha_compra"`
	IDProducto     int     `gorm:"column:id_producto;->" json:"id_producto"`
	Productos      string  `gorm:"column:productos;->" json:"productos"`
	TotalCompra    float64 `gorm:"column:total_compra" json:"total_compra"`
}

func (Compra) TableName() string {
	return "compra"
}

type CompraWrite struct {
	IDEmpleado  int     `json:"id_empleado"`
	IDCliente   int     `json:"id_cliente"`
	FechaCompra string  `json:"fecha_compra"`
	IDProducto  int     `json:"id_producto"`
	TotalCompra float64 `json:"total_compra"`
}
