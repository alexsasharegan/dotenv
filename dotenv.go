/*
Package dotenv is an implementation of the Ruby dotenv library.
The purpose of the library is to load variables from a file into the environment.
*/
package dotenv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const (
	posKey int = iota
	posVal
)

// Errors occurred during parsing.
var (
	ErrInvalidln = errors.New("invalid line")
	ErrEmptyln   = errors.New("empty line")
	ErrCommentln = errors.New("comment line")
)

var varRE = regexp.MustCompile(`\${\w+}`)

// ReadFile reads an env file at a given path, and return values as a map.
func ReadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

// Read parses the given reader's contents and return values as a map.
func Read(rd io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(rd)
	envMap := make(map[string]string)
	var (
		line, k, v string
		err        error
	)
	for scanner.Scan() {
		line = scanner.Text()
		if varRE.MatchString(line) {
			line = varRE.ReplaceAllStringFunc(line, func(s string) string {
				return envMap[strings.Trim(s, "${}")]
			})
		}
		k, v, err = ParseString(line)
		if err == ErrInvalidln {
			return nil, fmt.Errorf("could not parse file: %v", err)
		}
		envMap[k] = v
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return envMap, nil
}

// ParseString parses a given string into a key, value pair.
// Returns the key, value, and an error.
func ParseString(s string) (key, value string, err error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "#") {
		err = ErrCommentln
		return
	}
	if s == "" {
		err = ErrEmptyln
		return
	}
	if !strings.Contains(s, "=") {
		err = ErrInvalidln
		return
	}

	var (
		buf            bytes.Buffer
		quoteType      rune
		hasEq, inQuo   bool
		mapPos, escPos int = posKey, -1
	)

	for i, r := range s {
		if inQuo {
			if i == escPos {
				switch r {
				case 'n':
					buf.WriteString("\n")
				case 'r':
					buf.WriteString("\r")
				default:
					buf.WriteRune(r)
				}
				continue
			}
			// Mark escapes
			if r == '\\' {
				escPos = i + 1
				continue
			}
		}
		// Check for quote delimiters
		if r == '\'' || r == '"' {
			// Look for closing delimiter
			if inQuo && r == quoteType {
				inQuo = false
				// Don't parse beyond a value's terminating quote
				if mapPos == posVal {
					break
				}
				continue
			}
			// Mark quote as delimiter if at start of key/val
			if !inQuo && buf.Len() == 0 {
				quoteType = r
				inQuo = true
				continue
			}
		}
		// If we're inside quotes and not being escaped,
		// ignore certain tokens.
		if !inQuo {
			if unicode.IsSpace(r) {
				continue
			}
			if r == '#' {
				break
			}
			if mapPos == posKey && r == '=' {
				key = buf.String()
				buf.Reset()
				hasEq = true
				mapPos++
				continue
			}
		}
		buf.WriteRune(r)
	}
	if !hasEq {
		err = ErrInvalidln
		return
	}
	value = buf.String()
	// Watch out for a values that include a quote
	if inQuo {
		value = string(quoteType) + value
	}
	return
}

// Load will load a variadic number of environment config files.
// Will not overwrite currently set env vars.
func Load(paths ...string) (err error) {
	if len(paths) == 0 {
		paths = append(paths, ".env")
	}
	for _, path := range paths {
		err = loadFile(path, false)
		if err != nil {
			return
		}
	}
	return
}

// Overload will load a variadic number of environment config files.
// Overwrites currently set env vars.
func Overload(paths ...string) (err error) {
	if len(paths) == 0 {
		paths = append(paths, ".env")
	}
	for _, path := range paths {
		err = loadFile(path, true)
		if err != nil {
			return
		}
	}
	return
}

// loadFile parses the environment config at the given path
// and loads it into the os environment.
func loadFile(path string, overload bool) error {
	env, err := ReadFile(path)
	if err != nil {
		return err
	}
	loadMap(env, overload)
	return nil
}

func loadMap(envMap map[string]string, overload bool) {
	currentEnv := make(map[string]bool)
	for _, rawEnvLine := range os.Environ() {
		currentEnv[strings.Split(rawEnvLine, "=")[0]] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			os.Setenv(key, value)
		}
	}
}
