# grin

[![CI](https://github.com/Yoshi325/grin/actions/workflows/ci.yml/badge.svg)](https://github.com/Yoshi325/grin/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Yoshi325/grin)](https://goreportcard.com/report/github.com/Yoshi325/grin)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/Yoshi325/grin.svg)](https://pkg.go.dev/github.com/Yoshi325/grin)
[![Release](https://img.shields.io/github/v/release/Yoshi325/grin)](https://github.com/Yoshi325/grin/releases/latest)
[![Website](https://img.shields.io/badge/website-yoshi325.github.io%2Fgrin-blue)](https://yoshi325.github.io/grin/)
[![Tiny Tool Town](https://img.shields.io/badge/%F0%9F%8F%98%EF%B8%8F_Tiny_Tool_Town-Resident-orange)](https://www.tinytooltown.com/)

Make INI files greppable.

grin transforms INI files into discrete assignments, making them easy to explore and filter with standard Unix tools like `grep`, `sed`, and `diff`.

Inspired by [gron](https://github.com/tomnomnom/gron).

> **Warning:** grin is currently beta software. Please make sure you have a backup before using grin on any important files.

## Usage

![Basic conversion](docs/demos/basic.gif)

<details>
<summary>Text version (click to expand)</summary>

```
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
```

</details>

Now you can use standard tools to find what you need:

![Grep filtering](docs/demos/filtering.gif)

<details>
<summary>Text version (click to expand)</summary>

```
$ grin -m testdata/complex.ini | grep database
ini.database = {};
ini.database.host = "db.example.com";
ini.database.name = "mydb";
ini.database.pool = {};
ini.database.pool.max = "20";
ini.database.pool.min = "5";
ini.database.port = "5432";
```

</details>

And reconstruct a filtered INI file with `--ungrin`:

![Ungrin round-trip](docs/demos/ungrin.gif)

<details>
<summary>Text version (click to expand)</summary>

```
$ grin -m testdata/complex.ini | grep database | grin -u
[database]
host = db.example.com
name = mydb
port = 5432

[database.pool]
max = 20
min = 5
```

</details>

### PowerShell

grin works great with PowerShell's `Select-String` (the `grep` equivalent):

```powershell
# Filter INI sections
grin -m config.ini | Select-String "database"

# Multiple patterns
grin -m config.ini | Select-String -Pattern @("database", "cache")

# Round-trip with filtering
grin -m config.ini | Select-String "database" | grin -u
```

## Installation

### Homebrew (macOS/Linux)

```
brew install Yoshi325/tap/grin
```

### Scoop (Windows)

```powershell
scoop bucket add yoshi325 https://github.com/Yoshi325/scoop-bucket
scoop install grin
```

### Go

```
go install github.com/Yoshi325/grin@latest
```

### Binary

Download a binary from the [latest release](https://github.com/Yoshi325/grin/releases/latest).

## Options

```
-u, --ungrin     Reverse the operation (turn assignments back into INI)
-v, --values     Print just the values of provided assignments
-c, --colorize   Colorize output (default on tty)
-m, --monochrome Monochrome (don't colorize output)
    --no-sort    Don't sort output (faster)
    --version    Print version information
```

## Community

> **🏘️ grin is a proud resident of [Tiny Tool Town](https://www.tinytooltown.com/)** — Scott Hanselman's neighborhood of free, fun & open source tiny tools.

If you find grin useful, please [star this repository](https://github.com/Yoshi325/grin)! Once we reach 75+ stars, we can submit grin to Homebrew core so anyone can install with just `brew install grin`.

## License

[MIT](LICENSE)
