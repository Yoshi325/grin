# VHS Demos & README Update — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add three animated terminal demo GIFs to the README using VHS tape files.

**Architecture:** Three `.tape` files in `docs/demos/` produce GIFs showing basic conversion, grep filtering, and ungrin round-trip. The README Usage section replaces static code blocks with GIFs, preserving code in collapsible `<details>` blocks.

**Tech Stack:** VHS (charmbracelet/vhs), grin binary (built from source)

---

### Task 1: Install VHS and build grin

**Files:**
- None (environment setup only)

**Step 1: Verify VHS is installed (or install it)**

Run: `which vhs || brew install vhs`
Expected: path to `vhs` binary

**Step 2: Build grin from source**

Run: `cd /Users/charles/Developer/GitHub/Yoshi325-Grin && go build -o grin .`
Expected: `grin` binary in project root

**Step 3: Verify grin works with test data**

Run: `./grin -m testdata/complex.ini`
Expected: 16 lines of monochrome assignments starting with `ini = {};`

---

### Task 2: Create basic.tape

**Files:**
- Create: `docs/demos/basic.tape`

**Step 1: Create the docs/demos directory**

Run: `mkdir -p /Users/charles/Developer/GitHub/Yoshi325-Grin/docs/demos`

**Step 2: Write basic.tape**

Create `docs/demos/basic.tape`:

```
Output docs/demos/basic.gif

Set Shell "bash"
Set FontSize 22
Set Width 1200
Set Height 600
Set Theme "Dracula"
Set TypingSpeed 50ms

Type "cat testdata/complex.ini"
Sleep 500ms
Enter
Sleep 2s

Type "grin -m testdata/complex.ini"
Sleep 500ms
Enter
Sleep 3s
```

**Step 3: Generate the GIF**

Run: `cd /Users/charles/Developer/GitHub/Yoshi325-Grin && vhs docs/demos/basic.tape`
Expected: `docs/demos/basic.gif` created

**Step 4: Verify the GIF exists and is reasonable size**

Run: `ls -lh docs/demos/basic.gif`
Expected: file exists, size in hundreds of KB to low MB range

**Step 5: Commit**

```bash
git add docs/demos/basic.tape docs/demos/basic.gif
git commit -m ":sparkles: Add basic conversion VHS demo

Animated terminal recording showing cat of complex.ini followed
by grin output, generated via VHS tape file.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 3: Create filtering.tape

**Files:**
- Create: `docs/demos/filtering.tape`

**Step 1: Write filtering.tape**

Create `docs/demos/filtering.tape`:

```
Output docs/demos/filtering.gif

Set Shell "bash"
Set FontSize 22
Set Width 1200
Set Height 600
Set Theme "Dracula"
Set TypingSpeed 50ms

Type "grin -m testdata/complex.ini | grep database"
Sleep 500ms
Enter
Sleep 3s
```

**Step 2: Generate the GIF**

Run: `cd /Users/charles/Developer/GitHub/Yoshi325-Grin && vhs docs/demos/filtering.tape`
Expected: `docs/demos/filtering.gif` created

**Step 3: Commit**

```bash
git add docs/demos/filtering.tape docs/demos/filtering.gif
git commit -m ":sparkles: Add grep filtering VHS demo

Animated terminal recording showing grin piped through grep to
filter for database-related config entries.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 4: Create ungrin.tape

**Files:**
- Create: `docs/demos/ungrin.tape`

**Step 1: Write ungrin.tape**

Create `docs/demos/ungrin.tape`:

```
Output docs/demos/ungrin.gif

Set Shell "bash"
Set FontSize 22
Set Width 1200
Set Height 600
Set Theme "Dracula"
Set TypingSpeed 50ms

Type "grin -m testdata/complex.ini | grep database | grin -u"
Sleep 500ms
Enter
Sleep 3s
```

**Step 2: Generate the GIF**

Run: `cd /Users/charles/Developer/GitHub/Yoshi325-Grin && vhs docs/demos/ungrin.tape`
Expected: `docs/demos/ungrin.gif` created

**Step 3: Commit**

```bash
git add docs/demos/ungrin.tape docs/demos/ungrin.gif
git commit -m ":sparkles: Add ungrin round-trip VHS demo

Animated terminal recording showing grin piped through grep then
back through ungrin to reconstruct a filtered INI file.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 5: Update README.md

**Files:**
- Modify: `README.md:17-50`

**Step 1: Replace the Usage section**

Replace the entire Usage section (lines 17–50) with GIFs and collapsible code blocks. The new section should be:

```markdown
## Usage

![Basic conversion](docs/demos/basic.gif)

<details>
<summary>Text version (click to expand)</summary>

\`\`\`
$ cat testdata/complex.ini
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

$ grin -m testdata/complex.ini
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
\`\`\`

</details>

Now you can use standard tools to find what you need:

![Grep filtering](docs/demos/filtering.gif)

<details>
<summary>Text version (click to expand)</summary>

\`\`\`
$ grin -m testdata/complex.ini | grep database
ini.database = {};
ini.database.host = "db.example.com";
ini.database.name = "mydb";
ini.database.pool = {};
ini.database.pool.max = "20";
ini.database.pool.min = "5";
ini.database.port = "5432";
\`\`\`

</details>

And reconstruct a filtered INI file with `--ungrin`:

![Ungrin round-trip](docs/demos/ungrin.gif)

<details>
<summary>Text version (click to expand)</summary>

\`\`\`
$ grin -m testdata/complex.ini | grep database | grin -u
[database]
host = db.example.com
name = mydb
port = 5432

[database.pool]
max = 20
min = 5
\`\`\`

</details>
```

**Step 2: Verify the README renders correctly**

Run: `head -80 README.md`
Expected: GIF references and details blocks visible

**Step 3: Commit**

```bash
git add README.md
git commit -m ":memo: Update README with animated VHS demos

Replace static code blocks with animated GIF demos, preserving
code in collapsible details blocks for accessibility and copy-paste.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 6: Final review

**Step 1: Verify all files are committed**

Run: `git status`
Expected: clean working tree

**Step 2: Verify all GIFs exist**

Run: `ls -lh docs/demos/*.gif`
Expected: three GIF files (basic.gif, filtering.gif, ungrin.gif)

**Step 3: Verify all tape files exist**

Run: `ls docs/demos/*.tape`
Expected: three tape files (basic.tape, filtering.tape, ungrin.tape)
