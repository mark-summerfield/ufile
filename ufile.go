// Copyright © 2024 Mark Summerfield. All rights reserved.

// This package provides file-related utility functions.
package ufile

import (
	"bufio"
	"bytes"
	_ "embed"
	"io"
	"iter"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mark-summerfield/utext"
)

//go:embed Version.dat
var Version string

const ModeURW = 0o600

// AbsPath returns the filename with its path absolute, or cleaned on error.
// See also [Relativized].
func AbsPath(filename string) string {
	absFilename, err := filepath.Abs(filename)
	if err == nil {
		return absFilename
	}
	return filepath.Clean(filename)
}

// Barename returns the filename without any path and without any suffix.
func Barename(path string) string {
	if i := strings.LastIndexAny(path, `/\`); i > -1 {
		path = path[i+1:]
	}
	for {
		if i := strings.LastIndexByte(path, '.'); i > -1 {
			path = path[:i]
		} else {
			break
		}
	}
	return path
}

// FileExists returns true if the filename exists and is a file.
// See also [PathExists].
func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
		return true
	}
	return false
}

// GetConfigFile given a domain name, say, "domain.com", and an application
// name, say, "myapp", and an extention, say, ".json", returns where the
// corresponding config file is located and true, or where the config file
// should be saved (i.e., if it doesn't exist) and false.
//
// When saving (at least for the first time) you may need to create the
// domain folder:
//
//	dir := filepath.Dir(configFilename)
//	if dir != "." {
//		_ = os.MkdirAll(dir, fs.ModePerm)
//	}
//	// now save to configFilename
func GetConfigFile(domain, appname, ext string) (string, bool) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	filename := appname + ext
	filenames := make([]string, 0, 8)
	var preferred string
	var fallback string
	configDir, err := os.UserConfigDir()
	if err == nil {
		if domain != "" {
			preferred = filepath.Join(configDir, domain, filename)
			filenames = append(filenames, preferred)
		}
		name := filepath.Join(configDir, filename)
		filenames = append(filenames, name)
		if preferred == "" {
			preferred = name
		}
	}
	homeDir, err := os.UserHomeDir()
	if err == nil {
		if domain != "" {
			fallback = filepath.Join(homeDir, "."+domain+"-"+filename)
			filenames = append(filenames, fallback)
		}
		name := filepath.Join(homeDir, "."+filename)
		filenames = append(filenames, name)
		if fallback == "" {
			fallback = name
		}
	}
	if len(filenames) == 0 { // if all else fails try current dir
		if domain != "" {
			filenames = append(filenames, domain+"-"+filename)
		}
		filenames = append(filenames, filename)
		if domain != "" {
			filenames = append(filenames, "."+domain+"-"+filename)
		}
		filenames = append(filenames, "."+filename)
	}
	for _, filename := range filenames {
		if FileExists(filename) {
			return filename, true // found
		}
	}
	if preferred != "" {
		return preferred, false
	}
	if fallback != "" {
		return fallback, false
	}
	if domain != "" {
		return domain + "-" + filename, false
	}
	return filename, false
}

// GetIniFile given a domain name, say, "domain.com", and an application
// name, say, "myapp", returns where the corresponding .ini file is located
// and true, or where the .ini should be saved (i.e., if it doesn't exist)
// and false.
//
// When saving (at least for the first time) you may need to create the
// domain folder:
//
//	dir := filepath.Dir(iniFilename)
//	if dir != "." {
//		_ = os.MkdirAll(dir, fs.ModePerm)
//	}
//	// now save to iniFilename
func GetIniFile(domain, appname string) (string, bool) {
	return GetConfigFile(domain, appname, ".ini")
}

// HomeDir returns the abs path of the home folder, e.g., `/home/mark`.
func HomeDir() string {
	name, err := os.UserHomeDir()
	if err != nil {
		return AbsPath(".")
	}
	return name
}

// IsDir returns true if name is a folder; otherwise returns false.
func IsDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// LongestCommonPath returns the longest common path, i.e., component,
// / or \ separated (which could be "" if there isn't one), and lowercased
// on Windows.
func LongestCommonPath(paths []string) string {
	caseInsensitive := runtime.GOOS == "windows" ||
		runtime.GOOS == "darwin"
	if len(paths) == 0 {
		return ""
	} else if len(paths) == 1 {
		if caseInsensitive {
			return strings.ToLower(paths[0])
		}
		return paths[0]
	}
	if caseInsensitive {
		for i := range len(paths) {
			paths[i] = strings.ToLower(paths[i])
		}
	}
	prefix := utext.LongestCommonPrefix(paths)
	if len(prefix) > 1 {
		i := strings.LastIndexByte(prefix, os.PathSeparator)
		if i == -1 { // no path separator to slice to
			prefix = ""
		} else {
			if i == 0 { // preserve root of / or \
				i = 1
			}
			prefix = prefix[:i]
		}
	}
	return prefix
}

// PathExists returns true if the path/filename exists.
// See also [FileExists].
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadTextFile reads the given file and returns a slices of lines with
// EOL stripped off. See also [ReadUtf8Lines]
func ReadTextFile(filename string) ([]string, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	raw = bytes.ReplaceAll(raw, []byte{'\r'}, []byte{})
	raw = bytes.TrimRight(raw, "\n")
	return strings.Split(string(raw), "\n"), nil
}

// ReadUtf8Lines reads the given file and returns an iterator of (line,
// error) for every line with EOL stripped off. See also [ReadTextFile].
func ReadUtf8Lines(filename string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		file, err := os.Open(filename)
		if err != nil {
			yield("", err) // failed to open file
			return         // we cannot progress from here
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				yield("", err) // read error
				return         // we cannot progress further
			}
			if line == "" && err == io.EOF {
				break // last (i.e., prev.) line ended with \n
			}
			if !yield(strings.TrimRight(line, "\r\n"), nil) {
				return // for loop break or return or panic
			}
			if err == io.EOF {
				break // last line did not end with \n
			}
		}
	}
}

// WriteTextFile writes the given lines to the given filename adding the
// platform-appropriate EOL to each line written.
func WriteTextFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}
	out := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err = out.WriteString(line); err != nil {
			return err
		}
		if _, err = out.WriteString(eol); err != nil {
			return err
		}
	}
	out.Flush()
	return nil
}
