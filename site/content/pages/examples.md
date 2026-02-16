Title: Examples
Slug: examples
Summary: Detailed usage examples for grin

## Basic Conversion

Transform an INI file into greppable assignments:

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

## Filtering with grep

Find all database-related settings:

    $ grin config.ini | grep database
    ini.database = {};
    ini.database.host = "db.example.com";
    ini.database.name = "mydb";
    ini.database.pool = {};
    ini.database.pool.max = "20";
    ini.database.pool.min = "5";
    ini.database.port = "5432";

## Reconstructing Filtered INI

Pipe filtered assignments through `grin -u` to get valid INI back:

    $ grin config.ini | grep database | grin -u
    [database]
    host = db.example.com
    name = mydb
    port = 5432

    [database.pool]
    max = 20
    min = 5

## Extracting Values Only

Use `--values` to print just the values without paths or quoting:

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

## Reading from stdin

Pipe INI content directly:

    $ cat config.ini | grin

Or from another command:

    $ curl -s https://example.com/config.ini | grin | grep database

## Unsorted Output

Use `--no-sort` to preserve the original INI order (faster for large files):

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

## Comparing Config Files

Use grin with `diff` to compare two INI files structurally:

    $ diff <(grin config-prod.ini) <(grin config-staging.ini)

## PowerShell Examples

grin works great with PowerShell's `Select-String` (the `grep` equivalent):

    PS> grin -m config.ini | Select-String "database"
    ini.database = {};
    ini.database.host = "db.example.com";
    ini.database.name = "mydb";
    ini.database.pool = {};
    ini.database.pool.max = "20";
    ini.database.pool.min = "5";
    ini.database.port = "5432";

Filter with multiple patterns:

    PS> grin -m config.ini | Select-String -Pattern @("database", "cache")

Round-trip with filtering:

    PS> grin -m config.ini | Select-String "database" | grin -u
    [database]
    host = db.example.com
    name = mydb
    port = 5432

    [database.pool]
    max = 20
    min = 5

Show context around matches:

    PS> grin -m config.ini | Select-String "database" -Context 2

Compare config files in PowerShell:

    PS> Compare-Object (grin config-prod.ini) (grin config-staging.ini)
