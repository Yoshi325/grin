package main

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestGrinFixtures(t *testing.T) {
	fixtures := []string{
		"simple",
		"global",
		"nested",
		"comments",
		"quoted",
		"empty",
		"complex",
	}

	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			iniPath := filepath.Join("testdata", name+".ini")
			grinPath := filepath.Join("testdata", name+".grin")

			iniData, err := os.ReadFile(iniPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", iniPath, err)
			}

			wantData, err := os.ReadFile(grinPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", grinPath, err)
			}

			var buf bytes.Buffer
			opts := optMonochrome
			exitCode, err := grinAction(bytes.NewReader(iniData), &buf, opts)
			if err != nil {
				t.Fatalf("grinAction error: %v", err)
			}
			if exitCode != exitOK {
				t.Fatalf("grinAction exit code: %d", exitCode)
			}

			got := buf.String()
			want := string(wantData)
			if got != want {
				t.Errorf("grin output mismatch for %s:\ngot:\n%s\nwant:\n%s", name, got, want)
			}
		})
	}
}

func TestGrinNoSort(t *testing.T) {
	input := `[zebra]
key = val

[alpha]
key = val
`
	var buf bytes.Buffer
	opts := optMonochrome | optNoSort
	exitCode, err := grinAction(strings.NewReader(input), &buf, opts)
	if err != nil {
		t.Fatalf("grinAction error: %v", err)
	}
	if exitCode != exitOK {
		t.Fatalf("grinAction exit code: %d", exitCode)
	}

	// With --no-sort, zebra should appear before alpha (insertion order)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	zebraIdx := -1
	alphaIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "zebra") && zebraIdx == -1 {
			zebraIdx = i
		}
		if strings.Contains(line, "alpha") && alphaIdx == -1 {
			alphaIdx = i
		}
	}

	if zebraIdx == -1 || alphaIdx == -1 {
		t.Fatal("expected both zebra and alpha in output")
	}
	if zebraIdx > alphaIdx {
		t.Errorf("with --no-sort, zebra (idx %d) should appear before alpha (idx %d)", zebraIdx, alphaIdx)
	}
}

func TestGrinValuesAction(t *testing.T) {
	input := `[section]
key1 = value1
key2 = value2
key3 = value3
`
	var buf bytes.Buffer
	exitCode, err := grinValuesAction(strings.NewReader(input), &buf, optMonochrome)
	if err != nil {
		t.Fatalf("grinValuesAction error: %v", err)
	}
	if exitCode != exitOK {
		t.Fatalf("grinValuesAction exit code: %d", exitCode)
	}

	got := strings.TrimSpace(buf.String())
	lines := strings.Split(got, "\n")
	sort.Strings(lines)

	want := []string{"value1", "value2", "value3"}
	if len(lines) != len(want) {
		t.Fatalf("got %d values, want %d", len(lines), len(want))
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("value[%d] = %q, want %q", i, lines[i], w)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := []string{
		"simple",
		"global",
		"nested",
		"quoted",
		"complex",
	}

	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			iniPath := filepath.Join("testdata", name+".ini")
			iniData, err := os.ReadFile(iniPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", iniPath, err)
			}

			// INI -> grin
			var grinBuf bytes.Buffer
			exitCode, err := grinAction(bytes.NewReader(iniData), &grinBuf, optMonochrome)
			if err != nil {
				t.Fatalf("grinAction error: %v", err)
			}
			if exitCode != exitOK {
				t.Fatalf("grinAction exit code: %d", exitCode)
			}

			firstGrin := grinBuf.String()

			// grin -> ungrin
			var unBuf bytes.Buffer
			exitCode, err = ungrinAction(strings.NewReader(firstGrin), &unBuf, optMonochrome)
			if err != nil {
				t.Fatalf("ungrinAction error: %v", err)
			}
			if exitCode != exitOK {
				t.Fatalf("ungrinAction exit code: %d", exitCode)
			}

			// ungrin -> grin again
			var secondGrinBuf bytes.Buffer
			exitCode, err = grinAction(strings.NewReader(unBuf.String()), &secondGrinBuf, optMonochrome)
			if err != nil {
				t.Fatalf("second grinAction error: %v", err)
			}
			if exitCode != exitOK {
				t.Fatalf("second grinAction exit code: %d", exitCode)
			}

			// The two grin outputs should be identical
			if firstGrin != secondGrinBuf.String() {
				t.Errorf("round-trip mismatch for %s:\nfirst grin:\n%s\nsecond grin:\n%s",
					name, firstGrin, secondGrinBuf.String())
			}
		})
	}
}

func TestUngrinAction(t *testing.T) {
	input := `ini = {};
ini.section = {};
ini.section.host = "localhost";
ini.section.port = "8080";
`
	var buf bytes.Buffer
	exitCode, err := ungrinAction(strings.NewReader(input), &buf, optMonochrome)
	if err != nil {
		t.Fatalf("ungrinAction error: %v", err)
	}
	if exitCode != exitOK {
		t.Fatalf("ungrinAction exit code: %d", exitCode)
	}

	want := "[section]\nhost = localhost\nport = 8080\n"
	if buf.String() != want {
		t.Errorf("ungrinAction output = %q, want %q", buf.String(), want)
	}
}

func TestGrinActionInvalidInput(t *testing.T) {
	input := "[invalid\nkey = value\n"
	var buf bytes.Buffer
	exitCode, err := grinAction(strings.NewReader(input), &buf, optMonochrome)
	if exitCode != exitFormStatements {
		t.Errorf("expected exit code %d, got %d", exitFormStatements, exitCode)
	}
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUngrinActionInvalidInput(t *testing.T) {
	input := "this is not valid grin output"
	var buf bytes.Buffer
	exitCode, err := ungrinAction(strings.NewReader(input), &buf, optMonochrome)
	if exitCode != exitParseStatements {
		t.Errorf("expected exit code %d, got %d", exitParseStatements, exitCode)
	}
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGrinFileInput(t *testing.T) {
	iniPath := filepath.Join("testdata", "simple.ini")
	grinPath := filepath.Join("testdata", "simple.grin")

	f, err := os.Open(iniPath)
	if err != nil {
		t.Fatalf("failed to open %s: %v", iniPath, err)
	}
	defer f.Close() //nolint:errcheck // test cleanup

	wantData, err := os.ReadFile(grinPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", grinPath, err)
	}

	var buf bytes.Buffer
	exitCode, err := grinAction(f, &buf, optMonochrome)
	if err != nil {
		t.Fatalf("grinAction error: %v", err)
	}
	if exitCode != exitOK {
		t.Fatalf("grinAction exit code: %d", exitCode)
	}

	if buf.String() != string(wantData) {
		t.Errorf("file input mismatch")
	}
}
