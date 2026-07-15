@echo off
echo ======================================
echo Git Status
echo ======================================
git status

echo.
echo ======================================
echo Adding changes...
echo ======================================
git add .

echo.
echo ======================================
echo Creating commit...
echo ======================================
git commit -m "fix: add nil checks to repositories and improve test error handling"

echo.
echo ======================================
echo Pushing to remote...
echo ======================================
git push

echo.
echo ======================================
echo Done!
echo ======================================
pause
