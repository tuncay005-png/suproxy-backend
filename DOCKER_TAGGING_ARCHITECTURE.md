# Docker Image Tagging Architecture

## Overview

This document explains our production-grade Docker image tagging strategy that enables:
- **Build once, deploy many times**
- **Fast releases without rebuilding**
- **Rollback to any version by tag**
- **Guaranteed consistency** between versions and tested images

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ Developer pushes code to main branch                             │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ CI Workflow (Automated)                                          │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│ 1. Run tests (unit, integration, lint)                          │
│ 2. Build Docker image once                                      │
│ 3. Push with multiple tags:                                     │
│    • sha-abc1234 (commit-based, immutable)                      │
│    • latest (mutable, always points to newest main)             │
│ 4. Store image digest for traceability                          │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             │ Image stored in GHCR
                             │ Digest: sha256:def456...
                             │
┌────────────────────────────┴────────────────────────────────────┐
│ Available Docker Tags (all point to SAME image):                │
│ • ghcr.io/owner/suproxy-backend:sha-abc1234                     │
│ • ghcr.io/owner/suproxy-backend:latest                          │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             │ (Time passes... ready to release)
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ Release Workflow (Manual Trigger)                               │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│ Input: version = v1.0.2                                          │
│                                                                  │
│ 1. Find commit SHA for this release                             │
│    → github.sha = abc1234...                                    │
│                                                                  │
│ 2. Calculate source image tag                                   │
│    → sha-abc1234                                                │
│                                                                  │
│ 3. Verify source image exists in GHCR                           │
│    ✓ ghcr.io/owner/suproxy-backend:sha-abc1234 exists          │
│                                                                  │
│ 4. Create version tag alias (NO REBUILD!)                       │
│    docker buildx imagetools create \                            │
│      --tag ghcr.io/owner/suproxy-backend:v1.0.2 \              │
│      ghcr.io/owner/suproxy-backend:sha-abc1234                 │
│                                                                  │
│ 5. Create Git tag v1.0.2                                        │
│ 6. Create GitHub Release                                        │
│ 7. Trigger deployment workflow                                  │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ Available Docker Tags (all point to SAME image):                │
│ • ghcr.io/owner/suproxy-backend:sha-abc1234                     │
│ • ghcr.io/owner/suproxy-backend:latest                          │
│ • ghcr.io/owner/suproxy-backend:v1.0.2  ← NEW!                  │
│                                                                  │
│ All share the same image digest: sha256:def456...               │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ Rollback Workflow                                               │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│ Input: version = v1.0.2                                          │
│                                                                  │
│ 1. Verify image exists                                          │
│    ✓ ghcr.io/owner/suproxy-backend:v1.0.2 exists               │
│                                                                  │
│ 2. Deploy to servers                                            │
│    docker pull ghcr.io/owner/suproxy-backend:v1.0.2            │
│                                                                  │
│ 3. Verify health checks                                         │
│    ✓ All servers healthy                                        │
└─────────────────────────────────────────────────────────────────┘
```

## Key Components

### 1. CI Workflow (ci.yml)

**Trigger:** Automatic on push to `main` branch

**Responsibilities:**
- Run all tests (unit, integration, linting)
- Build Docker image **exactly once**
- Push image with two tags:
  - `sha-<commit>` (e.g., `sha-abc1234`) - Immutable, content-addressable
  - `latest` - Mutable, always points to newest main build

**Why SHA tags?**
- Immutable: Never changes once created
- Traceable: Direct link to source code commit
- Unique: Every commit gets a unique image
- Reliable: Can always find the exact build for any commit

### 2. Release Workflow (release.yml)

**Trigger:** Manual via GitHub Actions UI

**Responsibilities:**
- Find the Docker image for the current commit (sha-xxx tag)
- **Re-tag** existing image with version tag (e.g., `v1.0.2`)
- Create Git tag
- Create GitHub Release
- Trigger deployment

**Key Innovation: Image Re-tagging**

Instead of rebuilding the image:

```bash
# OLD WAY (slow, wasteful, risky):
docker build -t app:v1.0.2 .  # Rebuild from scratch
docker push app:v1.0.2

# NEW WAY (fast, efficient, safe):
docker buildx imagetools create \
  --tag app:v1.0.2 \
  app:sha-abc1234  # Just copy the manifest pointer
```

**Benefits:**
- ⚡ **Fast**: Completes in seconds (no build time)
- 💾 **Efficient**: No duplicate storage (same image digest)
- ✅ **Safe**: Release EXACTLY what was tested in CI
- 🔒 **Immutable**: Version tag points to tested, validated image

### 3. Rollback Workflow (rollback.yml)

**Trigger:** Manual via GitHub Actions UI

**Responsibilities:**
- Verify version tag exists (e.g., `v1.0.2`)
- Deploy that exact image to specified servers
- Verify health checks pass

**Why it works:**
- Every release creates a version tag
- Version tags are permanent and never change
- Can rollback to any previous version instantly

## Image Tag Relationships

```
Commit abc1234 on main branch
    │
    ├─► CI builds image once
    │     Image Digest: sha256:def456789...
    │
    ├─► Tagged in CI workflow:
    │     • sha-abc1234  (immutable source of truth)
    │     • latest       (mutable pointer)
    │
    └─► Tagged in Release workflow:
          • v1.0.2       (immutable version pointer)
          • v1.0.3       (can be added later for hotfixes)

All tags point to SAME image blob (sha256:def456789...)
```

## Workflow Sequence

### Happy Path: Feature → Production

1. **Developer** pushes code to `main` branch
   ```
   git push origin main
   ```

2. **CI Workflow** (automatic)
   - Runs tests
   - Builds image
   - Tags: `sha-abc1234`, `latest`
   - Pushes to GHCR

3. **Release Workflow** (manual trigger)
   ```
   Workflow Dispatch:
     version: v1.0.2
   ```
   - Finds image `sha-abc1234`
   - Re-tags as `v1.0.2`
   - Creates GitHub Release
   - Triggers deployment

4. **Deploy Workflow** (automatic)
   - Pulls `v1.0.2` image
   - Deploys to all servers
   - Verifies health

### Rollback Scenario

1. **Problem detected** in production running `v1.0.2`

2. **Rollback Workflow** (manual trigger)
   ```
   Workflow Dispatch:
     version: v1.0.1
     servers: all
     reason: "API latency spike"
   ```

3. **Rollback executes**
   - Verifies `v1.0.1` image exists ✓
   - Deploys to all servers
   - Verifies health checks ✓

4. **Production restored** to previous stable version

## Comparison: Old vs New Architecture

### Old Architecture (Build in Release)

```
❌ CI Workflow:
   • Tests only
   • No Docker build

❌ Release Workflow:
   • Build image from scratch
   • Push v1.0.2 tag
   • Takes 5-10 minutes
   • Different environment than CI
   • Might have different dependencies
   • Not tested before release

❌ Rollback:
   • Version tags don't exist
   • Can only rollback to commits
   • Must use sha-xxx tags
   • Confusing version management
```

### New Architecture (Build Once, Tag Many)

```
✅ CI Workflow:
   • Tests
   • Build image ONCE
   • Push sha-xxx and latest
   • 5-10 minutes

✅ Release Workflow:
   • Re-tag existing image
   • Push v1.0.2 tag (alias)
   • Takes 10-30 seconds
   • Same image that passed CI tests
   • Guaranteed consistency

✅ Rollback:
   • Version tags exist
   • Easy to understand (v1.0.1, v1.0.2)
   • Fast rollback (image already in registry)
   • Clear version management
```

## Technical Details

### Docker Buildx Imagetools

The `docker buildx imagetools create` command creates a new manifest that points to an existing image:

```bash
# Source image (already exists)
SOURCE="ghcr.io/owner/app:sha-abc1234"

# Create new tag pointing to same image
docker buildx imagetools create \
  --tag ghcr.io/owner/app:v1.0.2 \
  $SOURCE

# Result: Both tags point to SAME image digest
# No data duplication
# No rebuild required
# Instant operation
```

### Image Digests

Every Docker image has a unique content-addressable digest:

```
ghcr.io/owner/app:sha-abc1234@sha256:def456789...
                  ↑                      ↑
                  Tag (mutable)          Digest (immutable)
```

When we re-tag an image, we're creating a new tag that points to the same digest.

### Storage Efficiency

```
# Traditional approach:
sha-abc1234:  1.2 GB
v1.0.2:       1.2 GB (duplicate)
v1.0.3:       1.2 GB (duplicate)
Total:        3.6 GB

# Re-tagging approach:
sha-abc1234:  1.2 GB
v1.0.2:       0 KB (pointer)
v1.0.3:       0 KB (pointer)
Total:        1.2 GB
```

## Best Practices

### 1. Never Delete SHA Tags

SHA tags are the source of truth. Version tags depend on them.

```bash
# ❌ NEVER DO THIS
docker rmi ghcr.io/owner/app:sha-abc1234

# ✅ Safe to delete old version tags if needed
docker rmi ghcr.io/owner/app:v0.9.0
```

### 2. Always Release from Main Branch

Release workflow should only run on `main` branch where CI has built and tested the image.

### 3. Use Semantic Versioning

- `v1.0.0` - Major release
- `v1.0.1` - Patch release
- `v1.1.0` - Minor release
- `v2.0.0` - Breaking changes

### 4. Version Tags Are Immutable

Once a version is released, never change what it points to. Create a new version instead.

```bash
# ❌ WRONG: Re-tagging existing version
docker buildx imagetools create \
  --tag app:v1.0.2 \
  app:sha-xyz9999  # Different source image

# ✅ CORRECT: Create new version
docker buildx imagetools create \
  --tag app:v1.0.3 \
  app:sha-xyz9999
```

## Troubleshooting

### "Source image not found: sha-abc1234"

**Cause:** CI workflow hasn't run yet or failed for this commit.

**Solution:**
1. Check CI workflow status on GitHub Actions
2. Wait for CI to complete if it's running
3. Fix CI failures if tests failed
4. Ensure you're releasing from a commit on `main` branch

### "Version tag already exists"

**Cause:** You're trying to create a release with a version that already exists.

**Solution:**
1. Check existing releases: `git tag -l`
2. Choose a new version number
3. Or delete the old tag if it was a mistake (not recommended)

### "Rollback image not found"

**Cause:** Trying to rollback to a version that was never released.

**Solution:**
1. Check available versions: `docker images` or GitHub Releases
2. Use an existing version tag
3. Or use a SHA tag directly: `sha-abc1234`

## Migration Guide

If you're migrating from the old architecture:

1. **Update CI workflow** to build and push Docker images
2. **Update Release workflow** to re-tag instead of rebuild
3. **Create version tags** for existing commits you want to rollback to:
   ```bash
   # For each important commit:
   docker buildx imagetools create \
     --tag app:v1.0.0 \
     app:sha-abc1234
   ```

## Additional Resources

- [Docker Buildx Documentation](https://docs.docker.com/build/buildx/)
- [Docker Multi-platform Images](https://docs.docker.com/build/building/multi-platform/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Semantic Versioning](https://semver.org/)

## Summary

This architecture provides:

✅ **Single source of truth** - SHA tags from CI  
✅ **Fast releases** - Re-tagging is instant  
✅ **Guaranteed consistency** - Release exactly what was tested  
✅ **Easy rollbacks** - Version tags always available  
✅ **Storage efficient** - No duplicate images  
✅ **Production-grade** - Used by major companies  

**Key principle**: Build once in CI, tag many times as needed for releases and rollbacks.
