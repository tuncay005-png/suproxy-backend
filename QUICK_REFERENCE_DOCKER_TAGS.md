# Quick Reference: Docker Image Tags

## Tag Types

### SHA Tags (Immutable, Created by CI)
```
ghcr.io/owner/suproxy-backend:sha-abc1234
```
- Created automatically by CI workflow
- One per commit to main branch
- Never changes or gets deleted
- Source of truth for all other tags

### Version Tags (Immutable, Created by Release)
```
ghcr.io/owner/suproxy-backend:v1.0.2
```
- Created manually via Release workflow
- Points to a specific SHA tag
- Used for rollbacks
- Never changes once created

### Latest Tag (Mutable, Updated by CI)
```
ghcr.io/owner/suproxy-backend:latest
```
- Updated automatically by CI workflow
- Always points to newest main branch build
- Used for development/staging
- NOT recommended for production rollbacks

---

## Common Commands

### Find Image for a Commit
```bash
# Get commit SHA
git rev-parse HEAD
# Output: abc1234567890abcdef1234567890abcdef12345

# Calculate SHA tag
SHA_SHORT=$(git rev-parse HEAD | cut -c1-7)
echo "sha-$SHA_SHORT"
# Output: sha-abc1234
```

### Check if Image Exists
```bash
docker manifest inspect ghcr.io/owner/suproxy-backend:v1.0.2
```

### List All Tags for Repository
```bash
# Via GitHub API
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/orgs/owner/packages/container/suproxy-backend/versions

# Or use GitHub Container Registry UI
```

### Create Version Tag Manually
```bash
# Login to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Re-tag existing image
docker buildx imagetools create \
  --tag ghcr.io/owner/suproxy-backend:v1.0.2 \
  ghcr.io/owner/suproxy-backend:sha-abc1234
```

### Pull Specific Version
```bash
# By version tag
docker pull ghcr.io/owner/suproxy-backend:v1.0.2

# By SHA tag
docker pull ghcr.io/owner/suproxy-backend:sha-abc1234

# Latest
docker pull ghcr.io/owner/suproxy-backend:latest
```

---

## Workflow Quick Reference

### Create a Release

1. **Ensure CI passed** for the commit you want to release
   - Check GitHub Actions
   - Verify `sha-xxx` tag exists

2. **Trigger Release Workflow**
   - Go to Actions → Create Release
   - Input version: `v1.0.2`
   - Click "Run workflow"

3. **Workflow will:**
   - Find image `sha-abc1234` for current commit
   - Create alias `v1.0.2` → `sha-abc1234`
   - Create Git tag `v1.0.2`
   - Create GitHub Release
   - Trigger deployment

4. **Result:**
   - Version tag `v1.0.2` available in GHCR
   - Can now rollback to this version

### Rollback to Previous Version

1. **Find version to rollback to**
   - Check GitHub Releases page
   - Or list Git tags: `git tag -l`

2. **Trigger Rollback Workflow**
   - Go to Actions → Rollback Deployment
   - Input version: `v1.0.1`
   - Input servers: `all` or specific servers
   - Input reason: "API latency spike"
   - Click "Run workflow"

3. **Workflow will:**
   - Verify `v1.0.1` image exists
   - Deploy to specified servers
   - Verify health checks
   - Create tracking issue

4. **Result:**
   - Servers running version `v1.0.1`
   - Issue closed automatically if successful

---

## Tag Naming Conventions

### Version Tags
- **Format:** `vMAJOR.MINOR.PATCH`
- **Examples:**
  - `v1.0.0` - Initial release
  - `v1.0.1` - Patch release (bug fixes)
  - `v1.1.0` - Minor release (new features)
  - `v2.0.0` - Major release (breaking changes)
  - `v1.0.0-rc.1` - Release candidate
  - `v1.0.0-beta.1` - Beta version

### SHA Tags
- **Format:** `sha-SHORTCOMMIT`
- **Example:** `sha-abc1234` (first 7 chars of commit SHA)
- **Auto-generated** by CI workflow

---

## Troubleshooting

### "Image not found" during Release

**Problem:** Release workflow can't find `sha-abc1234`

**Causes:**
1. CI workflow hasn't completed yet
2. CI workflow failed
3. Releasing from wrong branch (not main)
4. Commit was never pushed to main

**Solutions:**
1. Check CI workflow status
2. Wait for CI to finish
3. Fix CI test failures
4. Ensure on main branch

### "Version tag already exists"

**Problem:** Trying to create `v1.0.2` but it already exists

**Solutions:**
1. Choose new version number: `v1.0.3`
2. Or delete existing tag (not recommended):
   ```bash
   git tag -d v1.0.2
   git push origin :refs/tags/v1.0.2
   ```

### "Rollback failed" - Image not found

**Problem:** Trying to rollback to `v1.0.2` but tag doesn't exist

**Causes:**
1. Version was never released (no version tag created)
2. Tag was deleted
3. Release workflow failed partway through

**Solutions:**
1. Check available versions:
   ```bash
   git tag -l
   ```
2. Use a different version that exists
3. Or create tag manually for old commits:
   ```bash
   docker buildx imagetools create \
     --tag ghcr.io/owner/suproxy-backend:v1.0.2 \
     ghcr.io/owner/suproxy-backend:sha-abc1234
   ```

---

## Best Practices

### ✅ DO

- Use version tags for production deployments
- Create releases for every production deployment
- Use semantic versioning
- Keep SHA tags forever (they're small)
- Document rollback reasons

### ❌ DON'T

- Don't use `latest` for production
- Don't delete SHA tags
- Don't re-tag version numbers to different images
- Don't skip releases for direct deployments
- Don't rollback without documenting why

---

## Tag Lifecycle

```
Commit pushed to main
        ↓
CI builds and tests
        ↓
    SHA tag created (sha-abc1234)
    Latest tag updated
        ↓
    [Time passes...]
        ↓
Manual release trigger
        ↓
Version tag created (v1.0.2 → sha-abc1234)
        ↓
Deployed to production
        ↓
    [Used in production]
        ↓
    [If problem occurs]
        ↓
Rollback to previous version (v1.0.1)
        ↓
Production restored
```

---

## Quick Checks

### Check Image Tags for Current Commit
```bash
COMMIT=$(git rev-parse HEAD)
SHA_SHORT=$(echo $COMMIT | cut -c1-7)
echo "SHA tag: sha-$SHA_SHORT"
docker manifest inspect ghcr.io/owner/suproxy-backend:sha-$SHA_SHORT
```

### Check All Tags for an Image Digest
```bash
DIGEST="sha256:abc123..."
# Find all tags pointing to this digest
docker buildx imagetools inspect ghcr.io/owner/suproxy-backend@$DIGEST
```

### Verify Version Tag Exists Before Rollback
```bash
VERSION="v1.0.2"
if docker manifest inspect ghcr.io/owner/suproxy-backend:$VERSION >/dev/null 2>&1; then
  echo "✅ Version $VERSION exists"
else
  echo "❌ Version $VERSION not found"
fi
```

---

## Summary

**Three tag types:**
1. **SHA tags** - One per commit, never changes, source of truth
2. **Version tags** - One per release, never changes, for rollbacks
3. **Latest tag** - Always newest, changes often, not for production

**Key principle:** Build once (CI), tag many times (Release)

**Rollback-ready:** Every release creates a permanent version tag for instant rollbacks.
