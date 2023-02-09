package filetool

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func OpenFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func ReadLineString(f *os.File) ([]string, error) {
	var data []string
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		data = append(data, line)
	}
	return data, nil
}
