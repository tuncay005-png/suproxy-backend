@echo off
echo ======================================
echo Building project...
echo ======================================
go build ./...
if %errorlevel% neq 0 (
    echo Build FAILED!
    pause
    exit /b %errorlevel%
)

echo.
echo ======================================
echo Build SUCCESS! Running integration tests...
echo ======================================
go test -v ./test/integration/... -run TestAdminClientHandler
if %errorlevel% neq 0 (
    echo Tests FAILED!
    pause
    exit /b %errorlevel%
)

echo.
echo ======================================
echo All tests PASSED!
echo ======================================
pause
