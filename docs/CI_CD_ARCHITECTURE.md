# 🏗️ GitHub Actions Workflow Architecture

## Overview
This document describes the modular CI/CD pipeline architecture for SuProxy Backend.

## Workflow Separation Strategy

### 🔬 test.yml - Quality Assurance
**Responsibility:** Code quality and correctness validation
- Unit Tests
- Integration Tests (with PostgreSQL)
- Linting (golangci-lint)
- Security Scanning (gosec)
- Build Verification
- Code Coverage

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main` or `develop`

**Status:** ✅ Production Ready

---

### 🐳 build.yml - Container Image Building
**Responsibility:** Docker image creation and registry push
- Checkout code
- Generate version tags (semantic + SHA)
- Docker login (GHCR)
- Build Docker image
- Push multiple tags:
  - `latest`
  - `v1.0.X` (semantic version)
  - `sha-abc123` (commit SHA)

**Triggers:**
- After `test.yml` completes successfully
- Manual dispatch (with optional version override)

**Dependencies:**
- ✅ Requires `test.yml` to pass
- Only runs on `main` branch

**Status:** ✅ Implemented

---

### 🚀 deploy.yml - Deployment (Future)
**Responsibility:** Deploy to production VPS
- SSH to VPS
- Execute `deploy.sh` script
- Health check verification
- Multi-server deployment support

**Triggers:**
- After `build.yml` completes successfully
- Manual dispatch

**Dependencies:**
- ✅ Requires `build.yml` to pass

**Status:** 🔜 Coming in STEP 3

---

### 📦 release.yml - GitHub Releases (Future)
**Responsibility:** Create GitHub releases
- Create Git tag
- Generate release notes
- Generate changelog
- Upload artifacts

**Triggers:**
- Manual dispatch
- Version tags

**Status:** 🔜 Coming in STEP 4

---

## Architecture Diagram

```
┌─────────────────────────────────────┐
│  Push to main / Pull Request        │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│  test.yml (Quality Gate)            │
│  • Unit Tests                       │
│  • Integration Tests                │
│  • Linting                          │
│  • Security Scan                    │
│  • Build Verification               │
└────────────┬────────────────────────┘
             │
             ▼ (only if success)
┌─────────────────────────────────────┐
│  build.yml (Image Builder)          │
│  • Docker Login                     │
│  • Build Image                      │
│  • Push latest                      │
│  • Push v1.0.X                      │
│  • Push sha-abc123                  │
└────────────┬────────────────────────┘
             │
             ▼ (future)
┌─────────────────────────────────────┐
│  deploy.yml (Deployment)            │
│  • SSH to VPS                       │
│  • docker pull                      │
│  • docker-compose up                │
│  • Health Check                     │
└─────────────────────────────────────┘
```

## Image Tagging Strategy

Each successful build produces **3 tags**:

1. **`latest`** - Always points to the most recent build
   ```bash
   ghcr.io/tuncay005-png/suproxy-backend:latest
   ```

2. **Semantic Version** - `v1.0.X` where X is the run number
   ```bash
   ghcr.io/tuncay005-png/suproxy-backend:v1.0.42
   ```

3. **SHA Tag** - First 7 characters of git commit
   ```bash
   ghcr.io/tuncay005-png/suproxy-backend:sha-a1b2c3d
   ```

## Benefits of This Architecture

### ✅ Modularity
- Each workflow has a single responsibility
- Easy to maintain and debug
- Clear separation of concerns

### ✅ Reusability
- `test.yml` runs on every PR and push
- `build.yml` only builds after tests pass
- `deploy.yml` can deploy any version

### ✅ Flexibility
- Can manually trigger any workflow
- Can override version numbers
- Can deploy specific versions

### ✅ Safety
- Tests must pass before building
- Build must succeed before deploying
- Each stage is independently verifiable

### ✅ Traceability
- Every image tagged with semantic version
- Every image tagged with commit SHA
- Easy to identify which code is running

## Workflow Dependencies

```
test.yml
  ↓ (workflow_run)
build.yml
  ↓ (workflow_run - future)
deploy.yml
```

## Manual Workflow Dispatch

All workflows support manual triggering:

```yaml
workflow_dispatch:
  inputs:
    version:
      description: 'Version override'
      required: false
```

This allows:
- Manual builds without waiting for tests
- Version number customization
- Emergency deployments
- Testing pipeline stages

## Future Enhancements

- [ ] `release.yml` - Automated GitHub releases
- [ ] `rollback.yml` - Quick rollback to previous version
- [ ] `backup.yml` - Pre-deployment database backups
- [ ] Multi-server deployment matrix
- [ ] Blue/Green deployment strategy
- [ ] Canary deployment support
- [ ] Automated monitoring integration
