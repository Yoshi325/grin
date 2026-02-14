# Homebrew Tap Distribution Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Publish grin on Homebrew via a personal tap so users can `brew install Yoshi325/tap/grin`.

**Architecture:** Add goreleaser's `brews` config to auto-push a formula to `Yoshi325/homebrew-tap` on each tagged release. Pass a dedicated PAT via GitHub Actions secrets.

**Tech Stack:** goreleaser v2, GitHub Actions, Homebrew

---

### Task 1: Add `brews` section to `.goreleaser.yml`

**Files:**
- Modify: `.goreleaser.yml:29` (after the `changelog` block, before `source`)

**Step 1: Add the brews configuration**

Open `.goreleaser.yml` and add the following block between the `changelog` section and the `source` section:

```yaml
brews:
  - repository:
      owner: Yoshi325
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
    homepage: "https://github.com/Yoshi325/grin"
    description: "Make INI files greppable"
    license: "MIT"
    install: |
      bin.install "grin"
    test: |
      assert_match "grin version", shell_output("#{bin}/grin --version")
```

**Step 2: Validate goreleaser config**

Run: `goreleaser check`
Expected: no errors. If goreleaser is not installed locally, skip — CI will validate on release.

**Step 3: Commit**

```bash
git add .goreleaser.yml
git commit -m ":package: Add Homebrew tap config to goreleaser

Configure goreleaser to auto-publish a Homebrew formula to
Yoshi325/homebrew-tap on each tagged release.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 2: Pass `HOMEBREW_TAP_TOKEN` in release workflow

**Files:**
- Modify: `.github/workflows/release.yml:31-32`

**Step 1: Add the token to the env block**

In `.github/workflows/release.yml`, find the goreleaser step's `env` block (line 32) and add the `HOMEBREW_TAP_TOKEN`:

Before:
```yaml
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

After:
```yaml
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

**Step 2: Commit**

```bash
git add .github/workflows/release.yml
git commit -m ":construction_worker: Pass Homebrew tap token to goreleaser

The HOMEBREW_TAP_TOKEN secret is needed for goreleaser to push
the formula to the Yoshi325/homebrew-tap repository.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 3: Update README with Homebrew install instructions and star CTA

**Files:**
- Modify: `README.md:52-58` (Installation section)
- Modify: `README.md:72` (before License section)

**Step 1: Update the Installation section**

Replace the current Installation section (lines 52-58) with:

```markdown
## Installation

### Homebrew (macOS/Linux)

```
brew install Yoshi325/tap/grin
```

### Go

```
go install github.com/Yoshi325/grin@latest
```

### Binary

Download a binary from the [latest release](https://github.com/Yoshi325/grin/releases/latest).
```

**Step 2: Add star CTA before License**

Insert the following section before the `## License` line:

```markdown
## Help grin get into Homebrew core

If you find grin useful, please [star this repository](https://github.com/Yoshi325/grin)! Once we reach 75+ stars, we can submit grin to Homebrew core so anyone can install with just `brew install grin`.

```

**Step 3: Commit**

```bash
git add README.md
git commit -m ":memo: Add Homebrew install instructions and star CTA

Add brew install command as the primary installation method.
Include a call-to-action for stars to eventually qualify for
Homebrew core submission (requires 75+ stars).

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 4: Manual setup (GitHub UI — not code changes)

These steps must be done by the repo owner in the GitHub web UI.

**Step 1: Create the tap repository**

1. Go to https://github.com/new
2. Repository name: `homebrew-tap`
3. Set to **Public**
4. Do NOT initialize with README (leave empty)
5. Click **Create repository**

**Step 2: Create a personal access token**

1. Go to https://github.com/settings/tokens
2. Click **Generate new token (classic)**
3. Note: `goreleaser homebrew tap`
4. Expiration: choose based on preference (no expiration or 1 year)
5. Scope: check **`repo`** (full control of private repositories — needed to push to the tap repo)
6. Click **Generate token**
7. Copy the token immediately

**Step 3: Add the secret to the grin repo**

1. Go to https://github.com/Yoshi325/grin/settings/secrets/actions
2. Click **New repository secret**
3. Name: `HOMEBREW_TAP_TOKEN`
4. Value: paste the PAT from Step 2
5. Click **Add secret**

---

### Task 5: Test the full pipeline

**Step 1: Tag a release**

```bash
git tag v0.X.0  # use the next appropriate version number
git push origin v0.X.0
```

**Step 2: Verify the release workflow**

1. Go to https://github.com/Yoshi325/grin/actions — confirm the Release workflow completes successfully
2. Check the goreleaser logs for the "Publishing to Homebrew tap" step

**Step 3: Verify the formula was pushed**

1. Go to https://github.com/Yoshi325/homebrew-tap
2. Confirm a `Formula/grin.rb` file exists with correct download URLs and SHA256 checksums

**Step 4: Test the install**

```bash
brew install Yoshi325/tap/grin
grin --version
```

Expected: prints the tagged version number.

**Step 5: Test the formula test**

```bash
brew test grin
```

Expected: passes (runs the `assert_match` from the formula's test block).
