package models

type Employee struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null"`
	Title     string `json:"title" gorm:"not null"`
	ManagerID *int   `json:"manager_id" gorm:"index"`
}

func (Employee) TableName() string {
	return "employees"
}
