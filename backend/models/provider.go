package models

type Provider struct {
	IDProveedor int    `gorm:"column:id_proveedor;primaryKey" json:"id_proveedor"`
	Nombre      string `gorm:"column:nombre" json:"nombre"`
	Telefono    string `gorm:"column:telefono" json:"telefono"`
	Correo      string `gorm:"column:correo" json:"correo"`
	Direccion   string `gorm:"column:direccion" json:"direccion"`
}

func (Provider) TableName() string {
	return "proveedor"
}

type ProviderWrite struct {
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
	Direccion string `json:"direccion"`
}
