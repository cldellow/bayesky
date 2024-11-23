package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

// Try to model the parts of Bluesky posts that we care about,
// while being gloriously ignorant of a lot of things.
type Blob struct {
	MimeType string
	Size     uint64
	Link     string
}

type Image struct {
	Alt      string
	Height   uint64
	Width    uint64
	Blob     Blob 
}

type Video struct {
	Height   uint64
	Width    uint64
	Blob     Blob 
}

type ExternalEmbed struct {
	Description string
	Thumb Blob
	Title string
	Uri string
}

type Post struct {
	Did     string // did
	Time_us uint64 // time_us
	Rkey    string // commit.rkey
	// TODO: should we bother with this timestamp? I believe
	// a user can put whatever they want in here, and thereby
	// cause shenanigans.
	// `time_us` might be better for our use case in practice.
	// createdAt string

	Text string // commit.record.reply.text

	// Is this a reply?
	Parent_uri string // commit.record.reply.parent.uri or empty
	Root_uri   string // commit.record.reply.root.uri or empty

	// Is this a quote?
	Quote_uri string // commit.record.embed.record.uri

	// Does this embed a link?
	ExternalEmbed ExternalEmbed

	Images []Image
	Video Video
}

func extractBlob(subimage map[string]interface{}) (Blob, error) {
	var ref = subimage["ref"].(map[string]interface{})
	var blob Blob

	blob.Link = ref["$link"].(string)
	blob.MimeType = subimage["mimeType"].(string)

	var sizeNumber = subimage["size"].(json.Number)
	size, err := strconv.ParseUint(string(sizeNumber), 10, 64)
	if err != nil {
		return Blob{}, err
	}
	blob.Size = size
	return blob, nil
}

func extractVideo(video map[string]interface{}) (Video, error) {
	var rv Video

	var aspectRatio = video["aspectRatio"].(map[string]interface{})
	var widthNumber = aspectRatio["width"].(json.Number)
	width, err := strconv.ParseUint(string(widthNumber), 10, 64)
	if err != nil {
		return Video{}, err
	}
	var heightNumber = aspectRatio["height"].(json.Number)
	height, err := strconv.ParseUint(string(heightNumber), 10, 64)
	if err != nil {
		return Video{}, err
	}
	rv.Width = width
	rv.Height = height

	blob, err := extractBlob(video["video"].(map[string]interface{}))
	if err != nil {
		return Video{}, err
	}
	rv.Blob = blob
	return rv, nil
}

func extractImage(image map[string]interface{}) (Image, error) {
	var img Image
	img.Alt = image["alt"].(string)

	if image["aspectRatio"] != nil {
		var aspectRatio = image["aspectRatio"].(map[string]interface{})
		var widthNumber = aspectRatio["width"].(json.Number)
		width, err := strconv.ParseUint(string(widthNumber), 10, 64)
		if err != nil {
			return Image{}, err
		}
		var heightNumber = aspectRatio["height"].(json.Number)
		height, err := strconv.ParseUint(string(heightNumber), 10, 64)
		if err != nil {
			return Image{}, err
		}
		img.Width = width
		img.Height = height
	}

	blob, err := extractBlob(image["image"].(map[string]interface{}))
	if err != nil {
		return Image{}, err
	}
	img.Blob = blob

	return img, nil
}

func extractImages(images []interface{}) ([]Image, error) {
	var rv []Image

	// This is gross, but I don't know how to do the type cast
	// in a syntactically cleaner way and my gpt-4o credits
	// have run out for today. :)
	for _, interfaceImage := range images {
		var image = interfaceImage.(map[string]interface{})
		var img, err = extractImage(image)
		if err != nil {
			return nil, err
		}
		rv = append(rv, img)
	}
	return rv, nil
}

func ParsePost(line []byte) (Post, error) {
	var post Post
	parsed := make(map[string]interface{})
	d := json.NewDecoder(bytes.NewReader(line))
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
	post.Did = did
	post.Time_us = time_us

	var commit = parsed["commit"].(map[string]interface{})
	var rkey = commit["rkey"].(string)
	post.Rkey = rkey

	var record = commit["record"].(map[string]interface{})
	var text = record["text"].(string)
	post.Text = text

	if record["reply"] != nil {
		var reply = record["reply"].(map[string]interface{})
		// CONSIDER: this is something like
		// at://did:plc:dwt2ntmuye3zb3w3ie3b5zgu/app.bsky.feed.post/3lb2kyw2q222t
		// but might be better modelled as its component DID and post ID.
		post.Parent_uri = reply["parent"].(map[string]interface{})["uri"].(string)
		post.Root_uri = reply["root"].(map[string]interface{})["uri"].(string)
	}

	if record["embed"] != nil {
		var embed = record["embed"].(map[string]interface{})
		var type_ = embed["$type"].(string)

		if type_ == "app.bsky.embed.images" {
			var rv, err = extractImages(embed["images"].([]interface{}))
			if err != nil {
				return Post{}, err
			}
			post.Images = rv
		}

		// eg https://gist.github.com/cldellow/d6f5e01a86aa290745e5995c42fd4c7e
		if type_ == "app.bsky.embed.record" {
			var embedRecord = embed["record"].(map[string]interface{})
			if embedRecord["uri"] != nil {
				post.Quote_uri = embedRecord["uri"].(string)
			}
		}

		// eg https://gist.github.com/cldellow/f86506d6ec0065a3ea5deb2732f0c0a0
		if type_ == "app.bsky.embed.recordWithMedia" {
			var embedRecord = embed["record"].(map[string]interface{})

			var embedRecordRecord = embedRecord["record"].(map[string]interface{})
			if embedRecordRecord["uri"] != nil {
				post.Quote_uri = embedRecordRecord["uri"].(string)
			}

			var media = embed["media"].(map[string]interface{})
			var mediaType = media["$type"].(string)

			if mediaType == "app.bsky.embed.images" {
				var rv, err = extractImages(media["images"].([]interface{}))
				if err != nil {
					return Post{}, nil
				}

				post.Images = rv
			} else if mediaType == "app.bsky.embed.external" {
				var external = media["external"].(map[string]interface{})

				if external["thumb"] != nil {
					var thumb = external["thumb"].(map[string]interface{})
					var img, err = extractBlob(thumb)
					if err != nil {
						return Post{}, nil
					}
					post.ExternalEmbed.Thumb = img
				}
				post.ExternalEmbed.Title = external["title"].(string)
				post.ExternalEmbed.Description = external["description"].(string)
				post.ExternalEmbed.Uri = external["uri"].(string)
			} else if mediaType == "app.bsky.embed.video" {
				video, err := extractVideo(media)
				if err != nil {
					return Post{}, nil
				}
				post.Video = video
			}
		}
	}

	return post, nil
}
