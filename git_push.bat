@echo off
echo ========================================
echo Git Push Script
echo ========================================
echo.

echo Checking git status...
git status
echo.

echo ========================================
echo Adding all changes...
git add .
echo.

echo ========================================
echo Creating commit...
set /p commit_msg="Commit message (or press Enter for default): "
if "%commit_msg%"=="" set commit_msg=Fix test compilation errors and FK violations

git commit -m "%commit_msg%"
echo.

echo ========================================
echo Pushing to origin main...
git push origin main
echo.

echo ========================================
echo Done!
echo ========================================
pause
