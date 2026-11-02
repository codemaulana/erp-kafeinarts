package database

import (
	"log"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// SeedProducts menginisialisasi sampel katalog produk inventory.
func SeedProducts(db *gorm.DB) {
	var count int64
	db.Model(&model.Product{}).Count(&count)
	if count > 0 {
		return
	}

	products := []model.Product{
		{
			SKU:         "KAF-BEANS-001",
			Name:        "House Blend Espresso Beans 1kg",
			Stock:       45,
			Price:       180000,
			Description: "Racikan arabika & robusta 70:30 khas Kafeinarts untuk racikan espresso seimbang.",
		},
		{
			SKU:         "KAF-BEANS-002",
			Name:        "Single Origin Gayo Honey 250g",
			Stock:       7,
			Price:       95000,
			Description: "Biji kopi proses honey dengan tasting notes peach, brown sugar, dan lingering jasmine.",
		},
		{
			SKU:         "KAF-SYRUP-001",
			Name:        "Artisan Vanilla Syrup 750ml",
			Stock:       22,
			Price:       125000,
			Description: "Sirup ekstrak vanilla murni untuk menu mocktail dan flavoured latte.",
		},
		{
			SKU:         "KAF-CUP-001",
			Name:        "Takeaway Paper Cup 8oz (Dus @1000)",
			Stock:       5,
			Price:       350000,
			Description: "Gelas kertas food-grade ramah lingkungan dengan logo grafis Kafeinarts.",
		},
		{
			SKU:         "KAF-MERCH-001",
			Name:        "Kafeinarts Studio Canvas Tote Bag",
			Stock:       30,
			Price:       85000,
			Description: "Tote bag kanvas tebal 14oz berdesain tipografi eksklusif edisi studio seni.",
		},
	}

	if err := db.Create(&products).Error; err != nil {
		log.Printf("[ERROR] Gagal membuat produk seeder: %v", err)
	} else {
		log.Printf("[SEEDER] %d produk default berhasil ditambahkan", len(products))
	}
}
