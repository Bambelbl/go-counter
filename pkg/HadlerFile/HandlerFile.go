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
		return 0, fmt.Errorf("error in open file: %s", err.Error())
	}

	defer func() {
		err = file.Close()
		if err != nil {
			err = fmt.Errorf("error in close file: %s", err.Error())
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("error read content in file %s: %s", f.Filename, err.Error())
	}

	count := len(stringPattern.FindAll(content, -1))
	return uint64(count), nil
}
