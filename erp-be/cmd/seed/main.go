package main

import (
	"log"

	"erp-kafeinart/config"
	"erp-kafeinart/database"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Gagal menghubungkan database: %v", err)
	}

	database.SeedDatabase(db)
	log.Println("[SUCCESS] Seluruh seeder database berhasil dijalankan!")
}
