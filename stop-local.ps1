Write-Host "===================================================" -ForegroundColor Cyan
Write-Host "  Menghentikan Kafeinarts ERP Lokal" -ForegroundColor Cyan
Write-Host "===================================================" -ForegroundColor Cyan
Write-Host ""

docker compose down

Write-Host ""
Write-Host "[BERHASIL] Seluruh container berhasil dihentikan." -ForegroundColor Green
