package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// PermissionList adalah tipe slice string dengan implementasi Scanner dan Valuer
// agar dapat disimpan sebagai string JSON di SQLite dan dibaca sebagai array di JSON response.
type PermissionList []string

// Value mengonversi slice string ke JSON string untuk penyimpanan database.
func (p PermissionList) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "[]", nil
	}
	return json.Marshal(p)
}

// Scan mengurai data dari database menjadi slice string.
func (p *PermissionList) Scan(value interface{}) error {
	if value == nil {
		*p = PermissionList{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("tipe data tidak valid untuk PermissionList")
		}
		bytes = []byte(str)
	}

	return json.Unmarshal(bytes, p)
}

// User merepresentasikan entitas pengguna dalam sistem ERP.
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Email       string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Role        string         `gorm:"size:50;not null;default:'staff'" json:"role"`
	Permissions PermissionList `gorm:"type:text" json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
