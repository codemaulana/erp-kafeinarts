@echo off
echo ===================================================
echo   Menjalankan Kafeinarts ERP di Lingkungan Lokal
echo ===================================================
echo.
docker compose up -d
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Gagal menjalankan Docker Compose. Pastikan Docker Desktop aktif.
    pause
    exit /b %ERRORLEVEL%
)
echo.
echo ===================================================
echo [BERHASIL] Seluruh container berjalan di background!
echo.
echo Web Aplikasi (Frontend) : http://localhost:3008
echo API Endpoint (Backend)  : http://localhost:8000/api/v1 (atau :8484)
echo.
echo Login Bawaan:
echo - Email    : admin@kafeinarts.com
echo - Password : password123
echo ===================================================
echo.
echo Membuka browser...
start http://localhost:3008
pause
