package test

import (
	"testing"

	"erp-kafeinart/utils"
)

func TestPasswordHashing(t *testing.T) {
	password := "Kafeinarts2026!"

	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Gagal melakukan hashing password: %v", err)
	}

	if !utils.CheckPasswordHash(password, hash) {
		t.Errorf("Password yang benar harusnya valid terhadap hash")
	}

	if utils.CheckPasswordHash("wrongpassword", hash) {
		t.Errorf("Password yang salah tidak boleh valid terhadap hash")
	}
}

func TestJWTGenerationAndValidation(t *testing.T) {
	secret := "test-secret-key-123"
	userID := uint(42)
	email := "hanif@kafeinarts.com"
	name := "Hanif"
	role := "superadmin"
	expHours := 1

	token, err := utils.GenerateJWT(userID, email, name, role, secret, expHours)
	if err != nil {
		t.Fatalf("Gagal generate JWT: %v", err)
	}

	if token == "" {
		t.Fatalf("Token JWT tidak boleh kosong")
	}

	claims, err := utils.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Gagal memvalidasi JWT: %v", err)
	}

	if claims.UserID != userID || claims.Email != email || claims.Role != role {
		t.Errorf("Klaim JWT tidak cocok dengan input")
	}

	// Test secret salah
	_, err = utils.ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Errorf("Token harus gagal divalidasi jika secret berbeda")
	}
}

func TestPaginationMeta(t *testing.T) {
	meta := utils.BuildPaginationMeta(45, 2, 10)

	if meta.Total != 45 {
		t.Errorf("Expected total 45, got %d", meta.Total)
	}
	if meta.Page != 2 || meta.CurrentPage != 2 {
		t.Errorf("Expected page 2, got %d / %d", meta.Page, meta.CurrentPage)
	}
	if meta.PerPage != 10 {
		t.Errorf("Expected per_page 10, got %d", meta.PerPage)
	}
	if meta.TotalPages != 5 {
		t.Errorf("Expected total_pages 5, got %d", meta.TotalPages)
	}
}
