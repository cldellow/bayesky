package events

import (
	"io/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Load(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, err
}

func TestParsePostSimple(t *testing.T) {
	postBytes, err := Load("testdata/simple-japanese.json")
	assert.Nil(t, err)

	post, err := ParsePost(postBytes)
	assert.Nil(t, err)

	assert.Equal(t, "did:plc:bgmth4w3sycbz7vvh5s47w7v", post.Did)
	assert.Equal(t, uint64(1731748069639244), post.Time_us)
	assert.Equal(t, "3lb2k4bihok2f", post.Rkey)
	assert.Equal(t, "米花町でいくら死人が出ても気にならないが池袋で殺人は気になるなあw", post.Text)
}

func TestParsePostReply(t *testing.T) {
	postBytes, err := Load("testdata/reply.json")
	assert.Nil(t, err)

	post, err := ParsePost(postBytes)
	assert.Nil(t, err)

	assert.Equal(t, "did:plc:m2uytfba5tsqiu36trjzuigs", post.Did)
	assert.Equal(t, uint64(1731748069623875), post.Time_us)
	assert.Equal(t, "3lb2k4cdgxk22", post.Rkey)
	assert.Equal(t, "it really does 🫶🏻 it's so ugly and what tf is the meaning behind gatekeeping profile pics???", post.Text)
	assert.Equal(t, "at://did:plc:3klfhj2bre3ahiorkyboeq6b/app.bsky.feed.post/3lazlg3srh22a", post.Parent_uri)
	assert.Equal(t, "at://did:plc:m2uytfba5tsqiu36trjzuigs/app.bsky.feed.post/3lazgx374d22n", post.Root_uri)
}

func TestParsePostQuoteNoMedia(t *testing.T) {
	postBytes, err := Load("testdata/quote-no-media.json")
	assert.Nil(t, err)

	post, err := ParsePost(postBytes)
	assert.Nil(t, err)

	assert.Equal(t, "did:plc:jzjln2ihjs67zvewxzqbiywz", post.Did)
	assert.Equal(t, uint64(1731749262362035), post.Time_us)
	assert.Equal(t, "3lb2l7tlcys25", post.Rkey)
	assert.Equal(t, "Yes til it get hot.. then get off me", post.Text)
	assert.Equal(t, "at://did:plc:y2sbxgnabyivers76io3twt5/app.bsky.feed.post/3lazesxaunl2g", post.Quote_uri)
}

func TestParsePostQuoteImages(t *testing.T) {
	postBytes, err := Load("testdata/quote-with-image.json")
	assert.Nil(t, err)

	post, err := ParsePost(postBytes)
	assert.Nil(t, err)

	assert.Equal(t, "did:plc:frwulzyxodbxnhsfb7d4ijes", post.Did)
	assert.Equal(t, uint64(1731748083356312), post.Time_us)
	assert.Equal(t, "3lb2k4ltbmk2b", post.Rkey)
	assert.Equal(t, "Post a picture of yourself to show the newbies what they are dealing with", post.Text)
	assert.Equal(t, "at://did:plc:rkuir54pi47jak3r6pokk3xi/app.bsky.feed.post/3lb2imirpc22x", post.Quote_uri)

	// TODO: extract images and add tests
}

func TestParsePostImages(t *testing.T) {
	postBytes, err := Load("testdata/images.json")
	assert.Nil(t, err)

	post, err := ParsePost(postBytes)
	assert.Nil(t, err)

	assert.Equal(t, "did:plc:di3tph55e3rrvuuijfbe7zow", post.Did)
	assert.Equal(t, uint64(1731748069629843), post.Time_us)
	assert.Equal(t, "3lb2k45g6qc2s", post.Rkey)
	assert.Equal(t, "10 years old next year :(", post.Text)

	// TODO: extract images and add tests
}
