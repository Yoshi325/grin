package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

const (
	exitOK             = 0
	exitOpenFile       = 1
	exitReadInput      = 2
	exitFormStatements = 3
	exitParseStatements = 5
)

const (
	optMonochrome = 1 << iota
	optNoSort
)

var grinVersion = "dev"

type actionFn func(io.Reader, io.Writer, int) (int, error)

func main() {
	var (
		ungrinFlag     bool
		colorizeFlag   bool
		monochromeFlag bool
		noSortFlag     bool
		versionFlag    bool
		valuesFlag     bool
	)

	flag.BoolVar(&ungrinFlag, "ungrin", false, "Reverse the operation (turn assignments back into INI)")
	flag.BoolVar(&ungrinFlag, "u", false, "Reverse the operation (turn assignments back into INI)")
	flag.BoolVar(&colorizeFlag, "colorize", false, "Colorize output (default on tty)")
	flag.BoolVar(&colorizeFlag, "c", false, "Colorize output (default on tty)")
	flag.BoolVar(&monochromeFlag, "monochrome", false, "Monochrome (don't colorize output)")
	flag.BoolVar(&monochromeFlag, "m", false, "Monochrome (don't colorize output)")
	flag.BoolVar(&noSortFlag, "no-sort", false, "Don't sort output (faster)")
	flag.BoolVar(&versionFlag, "version", false, "Print version information")
	flag.BoolVar(&valuesFlag, "values", false, "Print just the values of provided assignments")
	flag.BoolVar(&valuesFlag, "v", false, "Print just the values of provided assignments")

	flag.Usage = func() {
		h := "Transform INI (from a file or stdin) into discrete assignments to make it greppable\n\n"

		h += "Usage:\n"
		h += "  grin [OPTIONS] [FILE|-]\n\n"

		h += "Options:\n"
		h += "  -u, --ungrin     Reverse the operation (turn assignments back into INI)\n"
		h += "  -v, --values     Print just the values of provided assignments\n"
		h += "  -c, --colorize   Colorize output (default on tty)\n"
		h += "  -m, --monochrome Monochrome (don't colorize output)\n"
		h += "      --no-sort    Don't sort output (faster)\n"
		h += "      --version    Print version information\n\n"

		h += "Exit Codes:\n"
		h += "  0\tOK\n"
		h += "  1\tFailed to open file\n"
		h += "  2\tFailed to read input\n"
		h += "  3\tFailed to form statements\n"
		h += "  5\tFailed to parse statements\n\n"

		h += "Examples:\n"
		h += "  grin /etc/config.ini\n"
		h += "  grin config.ini | grep database\n"
		h += "  cat config.ini | grin\n"
		h += "  grin config.ini | grep host | grin --ungrin\n"

		fmt.Fprint(os.Stderr, h)
	}

	flag.Parse()

	if versionFlag {
		fmt.Printf("grin version %s\n", grinVersion)
		os.Exit(exitOK)
	}

	// Auto-detect ungrin mode if invoked as "ungrin"
	if strings.HasSuffix(os.Args[0], "ungrin") {
		ungrinFlag = true
	}

	// Determine input source
	var rawInput io.Reader
	filename := flag.Arg(0)
	if filename == "" || filename == "-" {
		rawInput = os.Stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			fatal(exitOpenFile, err)
		}
		defer f.Close()
		rawInput = f
	}

	// Build options
	var opts int
	switch {
	case colorizeFlag:
		color.NoColor = false
	case monochromeFlag || color.NoColor:
		opts |= optMonochrome
	}
	if noSortFlag {
		opts |= optNoSort
	}

	// Select action
	var a actionFn = grinAction
	if ungrinFlag {
		a = ungrinAction
	} else if valuesFlag {
		a = grinValuesAction
	}

	exitCode, err := a(rawInput, colorable.NewColorableStdout(), opts)
	if exitCode != exitOK {
		fatal(exitCode, err)
	}
	os.Exit(exitOK)
}

func grinAction(r io.Reader, w io.Writer, opts int) (int, error) {
	var conv statementconv
	if opts&optMonochrome > 0 {
		conv = statementToString
	} else {
		conv = statementToColorString
	}

	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(r, prefix)
	if err != nil {
		return exitFormStatements, err
	}

	if opts&optNoSort == 0 {
		sort.Sort(ss)
	}

	for _, s := range ss {
		fmt.Fprintln(w, conv(s))
	}

	return exitOK, nil
}

func ungrinAction(r io.Reader, w io.Writer, opts int) (int, error) {
	ss, err := ungrinStatements(r)
	if err != nil {
		return exitParseStatements, err
	}

	if err := ungrinFromStatements(ss, w); err != nil {
		return exitParseStatements, err
	}

	return exitOK, nil
}

func grinValuesAction(r io.Reader, w io.Writer, opts int) (int, error) {
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(r, prefix)
	if err != nil {
		return exitFormStatements, err
	}

	for _, s := range ss {
		for _, t := range s {
			if t.typ == typString {
				fmt.Fprintln(w, unquoteString(t.text))
			}
		}
	}

	return exitOK, nil
}

func fatal(code int, err error) {
	fmt.Fprintf(os.Stderr, "grin: %s\n", err)
	os.Exit(code)
}
