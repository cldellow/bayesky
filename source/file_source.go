package source

import (
	"bufio"
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

// func (fs *FileSource) Next() (string, error) {
func (fs *FileSource) Next() ([]byte, error) {
	if fs.scanner.Scan() {
		//line := fs.scanner.Text()
		line := fs.scanner.Bytes()

		return line, nil
	}

	// If the scanner reaches the end of the file
	if err := fs.scanner.Err(); err != nil {
		return nil, err
	}

	// No more lines to read
	return nil, nil
}

func (fs *FileSource) Close() error {
	return fs.file.Close()
}
