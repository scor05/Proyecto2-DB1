package models

type Client struct {
	IDCliente int    `gorm:"column:id_cliente;primaryKey" json:"id_cliente"`
	Nombre    string `gorm:"column:nombre" json:"nombre"`
	Telefono  string `gorm:"column:telefono" json:"telefono"`
	Correo    string `gorm:"column:correo" json:"correo"`
}

func (Client) TableName() string {
	return "cliente"
}

type ClientWrite struct {
	Nombre   string `json:"nombre"`
	Telefono string `json:"telefono"`
	Correo   string `json:"correo"`
}
