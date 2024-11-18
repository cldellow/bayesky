package source

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type FileSource struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewFileSource(filename string) (*FileSource, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	return &FileSource{
		file:    file,
		scanner: scanner,
	}, nil
}

func (fs *FileSource) Next() (string, error) {
	if fs.scanner.Scan() {
		// Read the next line and return it as a string
		line := fs.scanner.Text()

		// Optionally, you can unmarshal the line into a struct if you want
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		// You can process or return the line here
		return line, nil
	}

	// If the scanner reaches the end of the file
	if err := fs.scanner.Err(); err != nil {
		return "", err
	}

	// No more lines to read
	return "", nil
}

func (fs *FileSource) Close() error {
	return fs.file.Close()
}
