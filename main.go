package main

import (
	"bayesky/source"
	"fmt"
	"strings"
)

func main() {
	// TODO: make the source configurable at the CLI
	// ...and eventually support Jetstream as a source
	fileSource, err := source.NewFileSource("data.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileSource.Close() // Ensure file is closed when done

	// Using the Next() method to read each line (JSON record) from the file
	for {
		line, err := fileSource.Next()
		if err != nil {
			fmt.Println("Error reading line:", err)
			break
		}
		if line == "" {
			// End of file
			break
		}

		// Hacky: we only care about English, apply a rough filter
		// very early on.
		if !strings.Contains(line, `"langs":["en"]`) {
			continue
		}

		fmt.Println(line)
	}
}
