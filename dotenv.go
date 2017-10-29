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
)

var errEmptyln = errors.New("empty line")
var errCommentln = errors.New("comment line")

var varRE = regexp.MustCompile("\\${\\w+}")

// ReadFile reads an env file at a given path, and return values as a map.
func ReadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer f.Close()
	return Read(f)
}

// Read parses the given reader's contents and return values as a map.
func Read(rd io.Reader) (map[string]string, error) {
	r := bufio.NewReader(rd)
	envMap := make(map[string]string)
	var (
		line, k, v string
		err        error
	)
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			break
		}
		if varRE.MatchString(line) {
			line = varRE.ReplaceAllStringFunc(line, func(s string) string {
				return envMap[strings.Trim(s, "${}")]
			})
		}
		k, v, err = parseln(line)
		if err != nil {
			continue
		}
		envMap[k] = v
	}

	if err != io.EOF {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return envMap, nil
}

func parseln(ln string) (key, value string, err error) {
	ln = strings.TrimSpace(ln)
	if strings.HasPrefix(ln, "#") {
		err = errEmptyln
		return
	}
	if ln == "" {
		err = errCommentln
		return
	}
	var (
		buf          bytes.Buffer
		quoteType    rune
		insideQuotes bool
	)

	for _, r := range ln {
		if r == '\'' || r == '"' {
			if insideQuotes && r == quoteType {
				insideQuotes = false
				continue
			}
			if !insideQuotes {
				quoteType = r
				insideQuotes = true
				continue
			}
		}
		if !insideQuotes {
			if r == '#' {
				break
			}
			if r == '=' {
				key = buf.String()
				buf.Reset()
				continue
			}
		}

		buf.WriteRune(r)
	}
	value = buf.String()
	return
}

// LoadFile parses the environment config at the given path
// and loads it into the os environment.
func LoadFile(path string, overload bool) error {
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
