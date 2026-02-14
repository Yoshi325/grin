# grin

[![CI](https://github.com/Yoshi325/grin/actions/workflows/ci.yml/badge.svg)](https://github.com/Yoshi325/grin/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Yoshi325/grin)](https://goreportcard.com/report/github.com/Yoshi325/grin)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/Yoshi325/grin.svg)](https://pkg.go.dev/github.com/Yoshi325/grin)
[![Release](https://img.shields.io/github/v/release/Yoshi325/grin)](https://github.com/Yoshi325/grin/releases/latest)

Make INI files greppable.

grin transforms INI files into discrete assignments, making them easy to explore and filter with standard Unix tools like `grep`, `sed`, and `diff`.

Inspired by [gron](https://github.com/tomnomnom/gron).

> **Warning:** grin is currently beta software. Please make sure you have a backup before using grin on any important files.

## Usage

```
$ cat config.ini
[database]
host = localhost
port = 5432

[database.pool]
max = 10

$ grin config.ini
ini = {};
ini.database = {};
ini.database.host = "localhost";
ini.database.pool = {};
ini.database.pool.max = "10";
ini.database.port = "5432";
```

Now you can use standard tools to find what you need:

```
$ grin config.ini | grep host
ini.database.host = "localhost";
```

And reconstruct a filtered INI file with `--ungrin`:

```
$ grin config.ini | grep pool | grin --ungrin
[database.pool]
max = 10
```

## Installation

Download a binary from the [latest release](https://github.com/Yoshi325/grin/releases/latest), or install with Go:

```
go install github.com/Yoshi325/grin@latest
```

## Options

```
-u, --ungrin     Reverse the operation (turn assignments back into INI)
-v, --values     Print just the values of provided assignments
-c, --colorize   Colorize output (default on tty)
-m, --monochrome Monochrome (don't colorize output)
    --no-sort    Don't sort output (faster)
    --version    Print version information
```

## License

[MIT](LICENSE)
