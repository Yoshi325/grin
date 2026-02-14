# Homebrew Tap Distribution Design

## Goal

Publish grin on Homebrew via a personal tap (`Yoshi325/homebrew-tap`) so users can install with `brew install Yoshi325/tap/grin`. Include a call-to-action for GitHub stars to eventually qualify for Homebrew core (requires 75+ stars).

## Approach

Use goreleaser's built-in `brews` section to automatically generate and push a Homebrew formula to a dedicated tap repository on every tagged release.

## Changes

### 1. `.goreleaser.yml` — add `brews` section

Add after the existing `changelog` block:

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

goreleaser will generate a `Formula/grin.rb` with correct download URLs and SHA256 checksums, then push it to the tap repo.

### 2. `.github/workflows/release.yml` — pass token

Add `HOMEBREW_TAP_TOKEN` to the goreleaser step's env block:

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

### 3. `README.md` — installation instructions and CTA

Update the Installation section to list Homebrew first, followed by Go install and binary download. Add a "Help grin get into Homebrew core" section before License with a call-to-action for stars.

### 4. Manual steps (outside this repo)

- Create the `Yoshi325/homebrew-tap` GitHub repository (public, empty)
- Create a GitHub personal access token (classic) with `repo` scope
- Add the token as `HOMEBREW_TAP_TOKEN` in the grin repo's Settings > Secrets and variables > Actions

## How it works end-to-end

1. Push a `v*` tag
2. GitHub Actions runs the release workflow
3. goreleaser builds binaries, creates the GitHub release with archives and checksums
4. goreleaser generates `Formula/grin.rb` and pushes it to `Yoshi325/homebrew-tap`
5. Users run `brew install Yoshi325/tap/grin`

## Future: Homebrew core

Once the repo reaches 75+ GitHub stars, submit the formula to `Homebrew/homebrew-core` via a PR. At that point users can install with just `brew install grin`.
