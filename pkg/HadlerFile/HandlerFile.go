package HadlerFile

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

type FileSource struct {
	Filename string
}

func (f FileSource) Handler(stringPattern *regexp.Regexp) (currentCount uint64, err error) {
	file, err := os.Open(f.Filename)
	if err != nil {
		return 0, fmt.Errorf("error in open file: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			err = fmt.Errorf("error in close file: %w", err)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("error read content in file %s: %w", f.Filename, err)
	}

	count := len(stringPattern.FindAll(content, -1))
	return uint64(count), nil
}
