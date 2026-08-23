# Migration Debug Guide

## Changes Applied

The `internal/infrastructure/database/migrator.go` file has been enhanced with comprehensive debug logging and improved path resolution to identify why migrations 008, 009, and 010 are not being applied.

### 1. Enhanced Path Resolution (getInstance method)

**Three-Strategy Approach:**
1. **Strategy 1**: Try relative to current working directory
2. **Strategy 2**: Try relative to executable location  
3. **Strategy 3**: Try parent directories (up to 3 levels)

**Debug Output Added:**
```go
// Logs current working directory and resolved path
m.logger.Info("Migration path resolution",
    "working_directory", wd,
    "relative_path", "migrations",
    "absolute_path", migrationsPath,
)

// Logs final resolved path
m.logger.Info("Final migrations path resolved", "path", migrationsPath)

// Lists all .up.sql files found
m.logger.Info("Migration files discovered",
    "count", len(migrationFiles),
    "files", migrationFiles,
)
```

### 2. Enhanced Up() Method

**Debug Output Added:**
```go
// Shows version before attempting migration
m.logger.Info("Current migration version before Up()", "version", currentVersion)

// Distinguishes between "no changes" and actual errors
if err == migrate.ErrNoChange {
    m.logger.Info("No new migrations to apply (already up-to-date)")
}

// Shows version after migration
m.logger.Info("Database migrations completed",
    "version", version,
    "dirty", dirty,
)

// Shows version change if it occurred
if currentVersion != version {
    m.logger.Info("Migration version updated",
        "from", currentVersion,
        "to", version,
    )
}
```

### 3. New CheckMigrationState() Method

Added utility method to check current migration state:
```go
func (m *Migrator) CheckMigrationState() error
```

This can be called independently to verify migration status without running migrations.

## Testing Instructions

### 1. Start the Backend

```bash
cd c:\Users\Tuncay\Desktop\suproxy-backend
.\api.exe
```

### 2. Expected Log Output

You should now see detailed migration logs showing:

#### Path Resolution Logs:
```
INFO Migration path resolution working_directory=C:\Users\Tuncay\Desktop\suproxy-backend relative_path=migrations absolute_path=C:\Users\Tuncay\Desktop\suproxy-backend\migrations
INFO Final migrations path resolved path=C:\Users\Tuncay\Desktop\suproxy-backend\migrations
INFO Migration files discovered count=10 files=[000001_init_schema.up.sql 000002_create_users_table.up.sql ... 000010_add_token_family.up.sql]
```

#### Migration Execution Logs:
```
INFO Running database migrations...
INFO Current migration version before Up() version=7
INFO Database migrations completed version=10 dirty=false
INFO Migration version updated from=7 to=10
```

OR if no migrations needed:
```
INFO Running database migrations...
INFO Current migration version before Up() version=10
INFO No new migrations to apply (already up-to-date)
INFO Database migrations completed version=10 dirty=false
```

### 3. Diagnosing Issues

Based on the logs, you'll be able to identify the root cause:

#### Issue A: Migrations Directory Not Found
**Log Pattern:**
```
WARN Migrations directory not found at working directory, searching alternatives attempted_path=...
ERROR migrations directory not found: tried ... and parent directories
```

**Solution:** The executable is running from an unexpected location. Check the working directory and ensure migrations folder is accessible.

#### Issue B: Migration Files Not Found
**Log Pattern:**
```
INFO Migration files discovered count=0 files=[]
```

**Solution:** The migrations directory is empty or doesn't contain .up.sql files. Verify the migrations folder contains all 10 migration files.

#### Issue C: Stuck at Version 7
**Log Pattern:**
```
INFO Current migration version before Up() version=7
INFO No new migrations to apply (already up-to-date)
INFO Database migrations completed version=7 dirty=false
```

**Possible Causes:**
1. Migration files 8, 9, 10 not found in the directory
2. golang-migrate not detecting new migrations (check file naming)
3. Database schema_migrations table corrupted

**Solution:** Check that:
- All 10 migration files are present (see list above)
- File names follow the pattern: `000008_description.up.sql`
- Database `schema_migrations` table is not corrupted

#### Issue D: Dirty Migration State
**Log Pattern:**
```
INFO Database migrations completed version=X dirty=true
```

**Solution:** A previous migration failed partway through. Run:
```bash
# Connect to database and fix manually
# Or use the Down() method to rollback, then try Up() again
```

## Migration Files Present

Confirmed migrations in `c:\Users\Tuncay\Desktop\suproxy-backend\migrations\`:

1. ✅ 000001_init_schema.up.sql
2. ✅ 000002_create_users_table.up.sql
3. ✅ 000003_add_security_fields.up.sql
4. ✅ 000004_create_plans_and_subscriptions.up.sql
5. ✅ 000005_create_servers_and_nodes.up.sql
6. ✅ 000006_create_xray_tables.up.sql
7. ✅ 000007_increase_token_hash_size.up.sql
8. ✅ 000008_add_users_created_at_index.up.sql
9. ✅ 000009_add_users_email_lower_index.up.sql
10. ✅ 000010_add_token_family.up.sql

## Next Steps

1. ✅ Code changes applied to `migrator.go`
2. ✅ Binary rebuilt (`api.exe`)
3. ⏳ **Start backend and check logs**
4. ⏳ **Identify root cause from debug output**
5. ⏳ **Apply specific fix based on diagnosis**

## Files Modified

- `internal/infrastructure/database/migrator.go` - Enhanced with debug logging and improved path resolution
- `api.exe` - Rebuilt with changes

## Rebuild Command

```bash
cd c:\Users\Tuncay\Desktop\suproxy-backend
go build -o api.exe ./cmd/api
```
