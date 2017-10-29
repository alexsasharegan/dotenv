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

// Read an env file at a given path, and return values as a map.
func Read(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	envMap := make(map[string]string)
	// var line, k, v string
	var line string
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			break
		}
		k, v, err := parseln(line)
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
				continue
			}
		}
		if r == hash {
			break
		}
		if r == eq {
			key = buf.String()
			buf.Reset()
			continue
		}
		buf.WriteRune(r)
	}
	value = buf.String()
	return
}
