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

    $ grin config.ini
    ini = {};
    ini.database = {};
    ini.database.host = "localhost";
    ini.database.port = "5432";

Filter for specific keys using grep:

    $ grin config.ini | grep host
    ini.database.host = "localhost";

Reconstruct a filtered INI file:

    $ grin config.ini | grep pool | grin --ungrin
    [database.pool]
    max = 10

Read from standard input:

    $ cat config.ini | grin

Extract just the values:

    $ grin -v config.ini
    localhost
    5432

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
