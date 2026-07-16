@echo off
set INTEGRATION_TEST=true
go test -v ./test/integration -run TestAdminHandler_ListClients
go test -v ./test/integration -run TestAdminHandler_GetClient
go test -v ./test/integration -run TestAdminHandler_DeleteClient
go test -v ./test/integration -run TestE2E_AuditFilteringFlow
