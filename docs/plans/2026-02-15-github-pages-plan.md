# GitHub Pages Site — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a Pelican-powered project website for grin at `yoshi325.github.io/grin`.

**Architecture:** Pelican source in `site/` on `main`, Papyrus theme via git submodule, built and deployed by GitHub Actions to an orphan `gh-pages` branch. Custom landing page template, plus Markdown pages for examples and the man page.

**Tech Stack:** Python 3, Pelican, Papyrus theme (Tailwind CSS), pelican-search, pelican-neighbors, pelican-readtime, pelican-toc, Stork (search indexer), GitHub Actions, peaceiris/actions-gh-pages

---

### Task 1: Scaffold the Pelican project

**Files:**
- Create: `site/pelicanconf.py`
- Create: `site/publishconf.py`
- Create: `site/content/pages/.gitkeep`
- Create: `site/requirements.txt`

**Step 1: Create directory structure**

Run:
```bash
mkdir -p site/content/pages site/content/images site/templates site/pelican-plugins
```

**Step 2: Create requirements.txt**

Create `site/requirements.txt`:

```
pelican[markdown]>=4.9
pelican-search
pelican-neighbors
pelican-readtime
beautifulsoup4
```

**Step 3: Create pelicanconf.py**

Create `site/pelicanconf.py`:

```python
AUTHOR = "Charles Y."
SITENAME = "grin"
SITEURL = ""

PATH = "content"
TIMEZONE = "America/New_York"
DEFAULT_LANG = "en"

# Theme
THEME = "themes/papyrus"

# Plugins
PLUGINS = ["readtime", "search", "neighbors", "pelican-toc"]
PLUGIN_PATHS = ["pelican-plugins"]

# Search
SEARCH_MODE = "output"
SEARCH_HTML_SELECTOR = "main"

# Static paths
STATIC_PATHS = ["images"]

# URL settings
PAGE_URL = "{slug}/"
PAGE_SAVE_AS = "{slug}/index.html"

# Disable blog features (this is a project site, not a blog)
DIRECT_TEMPLATES = ["archives"]
ARTICLE_PATHS = []
DISPLAY_PAGES_ON_MENU = True

# Template overrides for custom landing page
THEME_TEMPLATES_OVERRIDES = ["templates"]

# Feed (disable for local dev)
FEED_ALL_ATOM = None
CATEGORY_FEED_ATOM = None
TRANSLATION_FEED_ATOM = None
AUTHOR_FEED_ATOM = None
AUTHOR_FEED_RSS = None
```

**Step 4: Create publishconf.py**

Create `site/publishconf.py`:

```python
import os
import sys

sys.path.append(os.curdir)
from pelicanconf import *  # noqa: F401,F403,E402

SITEURL = "https://yoshi325.github.io/grin"
RELATIVE_URLS = False

FEED_ALL_ATOM = "feeds/all.atom.xml"
```

**Step 5: Commit**

```bash
git add site/
git commit -m ":sparkles: Scaffold Pelican project structure

Basic Pelican config with Papyrus theme settings, requirements,
and directory structure for the GitHub Pages site.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 2: Add Papyrus theme and pelican-toc plugin as git submodules

**Files:**
- Create: `site/themes/papyrus` (submodule)
- Create: `site/pelican-plugins/pelican-toc` (submodule)

**Step 1: Add Papyrus theme as submodule**

Run:
```bash
cd /Users/charles/Developer/GitHub/Yoshi325-Grin
git submodule add https://github.com/pelican-themes/papyrus.git site/themes/papyrus
```

**Step 2: Add pelican-toc plugin as submodule**

Run:
```bash
git submodule add https://github.com/ingwinlu/pelican-toc.git site/pelican-plugins/pelican-toc
```

**Step 3: Commit**

```bash
git add .gitmodules site/themes/papyrus site/pelican-plugins/pelican-toc
git commit -m ":sparkles: Add Papyrus theme and pelican-toc as submodules

Papyrus provides the Tailwind CSS theme with dark mode. pelican-toc
is required by Papyrus but not available via pip.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 3: Create the custom landing page template

**Files:**
- Create: `site/templates/index.html`

**Step 1: Write the custom index template**

Create `site/templates/index.html`. This extends the Papyrus base template and provides a hero section with tagline, demo GIFs, install instructions, and GitHub link:

```html
{% extends "base.html" %}

{% block title %}{{ SITENAME }} — Make INI files greppable{% endblock %}

{% block content %}
<div class="mx-auto max-w-3xl px-4">
  <header class="text-center py-12">
    <h1 class="text-5xl font-bold mb-4">grin</h1>
    <p class="text-xl text-gray-600 dark:text-gray-400 mb-6">Make INI files greppable.</p>
    <p class="text-gray-600 dark:text-gray-400 mb-8">
      grin transforms INI files into discrete assignments, making them easy to explore
      and filter with standard Unix tools like <code>grep</code>, <code>sed</code>, and <code>diff</code>.
      Inspired by <a href="https://github.com/tomnomnom/gron" class="underline">gron</a>.
    </p>
    <a href="https://github.com/Yoshi325/grin"
       class="inline-block bg-gray-800 dark:bg-gray-200 text-white dark:text-gray-800 px-6 py-3 rounded-lg font-medium hover:opacity-90 transition">
      View on GitHub
    </a>
  </header>

  <section class="py-8">
    <h2 class="text-2xl font-bold mb-4">Convert INI to greppable assignments</h2>
    <img src="https://raw.githubusercontent.com/Yoshi325/grin/main/docs/demos/basic.gif"
         alt="grin basic conversion demo"
         class="rounded-lg shadow-lg mb-8 w-full" loading="lazy">
  </section>

  <section class="py-8">
    <h2 class="text-2xl font-bold mb-4">Filter with standard Unix tools</h2>
    <img src="https://raw.githubusercontent.com/Yoshi325/grin/main/docs/demos/filtering.gif"
         alt="grin grep filtering demo"
         class="rounded-lg shadow-lg mb-8 w-full" loading="lazy">
  </section>

  <section class="py-8">
    <h2 class="text-2xl font-bold mb-4">Reconstruct filtered INI files</h2>
    <img src="https://raw.githubusercontent.com/Yoshi325/grin/main/docs/demos/ungrin.gif"
         alt="grin ungrin round-trip demo"
         class="rounded-lg shadow-lg mb-8 w-full" loading="lazy">
  </section>

  <section class="py-8">
    <h2 class="text-2xl font-bold mb-4">Install</h2>
    <div class="space-y-4">
      <div>
        <h3 class="font-semibold mb-2">Homebrew (macOS/Linux)</h3>
        <pre class="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 overflow-x-auto"><code>brew install Yoshi325/tap/grin</code></pre>
      </div>
      <div>
        <h3 class="font-semibold mb-2">Go</h3>
        <pre class="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 overflow-x-auto"><code>go install github.com/Yoshi325/grin@latest</code></pre>
      </div>
      <div>
        <h3 class="font-semibold mb-2">Binary</h3>
        <p>Download from the <a href="https://github.com/Yoshi325/grin/releases/latest" class="underline">latest release</a>.</p>
      </div>
    </div>
  </section>
</div>
{% endblock %}
```

**Step 2: Commit**

```bash
git add site/templates/index.html
git commit -m ":sparkles: Add custom landing page template

Hero section with tagline, demo GIFs from main branch, install
instructions, and GitHub link. Uses Tailwind CSS classes from Papyrus.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 4: Create the Examples page

**Files:**
- Create: `site/content/pages/examples.md`

**Step 1: Write the examples page**

Create `site/content/pages/examples.md`:

```markdown
Title: Examples
Slug: examples
Summary: Detailed usage examples for grin

## Basic Conversion

Transform an INI file into greppable assignments:

```
$ cat config.ini
; Complex configuration example
app-name = SuperApp
version = 2.1

[database]
host = db.example.com
port = 5432
name = mydb

[database.pool]
min = 5
max = 20

[cache]
enabled = true
ttl = 3600

[logging]
level = info
file = /var/log/app.log

$ grin config.ini
ini = {};
ini.app-name = "SuperApp";
ini.cache = {};
ini.cache.enabled = "true";
ini.cache.ttl = "3600";
ini.database = {};
ini.database.host = "db.example.com";
ini.database.name = "mydb";
ini.database.pool = {};
ini.database.pool.max = "20";
ini.database.pool.min = "5";
ini.database.port = "5432";
ini.logging = {};
ini.logging.file = "/var/log/app.log";
ini.logging.level = "info";
ini.version = "2.1";
```

## Filtering with grep

Find all database-related settings:

```
$ grin config.ini | grep database
ini.database = {};
ini.database.host = "db.example.com";
ini.database.name = "mydb";
ini.database.pool = {};
ini.database.pool.max = "20";
ini.database.pool.min = "5";
ini.database.port = "5432";
```

## Reconstructing Filtered INI

Pipe filtered assignments through `grin -u` to get valid INI back:

```
$ grin config.ini | grep database | grin -u
[database]
host = db.example.com
name = mydb
port = 5432

[database.pool]
max = 20
min = 5
```

## Extracting Values Only

Use `--values` to print just the values without paths or quoting:

```
$ grin -v config.ini
SuperApp
2.1
db.example.com
5432
mydb
5
20
true
3600
info
/var/log/app.log
```

## Reading from stdin

Pipe INI content directly:

```
$ cat config.ini | grin
```

Or from another command:

```
$ curl -s https://example.com/config.ini | grin | grep database
```

## Unsorted Output

Use `--no-sort` to preserve the original INI order (faster for large files):

```
$ grin --no-sort config.ini
ini = {};
ini.app-name = "SuperApp";
ini.version = "2.1";
ini.database = {};
ini.database.host = "db.example.com";
ini.database.port = "5432";
ini.database.name = "mydb";
ini.database.pool = {};
ini.database.pool.min = "5";
ini.database.pool.max = "20";
ini.cache = {};
ini.cache.enabled = "true";
ini.cache.ttl = "3600";
ini.logging = {};
ini.logging.level = "info";
ini.logging.file = "/var/log/app.log";
```

## Comparing Config Files

Use grin with `diff` to compare two INI files structurally:

```
$ diff <(grin config-prod.ini) <(grin config-staging.ini)
```
```

**Step 2: Commit**

```bash
git add site/content/pages/examples.md
git commit -m ":sparkles: Add examples page for project site

Detailed usage examples covering basic conversion, grep filtering,
ungrin round-trip, values-only mode, stdin, no-sort, and diff.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 5: Create the Man Page as a Pelican page

**Files:**
- Create: `site/content/pages/manpage.md`

**Step 1: Convert man page to Markdown**

Run pandoc to convert `doc/grin.1` to Markdown, then use the output to create the page. If pandoc is not installed, convert manually from the man page source.

Run:
```bash
pandoc doc/grin.1 -f man -t markdown --standalone 2>/dev/null || echo "PANDOC_NOT_INSTALLED"
```

**Step 2: Create the man page Markdown file**

Create `site/content/pages/manpage.md` with Pelican metadata header and the converted content:

```markdown
Title: Man Page
Slug: manpage
Summary: grin(1) manual page

## NAME

grin — make INI files greppable

## SYNOPSIS

**grin** [*OPTIONS*] [*FILE* | **-**]

## DESCRIPTION

**grin** transforms INI files into discrete assignments, making them
easy to explore and filter with standard Unix tools like **grep**(1),
**sed**(1), and **diff**(1).

Each INI key becomes a dot-separated path rooted at **ini**, with
string values quoted and lines terminated by semicolons. Section
headers produce empty-object assignments (`ini.section = {};`) so the
full structure is preserved.

When invoked with **--ungrin** (or as **ungrin**), the process is
reversed: assignment lines are parsed back into INI format.

If no *FILE* is given, or if *FILE* is **-**, **grin** reads from
standard input.

## OPTIONS

**-u**, **--ungrin**
:   Reverse the operation: turn grin assignments back into INI format.

**-v**, **--values**
:   Print just the values of provided assignments, one per line, without paths or quoting.

**-c**, **--colorize**
:   Colorize output. This is the default when output is a terminal.

**-m**, **--monochrome**
:   Monochrome output. Disables colorization.

**--no-sort**
:   Don't sort output. Preserves the original order of the INI file, which can be faster for large files.

**--version**
:   Print version information and exit.

## EXAMPLES

Transform an INI file into greppable assignments:

```
$ grin config.ini
ini = {};
ini.database = {};
ini.database.host = "localhost";
ini.database.port = "5432";
```

Filter for specific keys using grep:

```
$ grin config.ini | grep host
ini.database.host = "localhost";
```

Reconstruct a filtered INI file:

```
$ grin config.ini | grep pool | grin --ungrin
[database.pool]
max = 10
```

Read from standard input:

```
$ cat config.ini | grin
```

Extract just the values:

```
$ grin -v config.ini
localhost
5432
```

## EXIT STATUS

**0**
:   Success.

**1**
:   Failed to open the specified file.

**2**
:   Failed to read input.

**3**
:   Failed to form statements from INI input.

**5**
:   Failed to parse assignment statements (during ungrin).

## SEE ALSO

**gron**(1), **grep**(1), **sed**(1), **diff**(1), **ini**(5)

## AUTHOR

Charles Y. (Yoshi325)

## LICENSE

MIT License. See the LICENSE file in the source distribution.

## BUGS

Report bugs at: [github.com/Yoshi325/grin/issues](https://github.com/Yoshi325/grin/issues)
```

**Step 3: Commit**

```bash
git add site/content/pages/manpage.md
git commit -m ":sparkles: Add man page as Pelican page

grin(1) man page converted to Markdown for the project site.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 6: Create the GitHub Actions workflow

**Files:**
- Create: `.github/workflows/pages.yml`

**Step 1: Write the workflow**

Create `.github/workflows/pages.yml`:

```yaml
name: GitHub Pages

on:
  push:
    branches: [main]
    paths:
      - 'site/**'
      - '.github/workflows/pages.yml'
  workflow_dispatch:

permissions:
  contents: write

jobs:
  deploy:
    name: Build and Deploy
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          submodules: true

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.12'

      - name: Install dependencies
        working-directory: site
        run: pip install -r requirements.txt

      - name: Install Stork (for search)
        run: |
          wget -qO stork https://files.stork-search.net/releases/v1.6.0/stork-amazon-linux
          chmod +x stork
          sudo mv stork /usr/local/bin/

      - name: Build site
        working-directory: site
        run: pelican content -s publishconf.py

      - name: Deploy to gh-pages
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: site/output
```

**Step 2: Commit**

```bash
git add .github/workflows/pages.yml
git commit -m ":construction_worker: Add GitHub Pages build and deploy workflow

Builds Pelican site and deploys to gh-pages branch on pushes to
main that touch site/ files. Uses peaceiris/actions-gh-pages.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 7: Test locally and verify

**Step 1: Set up local Python environment**

Run:
```bash
cd /Users/charles/Developer/GitHub/Yoshi325-Grin/site
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

**Step 2: Build the site locally**

Run:
```bash
cd /Users/charles/Developer/GitHub/Yoshi325-Grin/site
pelican content -s pelicanconf.py
```
Expected: `output/` directory created with HTML files

**Step 3: Verify output files exist**

Run:
```bash
ls site/output/index.html site/output/examples/index.html site/output/manpage/index.html
```
Expected: all three files exist

**Step 4: Preview locally**

Run:
```bash
cd /Users/charles/Developer/GitHub/Yoshi325-Grin/site
pelican --listen
```
Expected: site served at http://localhost:8000, visually verify landing page, examples, and man page render correctly

**Step 5: Add output/ to .gitignore**

Create or update `site/.gitignore`:

```
output/
venv/
__pycache__/
*.pyc
```

**Step 6: Commit**

```bash
git add site/.gitignore
git commit -m ":see_no_evil: Add site .gitignore for build artifacts

Ignore Pelican output directory, Python venv, and cache files.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 8: Push and enable GitHub Pages

**Step 1: Push all commits**

Run: `git push`

**Step 2: Verify the workflow runs**

Run: `gh run list --workflow=pages.yml --limit=1`
Expected: workflow triggered and running/succeeded

**Step 3: Enable GitHub Pages in repo settings**

Run:
```bash
gh api repos/Yoshi325/grin/pages -X PUT -f source.branch=gh-pages -f source.path="/" 2>/dev/null || \
gh api repos/Yoshi325/grin/pages -X POST -f source.branch=gh-pages -f source.path="/"
```

**Step 4: Verify the site is live**

Run: `curl -sI https://yoshi325.github.io/grin/ | head -5`
Expected: HTTP 200 response (may take a few minutes after first deploy)
