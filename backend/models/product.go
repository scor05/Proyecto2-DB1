package models

type Product struct {
	IDProducto  int     `gorm:"column:id_producto;primaryKey" json:"id_producto"`
	IDCategoria int     `gorm:"column:id_categoria" json:"id_categoria"`
	IDProveedor int     `gorm:"column:id_proveedor" json:"id_proveedor"`
	Categoria   string  `gorm:"column:categoria;->" json:"categoria"`
	Precio      float64 `gorm:"column:precio" json:"precio"`
	Stock       int     `gorm:"column:stock" json:"stock"`
	Nombre      string  `gorm:"column:nombre" json:"nombre"`
	Imagen      *string `gorm:"column:imagen" json:"imagen"`
	Descripcion string  `gorm:"column:descripcion" json:"descripcion"`
}

func (Product) TableName() string {
	return "producto"
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
