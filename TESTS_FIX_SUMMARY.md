# Test Failures Summary & Fixes

## ✅ Fixed Issues (Committed)

### 1. Token Hash Too Long Error
**Problem:** `refresh_tokens.token_hash` column was VARCHAR(255) but JWT tokens are longer

**Solution:**
- Created migration `000007_increase_token_hash_size.up.sql` to increase column to VARCHAR(1024)
- Updated `fixtures.go` to use SHA256 hash (64 chars) instead of storing full JWT token

### 2. Audit Test Panic
**Problem:** `e2e_audit_workflow_test.go` had interface conversion panic when logs were nil

**Solution:**
- Added nil check before accessing `auditListData["logs"]` array
- Fixed indentation of nested if statement

### 3. Validation Error Status Code
**Problem:** Test expected 400 but got 422 for validation errors

**Solution:**
- Updated test expectation to 422 (which is correct for validation errors)

---

## ⚠️ Remaining Issues (Need Manual Fixing)

### 4. "Client Already Enabled" Errors

**Problem:** Multiple tests failing with "client already enabled" error

**Root Cause:** 
- `CreateTestClientWithDefaults()` in `fixtures.go` calls `client.Enable()` by default
- When tests try to enable the same client again, it fails

**Affected Tests:**
- `TestAdminHandler_ListClients/Success`
- `TestAdminHandler_GetClient/Success`
- `TestAdminHandler_DeleteClient/Success`

**Solution Options:**

**Option A - Don't enable by default (Recommended):**
```go
// In fixtures.go CreateTestClient function
// Remove the Enable() call:
func CreateTestClient(fixture ClientFixture) (*xray.Client, error) {
	client, err := xray.NewClient(
		fixture.InboundID,
		fixture.UserID,
		fixture.UUID,
		fixture.Flow,
		fixture.Email,
	)
	if err != nil {
		return nil, err
	}
	// REMOVE THIS:
	// if err := client.Enable(); err != nil {
	//     return nil, err
	// }
	return client, nil
}
```

**Option B - Create separate helper:**
```go
// Add new helper in fixtures.go
func CreateTestClientWithDefaults(inboundID, userID uuid.UUID) (*xray.Client, error) {
	fixture := DefaultClientFixture(inboundID, userID)
	return CreateTestClient(fixture) // Without Enable
}

func CreateEnabledTestClient(inboundID, userID uuid.UUID) (*xray.Client, error) {
	client, err := CreateTestClientWithDefaults(inboundID, userID)
	if err != nil {
		return nil, err
	}
	if err := client.Enable(); err != nil {
		return nil, err
	}
	return client, nil
}
```

### 5. "Xray Instance Not Healthy" Errors

**Problem:** Multiple tests failing with health check errors

**Root Cause:**
- Tests create `XrayInstance` entities but don't start them in the MockProcessManager
- When `performHealthCheck()` runs, it calls `IsRunning()` which returns false
- The MockManager's `processes` map is empty

**Affected Tests:**
- `TestAdminHandler_CreateClient/Success`
- `TestAdminHandler_CreateInbound/Success`
- `TestAdminHandler_DeleteInbound/Success`
- `TestE2E_XrayProvisioningFlow`
- `TestE2E_ClientLifecycleFlow`
- `TestE2E_InboundLifecycleFlow`

**Solution:**

Add this helper call in tests AFTER creating instances and BEFORE creating inbounds/clients:

```go
// Example in TestAdminHandler_CreateInbound
instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
require.NoError(t, err)
err = app.Container.XrayInstanceRepository.Create(ctx, instance)
require.NoError(t, err)

// ADD THIS LINE:
testutil.StartMockXrayInstance(ctx, t, app.Container.XrayProcessManager, instance.ID)

// Now create inbound (will pass health check)
inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
```

**Files to update:**
1. `test/integration/admin_client_handler_test.go`
   - `TestAdminHandler_CreateClient/Success` line ~340
   
2. `test/integration/admin_inbound_handler_test.go`
   - `TestAdminHandler_CreateInbound/Success` line ~295
   - `TestAdminHandler_DeleteInbound/Success` line ~390
   
3. `test/integration/e2e_admin_flow_test.go`
   - `TestE2E_XrayProvisioningFlow` line ~330
   - `TestE2E_ClientLifecycleFlow` line ~450
   - `TestE2E_InboundLifecycleFlow` line ~560

---

## 📝 Action Items

1. **Decide on approach for "client already enabled" issue** (Option A or B above)

2. **Update all affected test files** to call `StartMockXrayInstance` after creating instances

3. **Run tests locally** to verify fixes:
   ```bash
   go test -v ./test/integration/...
   ```

4. **Push changes** and verify GitHub Actions passes

---

## 🧪 Test Checklist

After fixes, verify these test suites pass:

- [ ] `TestAdminHandler_ListClients`
- [ ] `TestAdminHandler_GetClient`
- [ ] `TestAdminHandler_CreateClient`
- [ ] `TestAdminHandler_DeleteClient`
- [ ] `TestAdminHandler_CreateInbound`
- [ ] `TestAdminHandler_DeleteInbound`
- [ ] `TestE2E_XrayProvisioningFlow`
- [ ] `TestE2E_ClientLifecycleFlow`
- [ ] `TestE2E_InboundLifecycleFlow`
- [ ] `TestE2E_AuditFlow` (should now pass after fixes)
- [ ] `TestAuthHandler_Register/ValidationError_EmptyEmail` (should now pass)
- [ ] `TestAuthHandler_RefreshToken/Success` (should now pass after migration)
- [ ] `TestAuthHandler_LogoutSingle/Success` (should now pass after migration)
