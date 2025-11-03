# Version Management Instructions for Claude Code

## IMPORTANT: Read This Before Making ANY Code Changes

**Repository:** https://github.com/johnzastrow/waterlogger

### Current Version: 1.5.1

---

## Version Bumping Policy

### ALWAYS increment the version when making code changes:

- **MAJOR (X.0.0)** - Breaking changes or incompatible API changes
  - Example: 1.4.0 → 2.0.0

- **MINOR (X.Y.0)** - New features (backward compatible)
  - Example: 1.4.0 → 1.5.0
  - Use this for: New features, significant enhancements, new middleware, new packages

- **PATCH (X.Y.Z)** - Bug fixes (backward compatible)
  - Example: 1.4.0 → 1.4.1
  - Use this for: Bug fixes, typo corrections, small improvements, documentation updates

### When NOT to bump version:
- Only updating comments
- Fixing typos in documentation (markdown files only)
- Formatting code (no logic changes)
- Updating README or non-code documentation

---

## Checklist: How to Bump Version

When making code changes, follow these steps IN ORDER:

### 1. **Determine New Version Number**
Based on the type of change (MAJOR/MINOR/PATCH), determine the new version.

Example:
- Adding new feature → Increment MINOR: 1.4.0 → 1.5.0
- Fixing bug → Increment PATCH: 1.4.0 → 1.4.1
- Breaking change → Increment MAJOR: 1.3.0 → 2.0.0

---

### 2. **Update These Files (ALL REQUIRED):**

#### A. `VERSION` file
```
1.4.0
```

#### B. `cmd/waterlogger/main.go` (2 locations)
```go
// Line ~52
fmt.Println("Waterlogger v1.4.0")

// Line ~112
logging.Info().Str("version", "1.4.0")...
```

#### C. `internal/config/config.go`
```go
// Line ~120 in Default() function
Version:   "1.4.0",
```

#### D. `config.yaml`
```yaml
app:
    name: Waterlogger
    version: 1.4.0
```

#### E. `config.example.yaml`
```yaml
app:
    name: "Waterlogger"
    version: "1.4.0"
```

#### F. `Makefile`
```makefile
VERSION := 1.4.0
```

---

### 3. **Update CHANGELOG.md**

Add a new section at the top (after `## [Unreleased]`):

```markdown
## [1.4.0] - YYYY-MM-DD

### Added
- Feature description here
- Another feature

### Changed
- What changed

### Fixed
- What was fixed

### Technical Details
- Technical notes
```

**Format for version anchor:**
- Version `1.4.0` becomes `#140---YYYY-MM-DD`
- Version `2.0.0` becomes `#200---YYYY-MM-DD`
- Version `1.3.1` becomes `#131---YYYY-MM-DD`

---

### 4. **Update ALL Template Files**

Update the CHANGELOG link in all 10 template files:

**Files to update:**
1. `web/templates/dashboard.html`
2. `web/templates/pools.html`
3. `web/templates/kits.html`
4. `web/templates/samples.html`
5. `web/templates/export.html`
6. `web/templates/settings.html`
7. `web/templates/adjustments.html`
8. `web/templates/setup.html`
9. `web/templates/login.html`
10. `web/templates/404.html`

**Find this in each file:**
```html
<a href="https://github.com/johnzastrow/waterlogger/blob/main/CHANGELOG.md#130---2025-10-26" target="_blank" title="View changelog">v{{.Version}}</a>
```

**Change to:**
```html
<a href="https://github.com/johnzastrow/waterlogger/blob/main/CHANGELOG.md#140---YYYY-MM-DD" target="_blank" title="View changelog">v{{.Version}}</a>
```

Replace:
- `#140` with the correct anchor (version without dots)
- `YYYY-MM-DD` with today's date

---

### 5. **Update TODO.md**

Update the header:
```markdown
**Last Updated:** YYYY-MM-DD
**Version:** 1.4.0
```

---

## Quick Reference: Files to Update

### Must Update (7 files):
1. ✅ `VERSION`
2. ✅ `cmd/waterlogger/main.go` (2 places)
3. ✅ `internal/config/config.go`
4. ✅ `config.yaml`
5. ✅ `config.example.yaml`
6. ✅ `Makefile`
7. ✅ `CHANGELOG.md`

### Must Update (10 template files):
8. ✅ `web/templates/dashboard.html`
9. ✅ `web/templates/pools.html`
10. ✅ `web/templates/kits.html`
11. ✅ `web/templates/samples.html`
12. ✅ `web/templates/export.html`
13. ✅ `web/templates/settings.html`
14. ✅ `web/templates/adjustments.html`
15. ✅ `web/templates/setup.html`
16. ✅ `web/templates/login.html`
17. ✅ `web/templates/404.html`

### Should Update:
18. ✅ `TODO.md` (header with date and version)

---

## Example: Bumping from 1.3.0 to 1.4.0

Let's say you're adding a new feature (Minor version bump):

**Today's Date:** 2025-10-27

### Step-by-Step:

1. **VERSION file:**
   ```
   1.4.0
   ```

2. **main.go** (line 52):
   ```go
   fmt.Println("Waterlogger v1.4.0")
   ```

3. **main.go** (line 112):
   ```go
   logging.Info().Str("version", "1.4.0")...
   ```

4. **config.go** (line 120):
   ```go
   Version:   "1.4.0",
   ```

5. **config.yaml**:
   ```yaml
   version: 1.4.0
   ```

6. **config.example.yaml**:
   ```yaml
   version: "1.4.0"
   ```

7. **Makefile**:
   ```makefile
   VERSION := 1.4.0
   ```

8. **CHANGELOG.md** (add after Unreleased):
   ```markdown
   ## [1.4.0] - 2025-10-27

   ### Added
   - New awesome feature description
   ```

9. **All 10 template files** - Update link:
   ```html
   <a href="https://github.com/johnzastrow/waterlogger/blob/main/CHANGELOG.md#140---2025-10-27" target="_blank" title="View changelog">v{{.Version}}</a>
   ```

10. **TODO.md**:
    ```markdown
    **Last Updated:** 2025-10-27
    **Version:** 1.4.0
    ```

---

## Verification Checklist

After updating version, verify:

- [ ] All 7 core files have new version
- [ ] All 10 template files have updated CHANGELOG link
- [ ] CHANGELOG.md has new entry with correct date
- [ ] CHANGELOG anchor format is correct (#140 not #1.4.0)
- [ ] TODO.md header updated
- [ ] Version string is consistent everywhere (no quotes in some places, quotes in others)

---

## Common Mistakes to Avoid

❌ **DON'T:**
- Forget to update template files (they have hardcoded CHANGELOG links)
- Use wrong anchor format in CHANGELOG link (#1.4.0 instead of #140)
- Forget to update the date in CHANGELOG
- Update version in only some files
- Forget to increment version at all when making code changes

✅ **DO:**
- Update ALL 17 files listed above
- Use correct anchor format (no dots, no spaces)
- Use today's date in CHANGELOG
- Double-check version consistency
- Always increment version for code changes

---

## GitHub Repository Info

- **Owner:** johnzastrow
- **Repository:** waterlogger
- **CHANGELOG URL:** https://github.com/johnzastrow/waterlogger/blob/main/CHANGELOG.md

---

## Current Version History

- **1.5.1** (2025-11-01) - Production logging path fix
- **1.5.0** (2025-10-27) - Database backup import & automated deployment
- **1.4.0** (2025-10-27) - Database schema migration system & enhanced settings UI
- **1.3.0** (2025-10-26) - Comprehensive logging system with zerolog
- **1.2.0** (2025-10-25) - Pool volume calculator, chemical adjustments
- **1.0.0** (2024-07-14) - Initial release

---

## Questions?

If unsure about version bump:
- Adding features? → MINOR (1.X.0)
- Fixing bugs? → PATCH (1.3.X)
- Breaking changes? → MAJOR (X.0.0)

When in doubt, use MINOR for any significant change.

---

**Last Updated:** 2025-11-01
**Maintained By:** Claude Code AI Assistant
