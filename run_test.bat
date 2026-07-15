@echo off
cd /d %~dp0
go test -v ./test/integration/... -run TestAdminClientHandler
pause
