package file

import (
	"bufio"
	"os"
)

type file struct {
	file     *os.File
	filename string
}

// New creates a new file handler.
func New(filename string) *file {
	return &file{filename: filename}
}

// Open opens the file.
func (f *file) Open() error {
	file, err := os.Open(f.filename)
	if err != nil {
		return err
	}

	f.file = file
	return nil
}

// Close closes the file.
func (f *file) Close() error {
	return f.file.Close()
}

// Read reads the file and returns its lines.
func (f *file) Read() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(f.file)
	scanner.Scan()
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
