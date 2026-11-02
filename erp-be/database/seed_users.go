package database

import (
	"log"

	"erp-kafeinart/model"
	"erp-kafeinart/utils"

	"gorm.io/gorm"
)

// SeedUsers membuat akun default untuk setiap modul jika belum ada di database.
func SeedUsers(db *gorm.DB) {
	defaultUsers := []struct {
		Name        string
		Email       string
		Role        string
		Permissions []string
	}{
		{
			Name:        "Super Admin",
			Email:       "admin@kafeinarts.com",
			Role:        "superadmin",
			Permissions: []string{"*"},
		},
		{
			Name:        "Admin Inventaris",
			Email:       "inventory@kafeinarts.com",
			Role:        "inventory",
			Permissions: []string{"inventory"},
		},
		{
			Name:        "Admin Finance",
			Email:       "finance@kafeinarts.com",
			Role:        "finance",
			Permissions: []string{"finance"},
		},
		{
			Name:        "Admin Human Resource",
			Email:       "hr@kafeinarts.com",
			Role:        "hr",
			Permissions: []string{"hr"},
		},
		{
			Name:        "Admin Supply Chain",
			Email:       "sc@kafeinarts.com",
			Role:        "supply_chain",
			Permissions: []string{"supply_chain"},
		},
	}

	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		log.Printf("[ERROR] Gagal hashing password seeder: %v", err)
		return
	}

	for _, u := range defaultUsers {
		var existing model.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err != nil {
			newUser := model.User{
				Name:        u.Name,
				Email:       u.Email,
				Password:    hashedPassword,
				Role:        u.Role,
				Permissions: model.PermissionList(u.Permissions),
			}
			if err := db.Create(&newUser).Error; err != nil {
				log.Printf("[ERROR] Gagal membuat user seeder %s: %v", u.Email, err)
			} else {
				log.Printf("[SEEDER] Default user berhasil dibuat: %s (role: %s, pass: password123)", u.Email, u.Role)
			}
		}
	}
}
