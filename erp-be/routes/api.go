package routes

import (
	"erp-kafeinart/config"
	"erp-kafeinart/handler"
	"erp-kafeinart/middleware"
	"erp-kafeinart/repository"
	"erp-kafeinart/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"gorm.io/gorm"
)

// SetupRoutes mendaftarkan seluruh endpoint API dan middleware pada aplikasi Fiber v3.
func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// ── 1. Inisialisasi Repositories (Data Access Layer) ─────────────────────────
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)

	financeMasterRepo := repository.NewFinanceMasterRepository(db)
	financeTrxRepo := repository.NewFinanceTransactionRepository(db)
	hrRepo := repository.NewHrRepository(db)
	scRepo := repository.NewScRepository(db)

	// ── 2. Inisialisasi Services (Business Logic Layer) ──────────────────────────
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	txService := service.NewTransactionService(txRepo, ledgerRepo)
	financeService := service.NewFinanceService(financeMasterRepo, financeTrxRepo)
	hrService := service.NewHrService(hrRepo)
	scService := service.NewScService(scRepo)

	// ── 3. Inisialisasi Handlers (Presentation Layer) ────────────────────────────
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	txHandler := handler.NewTransactionHandler(txService)
	financeHandler := handler.NewFinanceHandler(financeService)
	hrHandler := handler.NewHrHandler(hrService)
	scHandler := handler.NewScHandler(scService)

	// ── 4. Static Routes & Middleware ───────────────────────────────────────────
	app.Use("/uploads", static.New("./storage/uploads"))
	jwtAuth := middleware.Protected(cfg)

	// ── 5. Health & Root Check ──────────────────────────────────────────────────
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":     cfg.AppName,
			"version": "1.2.0",
			"status":  "running",
			"docs":    "/docs",
		})
	})
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// ── 6. Pendaftaran Sub-Routes Modular (/api/v1) ──────────────────────────────
	api := app.Group("/api/v1")

	api.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":     cfg.AppName,
			"version": "1.2.0",
			"status":  "running",
			"message": "Kafeinarts ERP REST API v1 aktif dan siap melayani permintaan",
			"endpoints": fiber.Map{
				"auth":      "/api/v1/auth/login",
				"users":     "/api/v1/users",
				"inventory": "/api/v1/inventory/products",
				"finance":   "/api/v1/finance",
				"hr":        "/api/v1/hr",
				"sc":        "/api/v1/sc",
			},
		})
	})

	RegisterAuthRoutes(api, app, authHandler, jwtAuth)
	RegisterUserRoutes(api, userHandler, jwtAuth)
	RegisterInventoryRoutes(api, productHandler, jwtAuth)
	RegisterFinanceRoutes(api, financeHandler, txHandler, jwtAuth)
	RegisterHrRoutes(api, hrHandler, jwtAuth)
	RegisterScRoutes(api, scHandler, jwtAuth)
}
