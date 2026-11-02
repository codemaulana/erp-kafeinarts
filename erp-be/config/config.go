package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config merepresentasikan variabel konfigurasi global aplikasi.
type Config struct {
	Port               string
	AppEnv             string
	AppName            string
	DBPath             string
	JWTSecret          string
	JWTExpirationHours int
	CORSOrigins        []string
}

// Global instance config
var AppConfig *Config

// LoadConfig membaca environment variables dari .env (jika ada) dan OS environment.
func LoadConfig() *Config {
	// Muat file .env jika ada di folder kerja
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan nilai bawaan atau environment OS")
	}

	port := getEnv("PORT", "8000")
	appEnv := getEnv("APP_ENV", "development")
	appName := getEnv("APP_NAME", "Kafeinarts ERP API")
	dbPath := getEnv("DB_PATH", "./storage/erp.db")
	jwtSecret := getEnv("JWT_SECRET", "kafeinarts-super-secret-key-2026-very-secure")

	jwtExpHoursStr := getEnv("JWT_EXPIRATION_HOURS", "72")
	jwtExpHours, err := strconv.Atoi(jwtExpHoursStr)
	if err != nil {
		jwtExpHours = 72
	}

	corsOriginsRaw := getEnv("CORS_ORIGINS", "http://localhost:9000,http://localhost:9001,http://localhost:8080,http://127.0.0.1:9000")
	corsOrigins := strings.Split(corsOriginsRaw, ",")
	for i, origin := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(origin)
	}

	AppConfig = &Config{
		Port:               port,
		AppEnv:             appEnv,
		AppName:            appName,
		DBPath:             dbPath,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: jwtExpHours,
		CORSOrigins:        corsOrigins,
	}

	return AppConfig
}

// getEnv membaca env var berdasarkan key, atau mengembalikan fallback jika kosong.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}
