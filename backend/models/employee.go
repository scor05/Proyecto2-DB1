package models

type Employee struct {
	IDEmpleado int    `gorm:"column:id_empleado;primaryKey" json:"id_empleado"`
	Nombre     string `gorm:"column:nombre" json:"nombre"`
	Estado     string `gorm:"column:estado" json:"estado"`
	Correo     string `gorm:"column:correo" json:"correo"`
}

func (Employee) TableName() string {
	return "empleado"
}

type EmployeeWrite struct {
	Nombre string `json:"nombre"`
	Estado string `json:"estado"`
	Correo string `json:"correo"`
}
