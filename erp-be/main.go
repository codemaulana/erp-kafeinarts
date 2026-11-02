package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"erp-kafeinart/config"
	"erp-kafeinart/database"
	"erp-kafeinart/middleware"
	"erp-kafeinart/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// 1. Muat Konfigurasi Lingkungan
	cfg := config.LoadConfig()

	// 2. Hubungkan Database SQLite & AutoMigrate
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Gagal menghubungkan database: %v", err)
	}

	// 3. Jalankan Database Seeder (User admin default, sample produk, finance, pegawai)
	database.SeedDatabase(db)

	// 4. Inisialisasi Instance Fiber v3
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	// 5. Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} | ${method} ${path}\n",
	}))
	app.Use(middleware.SetupCORS(cfg))

	// 6. Daftarkan Routes
	routes.SetupRoutes(app, db, cfg)

	// 7. Saluran untuk Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("=====================================================")
		log.Printf("   %s berjalan di http://localhost%s", cfg.AppName, addr)
		log.Printf("   API Base URL: http://localhost%s/api/v1", addr)
		log.Printf("   Default Login: admin@kafeinarts.com / password123")
		log.Printf("=====================================================")

		if err := app.Listen(addr, fiber.ListenConfig{
			DisableStartupMessage: true,
		}); err != nil {
			log.Printf("[INFO] Server berhenti: %v", err)
		}
	}()

	// Menunggu sinyal interrupt untuk shutdown aplikasi secara bersih
	<-stop
	log.Println("\n[INFO] Mematikan server secara graceful...")
	if err := app.Shutdown(); err != nil {
		log.Printf("[ERROR] Gagal mematikan server: %v", err)
	}
	log.Println("[SUCCESS] Server berhasil dinonaktifkan.")
}