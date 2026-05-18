package models

type Category struct {
	IDCategoria int    `gorm:"column:id_categoria;primaryKey" json:"id_categoria"`
	Nombre      string `gorm:"column:nombre" json:"nombre"`
}

func (Category) TableName() string {
	return "categoria"
}

type CategoryWrite struct {
	Nombre string `json:"nombre"`
}
