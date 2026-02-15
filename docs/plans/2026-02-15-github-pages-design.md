# GitHub Pages Site — Design

## Summary

Create a project website for grin at `yoshi325.github.io/grin` using Pelican (Python static site generator) with the Papyrus theme. The site combines a landing page with a few documentation pages. Built and deployed via GitHub Actions to an orphan `gh-pages` branch.

## Site Structure

1. **Index (landing page)** — Custom Jinja2 template with hero section: tagline ("Make INI files greppable"), the 3 demo GIFs, install instructions, link to GitHub repo.
2. **Examples page** — Detailed usage examples: text versions of demos plus additional scenarios (`--values`, `--no-sort`, stdin piping).
3. **Man page** — The `grin.1` man page converted to Markdown and rendered as a Pelican page.

## Theme

Papyrus (Tailwind CSS, built-in dark mode, modern aesthetic). Landing page uses a custom template override for the hero section.

## Source Layout

Pelican source files live on `main` in a `site/` directory:

```
site/
  pelicanconf.py          # Local dev config (SITEURL = "")
  publishconf.py          # Production config (SITEURL = "https://yoshi325.github.io/grin")
  content/
    pages/
      examples.md
      manpage.md
  templates/
    index.html            # Custom landing page template
```

## Build & Deploy

GitHub Actions workflow (`.github/workflows/pages.yml` on `main`):

1. Triggers on pushes to `main` that touch `site/**`
2. Sets up Python, installs Pelican + plugins + Papyrus theme
3. Runs `pelican content -s publishconf.py`
4. Deploys built output to orphan `gh-pages` branch via `peaceiris/actions-gh-pages`

## Assets

- **GIF demos:** Referenced via raw GitHub URLs from `main` branch — no duplication
- **Man page:** Converted once from `grin.1` via pandoc, committed as Markdown page
- **Theme:** Installed as pip dependency at build time

## Out of Scope

- Custom domain (uses default `yoshi325.github.io/grin`)
- Blog/news section
- Auto-updating man page conversion
