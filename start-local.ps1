Write-Host "===================================================" -ForegroundColor Cyan
Write-Host "  Menjalankan Kafeinarts ERP di Lingkungan Lokal" -ForegroundColor Cyan
Write-Host "===================================================" -ForegroundColor Cyan
Write-Host ""

docker compose up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Gagal menjalankan Docker Compose. Pastikan Docker Desktop aktif." -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "===================================================" -ForegroundColor Green
Write-Host "[BERHASIL] Seluruh container berjalan di background!" -ForegroundColor Green
Write-Host ""
Write-Host "Web Aplikasi (Frontend) : http://localhost:3008" -ForegroundColor Yellow
Write-Host "API Endpoint (Backend)  : http://localhost:8000/api/v1 (atau :8484)" -ForegroundColor Yellow
Write-Host ""
Write-Host "Login Bawaan:"
Write-Host "- Email    : admin@kafeinarts.com"
Write-Host "- Password : password123"
Write-Host "===================================================" -ForegroundColor Green
Write-Host ""
Start-Process "http://localhost:3008"
