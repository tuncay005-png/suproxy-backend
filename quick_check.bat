@echo off
echo ======================================
echo Quick Build Check
echo ======================================
go build ./...
if %errorlevel% neq 0 (
    echo Build FAILED!
    pause
    exit /b %errorlevel%
)

echo.
echo ======================================
echo Build SUCCESS!
echo ======================================
echo.
echo Run full tests with: go test -v ./test/integration/...
pause
