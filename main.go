package main

import (
	"bayesky/source"
	//	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Try to model the parts of Bluesky posts that we care about,
// while being gloriously ignorant of a lot of things.
type Post struct {
	did     string // did
	time_us uint64 // time_us
	rkey    string // commit.rkey
	// TODO: should we bother with this timestamp? I believe
	// a user can put whatever they want in here, and thereby
	// cause shenanigans.
	// `time_us` might be better for our use case in practice.
	// createdAt string

	text string // commit.record.reply.text

	// Is this a reply?
	parent_uri string // commit.record.reply.parent.uri or empty
	root_uri   string // commit.record.reply.root.uri or empty

	// Is this a quote?
	quote_uri string // commit.record.embed.record.uri
}

func ParsePost(line string) (Post, error) {
	var post Post
	parsed := make(map[string]interface{})
	// TODO: can we defer turning the bytes into a string
	// in Source for an optimization here?
	d := json.NewDecoder(strings.NewReader(line))
	d.UseNumber()
	err := d.Decode(&parsed)
	if err != nil {
		return Post{}, err
	}

	var kind = parsed["kind"].(string)
	if kind != "commit" {
		return post, errors.New("unexpected kind")
	}

	// TODO: we should confirm these match.
	// Currently code elsewhere filters non-posts, but that
	// may not always be the case.
	// commit.operation == create
	// commit.collection == app.bsky.feed.post

	var did = parsed["did"].(string)
	var time_usNumber = parsed["time_us"].(json.Number)
	time_us, err := strconv.ParseUint(string(time_usNumber), 10, 64)

	if err != nil {
		return Post{}, err
	}
	post.did = did
	post.time_us = time_us

	var commit = parsed["commit"].(map[string]interface{})
	var rkey = commit["rkey"].(string)
	post.rkey = rkey

	var record = commit["record"].(map[string]interface{})
	var text = record["text"].(string)
	post.text = text

	if record["reply"] != nil {
		var reply = record["reply"].(map[string]interface{})
		// CONSIDER: this is something like
		// at://did:plc:dwt2ntmuye3zb3w3ie3b5zgu/app.bsky.feed.post/3lb2kyw2q222t
		// but might be better modelled as its component DID and post ID.
		post.parent_uri = reply["parent"].(map[string]interface{})["uri"].(string)
		post.root_uri = reply["root"].(map[string]interface{})["uri"].(string)
	}

	if record["embed"] != nil {
		var embed = record["embed"].(map[string]interface{})
		var type_ = embed["$type"].(string)

		// eg https://gist.github.com/cldellow/d6f5e01a86aa290745e5995c42fd4c7e
		if type_ == "app.bsky.embed.record" {
			var embedRecord = embed["record"].(map[string]interface{})
			if embedRecord["uri"] != nil {
				post.quote_uri = embedRecord["uri"].(string)
			}
		}

		// eg https://gist.github.com/cldellow/f86506d6ec0065a3ea5deb2732f0c0a0
		if type_ == "app.bsky.embed.recordWithMedia" {
			var embedRecord = embed["record"].(map[string]interface{})
			if embedRecord["uri"] != nil {
				post.quote_uri = embedRecord["uri"].(string)
			}
		}

	}

	return post, nil
}

func main() {
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
		if line == "" {
			// End of file
			break
		}

		// Hacky: we only care about English, apply a rough filter
		// very early on.
		//
		// NOTE: This has the side effect of filtering out non-posts,
		//       so we probably want to loosen that up eventually.
		if !strings.Contains(line, `"langs":["en"]`) {
			continue
		}

		post, err := ParsePost(line)
		if err != nil {
			fmt.Println("Error parsing post:", err)
			return
		}

		fmt.Println(post)
		fmt.Printf("https://bsky.app/profile/%s/post/%s\n", post.did, post.rkey)

		//		fmt.Println(line)
		//		fmt.Println(parsed)
	}
}
