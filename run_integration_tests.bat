@echo off
echo ======================================
echo Running Integration Tests
echo ======================================
echo.
echo Test will run with -count=1 flag to bypass cache
echo.

go test ./test/integration/... -count=1 -v

if %errorlevel% neq 0 (
    echo.
    echo ======================================
    echo Integration Tests FAILED!
    echo ======================================
    pause
    exit /b %errorlevel%
)

echo.
echo ======================================
echo Integration Tests PASSED!
echo ======================================
pause
