# VHS Demos & README Update — Design

## Summary

Add animated terminal demos to the README using [VHS by Charmbracelet](https://github.com/charmbracelet/vhs). Three scriptable `.tape` files produce GIF recordings that replace the static code blocks in the Usage section, with the original code preserved in collapsible `<details>` sections.

## Demos

Three tape files in `docs/demos/`, each producing a GIF in the same directory:

1. **`basic.tape`** → `basic.gif` — `cat testdata/complex.ini` then `grin testdata/complex.ini`. Core INI-to-assignments transformation.
2. **`filtering.tape`** → `filtering.gif` — `grin testdata/complex.ini | grep database`. The "greppable" value proposition.
3. **`ungrin.tape`** → `ungrin.gif` — `grin testdata/complex.ini | grep database | grin -u`. Round-trip: filter then reconstruct valid INI.

All tapes use `-m` (monochrome) so output renders cleanly without terminal color escape issues.

## VHS Settings (shared across all tapes)

- Theme: dark (VHS default or Dracula)
- Font size: 20–22
- Dimensions: 1200x600
- Typing speed: fast but readable, short pauses between commands and after output
- Monochrome flag on grin commands

## README Changes

The Usage section restructured:

1. GIF for basic conversion replaces first code block
2. GIF for grep filtering replaces the grep example code block
3. GIF for ungrin round-trip replaces the ungrin example code block

Each GIF followed by a collapsible `<details>` block containing the equivalent code for copy-paste and accessibility.

Example data shifts from the current simple `config.ini` snippet to `testdata/complex.ini` (richer: app globals, database, cache, logging).

## Out of Scope

- No CI workflow for GIF regeneration — generate locally, commit the GIFs
- No new example INI files — using existing `testdata/complex.ini`
- No changes to grin functionality or flags
