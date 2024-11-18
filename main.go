package main

import (
	"bayesky/events"
	"bayesky/source"
	"bytes"
	"fmt"
	// "strings"
)

func main() {
	EnLangsBytes := []byte(`"langs":["en"]`)

	// TODO: make the source configurable at the CLI
	// ...and eventually support Jetstream as a source
	fileSource, err := source.NewFileSource("data.json")
	//fileSource, err := source.NewFileSource("24ish.jsonl")
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
		if line == nil {
			// End of file
			break
		}

		// Hacky: we only care about English, apply a rough filter
		// very early on.
		//
		// NOTE: This has the side effect of filtering out non-posts,
		//       so we probably want to loosen that up eventually.
		if !bytes.Contains(line, EnLangsBytes) {
			continue
		}
		/*
			if !strings.Contains(line, `"langs":["en"]`) {
				continue
			}
		*/

		_, err = events.ParsePost(line)
		if err != nil {
			fmt.Println("Error parsing post:", err)
			return
		}

		// fmt.Println(post)
		// fmt.Printf("https://bsky.app/profile/%s/post/%s\n", post.Did, post.Rkey)

		//		fmt.Println(line)
		//		fmt.Println(parsed)
	}
}
