@echo off
echo Running govulncheck...
govulncheck ./... > govulncheck_after.txt 2>&1
echo.
echo Govulncheck completed. Output saved to govulncheck_after.txt
type govulncheck_after.txt
