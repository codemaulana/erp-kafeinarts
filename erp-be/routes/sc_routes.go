package routes

import (
	"erp-kafeinart/handler"

	"github.com/gofiber/fiber/v3"
)

// RegisterScRoutes mendaftarkan seluruh endpoint modul Rantai Pasok (Supply Chain).
func RegisterScRoutes(api fiber.Router, scHandler *handler.ScHandler, jwtAuth fiber.Handler) {
	sc := api.Group("/sc", jwtAuth)

	// ── A. Data Master ──────────────────────────────────────────────────────────
	sc.Get("/vendors", scHandler.GetVendors)
	sc.Get("/vendors/:id", scHandler.GetVendorByID)
	sc.Post("/vendors", scHandler.CreateVendor)
	sc.Put("/vendors/:id", scHandler.UpdateVendor)
	sc.Delete("/vendors/:id", scHandler.DeleteVendor)

	sc.Get("/item-categories", scHandler.GetItemCategories)
	sc.Get("/item-categories/:id", scHandler.GetItemCategoryByID)
	sc.Post("/item-categories", scHandler.CreateItemCategory)
	sc.Put("/item-categories/:id", scHandler.UpdateItemCategory)
	sc.Delete("/item-categories/:id", scHandler.DeleteItemCategory)
	// Categories Aliases
	sc.Get("/categories", scHandler.GetItemCategories)
	sc.Get("/categories/:id", scHandler.GetItemCategoryByID)
	sc.Post("/categories", scHandler.CreateItemCategory)
	sc.Put("/categories/:id", scHandler.UpdateItemCategory)
	sc.Delete("/categories/:id", scHandler.DeleteItemCategory)

	sc.Get("/items", scHandler.GetItems)
	sc.Get("/items/:id", scHandler.GetItemByID)
	sc.Post("/items", scHandler.CreateItem)
	sc.Put("/items/:id", scHandler.UpdateItem)
	sc.Delete("/items/:id", scHandler.DeleteItem)

	sc.Get("/locations", scHandler.GetLocations)
	sc.Get("/locations/:id", scHandler.GetLocationByID)
	sc.Post("/locations", scHandler.CreateLocation)
	sc.Put("/locations/:id", scHandler.UpdateLocation)
	sc.Delete("/locations/:id", scHandler.DeleteLocation)

	// ── B. Data Transaksi ───────────────────────────────────────────────────────
	// Purchase Requests (PR)
	sc.Get("/purchase-requests", scHandler.GetPurchaseRequests)
	sc.Get("/purchase-requests/:id", scHandler.GetPurchaseRequestByID)
	sc.Post("/purchase-requests", scHandler.CreatePurchaseRequest)
	sc.Put("/purchase-requests/:id", scHandler.UpdatePurchaseRequest)
	sc.Patch("/purchase-requests/:id/status", scHandler.UpdatePurchaseRequestStatus)
	sc.Delete("/purchase-requests/:id", scHandler.DeletePurchaseRequest)

	// Purchase Orders (PO)
	sc.Get("/purchase-orders", scHandler.GetPurchaseOrders)
	sc.Get("/purchase-orders/:id", scHandler.GetPurchaseOrderByID)
	sc.Post("/purchase-orders", scHandler.CreatePurchaseOrder)
	sc.Put("/purchase-orders/:id", scHandler.UpdatePurchaseOrder)
	sc.Patch("/purchase-orders/:id/status", scHandler.UpdatePurchaseOrderStatus)
	sc.Delete("/purchase-orders/:id", scHandler.DeletePurchaseOrder)

	// Goods Receipts (GR / Penerimaan Barang)
	sc.Get("/goods-receipts", scHandler.GetGoodsReceipts)
	sc.Get("/goods-receipts/:id", scHandler.GetGoodsReceiptByID)
	sc.Post("/goods-receipts", scHandler.CreateGoodsReceipt)
	sc.Put("/goods-receipts/:id", scHandler.UpdateGoodsReceipt)
	sc.Delete("/goods-receipts/:id", scHandler.DeleteGoodsReceipt)

	// Manajemen Aset
	sc.Get("/assets", scHandler.GetAssets)
	sc.Get("/assets/:id", scHandler.GetAssetByID)
	sc.Post("/assets", scHandler.CreateAsset)
	sc.Put("/assets/:id", scHandler.UpdateAsset)
	sc.Patch("/assets/:id/status", scHandler.UpdateAssetStatus)
	sc.Delete("/assets/:id", scHandler.DeleteAsset)

	// Langganan / Lisensi Software
	sc.Get("/subscriptions", scHandler.GetSubscriptions)
	sc.Get("/subscriptions/:id", scHandler.GetSubscriptionByID)
	sc.Post("/subscriptions", scHandler.CreateSubscription)
	sc.Put("/subscriptions/:id", scHandler.UpdateSubscription)
	sc.Patch("/subscriptions/:id/status", scHandler.UpdateSubscriptionStatus)
	sc.Delete("/subscriptions/:id", scHandler.DeleteSubscription)

	// ── C. Upload File Bukti / Lampiran SC ──────────────────────────────────────
	sc.Post("/upload", scHandler.UploadFile)

	// ── D. Laporan Supply Chain ─────────────────────────────────────────────────
	sc.Get("/laporan/purchase-orders", scHandler.GetLaporanPO)
	sc.Get("/laporan/aset", scHandler.GetLaporanAset)
	sc.Get("/laporan/lisensi", scHandler.GetLaporanLisensi)
	sc.Get("/laporan/riwayat-vendor", scHandler.GetRiwayatVendor)
	sc.Get("/laporan/riwayat-vendor/:vendor_id", scHandler.GetRiwayatVendor)
	// Reports aliases
	sc.Get("/reports/purchase-orders", scHandler.GetLaporanPO)
	sc.Get("/reports/asset-depreciation", scHandler.GetLaporanAset)
	sc.Get("/reports/subscriptions", scHandler.GetLaporanLisensi)
	sc.Get("/reports/vendor-history", scHandler.GetRiwayatVendor)
	sc.Get("/reports/vendor-history/:vendor_id", scHandler.GetRiwayatVendor)
}
