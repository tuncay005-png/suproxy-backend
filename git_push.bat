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
git commit -m "fix: resolve test port conflicts with atomic counter

- Add atomic counter for unique test port generation
- Fix idx_nodes_server_port duplicate key violation
- Add servers table to cleanup for proper test isolation
- Affects only test infrastructure, no production code changes"

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
