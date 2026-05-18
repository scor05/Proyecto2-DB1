package models

type ProductCompra struct {
	IDCompra         int `gorm:"column:id_compra;primaryKey" json:"id_compra"`
	IDProducto       int `gorm:"column:id_producto;primaryKey" json:"id_producto"`
	CantidadProducto int `gorm:"column:cantidad_producto" json:"cantidad_producto"`
}

func (ProductCompra) TableName() string {
	return "producto_compra"
}
