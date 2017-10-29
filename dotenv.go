package dotenv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	squo = '\''
	dquo = '"'
	hash = '#'
	eq   = '='
)

var errEmptyln = errors.New("empty line")
var errCommentln = errors.New("comment line")

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
		k, v, err = parseln(line, envMap)
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
	var buf bytes.Buffer
	var quoteType rune
	insideQuotes := false
	for _, r := range ln {
		if r == squo || r == dquo {
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
			if r == hash {
				break
			}
			if r == eq {
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
