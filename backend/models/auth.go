package models

type AuthUser struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Correo string `json:"correo"`
	Rol    string `json:"rol"`
}

type Gerente struct {
	IDGerente int    `gorm:"column:id_gerente;primaryKey" json:"id_gerente"`
	Nombre    string `gorm:"column:nombre" json:"nombre"`
	Correo    string `gorm:"column:correo" json:"correo"`
}

func (Gerente) TableName() string {
	return "gerente"
}

type SuperAdmin struct {
	IDSuperAdmin int    `gorm:"column:id_superadmin;primaryKey" json:"id_superadmin"`
	Nombre       string `gorm:"column:nombre" json:"nombre"`
	Correo       string `gorm:"column:correo" json:"correo"`
}

func (SuperAdmin) TableName() string {
	return "superadmin"
}
