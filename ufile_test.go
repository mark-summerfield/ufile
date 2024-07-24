package ufile

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

func Test_Barename(t *testing.T) {
	names := []string{"/home/mark/data.dat", `C:\Users\mark\config.ini`,
		"archive.tar.gz", "master.zip", "README", "weird.named.file.xz"}
	bares := []string{"data", "config", "archive", "master", "README",
		"weird"}
	for i, name := range names {
		bare := Barename(name)
		if bare != bares[i] {
			t.Errorf("expected %q,\ngot %q", bares[i], bare)
		}
	}
}

func Test_LongestCommonPath1(t *testing.T) {
	items := []string{"/home/mark/app/go/ufile",
		"/home/mark/app/py/accelhints", "/home/mark/app/rs"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\home\mark\app` {
			t.Errorf(`expected \home\mark\app got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/home/mark/app" {
			t.Errorf("expected /home/mark/app got %q", prefix)
		}
	}
}

func Test_LongestCommonPath2(t *testing.T) {
	items := []string{"/users/mark/app/go/ufile",
		"/Users/mark/app/py/accelhints", "/home/mark/app/rs"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\` {
			t.Errorf("expected \"\\\" got %q", prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/" {
			t.Errorf("expected \"/\" got %q", prefix)
		}
	}
}

func Test_LongestCommonPath3(t *testing.T) {
	items := []string{"C:\\users\\mark\\app\\go\\ufile",
		"/Users/mark/app/py/accelhints"}
	prefix := LongestCommonPath(items)
	if prefix != "" {
		t.Errorf("expected \"\" got %s", prefix)
	}
}

func Test_LongestCommonPath4(t *testing.T) {
	items := []string{"mark/app/go/ufile", "mark/app/py/accelhints",
		"mark/app/rs", "mark/app/rsc"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `mark\app` {
			t.Errorf(`expected mark\app got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "mark/app" {
			t.Errorf("expected mark/app got %q", prefix)
		}
	}
}

func Test_LongestCommonPath5(t *testing.T) {
	items := []string{"/home/mark/app/go/ufile", "/home/mark/apps/bin",
		"/home/mark/app/py/accelhints", "/home/mark/app/rs"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\home\mark` {
			t.Errorf(`expected \home\mark got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/home/mark" {
			t.Errorf("expected /home/mark got %q", prefix)
		}
	}
}

func Test_LongestCommonPath6(t *testing.T) {
	items := []string{"/home/mark/bin/checksum", "/home/mark/bin/checkkey"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\home\mark\bin` {
			t.Errorf(`expected \home\mark\bin got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/home/mark/bin" {
			t.Errorf("expected /home/mark/bin got %q", prefix)
		}
	}
	items = []string{"/home/mark/bin/checksum", "/home/mark/bip/checkkey"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\home\mark` {
			t.Errorf(`expected \home\mark got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/home/mark" {
			t.Errorf("expected /home/mark got %q", prefix)
		}
	}
	items = []string{"/home/marm/bin/checksum", "/home/mark/bin/checkkey"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\home` {
			t.Errorf(`expected \home got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/home" {
			t.Errorf("expected /home got %q", prefix)
		}
	}
	items = []string{"/home/marm/bin/checksum", "/homer/mark/bin/checkkey"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\` {
			t.Errorf(`expected \ got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/" {
			t.Errorf("expected / got %q", prefix)
		}
	}
	items = []string{"/home/page", "/homeric/poem"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != `\` {
			t.Errorf(`expected \ got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "/" {
			t.Errorf("expected / got %q", prefix)
		}
	}
	items = []string{"home", "homeric"}
	if runtime.GOOS == "windows" {
		for i := range len(items) {
			items[i] = filepath.FromSlash(items[i])
		}
		prefix := LongestCommonPath(items)
		if prefix != "" {
			t.Errorf(`expected \"\" got %q`, prefix)
		}
	} else {
		prefix := LongestCommonPath(items)
		if prefix != "" {
			t.Errorf("expected \"\" got %q", prefix)
		}
	}
}

func Test_read_write_text(t *testing.T) {
	Lines := []string{"This text file has non-ASCII characters:",
		"é and ä, π and α²."}

	filename := filepath.Join(os.TempDir(), "read-write-text-go.txt")
	if err := WriteTextFile(filename, Lines); err != nil {
		t.Error(err)
	}
	defer func() { os.Remove(filename) }()
	newLines, err := ReadTextFile(filename)
	if err != nil {
		t.Error(err)
	}
	if slices.Compare(Lines, newLines) != 0 {
		log.Fatalf("%q split !=\n%q\n", Lines, newLines)
	}
	lines := []string{}
	for line, err := range ReadUtf8Lines(filename) {
		if err != nil { // read error, not io.EOF
			t.Error(err)
		}
		lines = append(lines, line)
	}
	if slices.Compare(Lines, lines) != 0 {
		log.Fatalf("%q ReadUtf8Lines !=\n%q\n", Lines, lines)
	}
}
