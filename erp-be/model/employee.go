package model

import "time"

// Employee merepresentasikan data staf / pegawai Kafeinarts Studio.
type Employee struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	NIP        string    `gorm:"size:50;uniqueIndex;not null" json:"nip"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Email      string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone      string    `gorm:"size:30" json:"phone"`
	Department string    `gorm:"size:100;not null" json:"department"`
	Position   string    `gorm:"size:100;not null" json:"position"`
	Status     string    `gorm:"size:30;not null;default:'active'" json:"status"` // "active", "probation", "inactive"
	JoinDate   time.Time `gorm:"not null" json:"join_date"`
	Salary     float64   `gorm:"type:decimal(14,2);not null;default:0" json:"salary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
