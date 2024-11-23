package events

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
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

	assert.Equal(t, 1, len(post.Images))
	assert.Equal(t, "Woman with a Luigi hat wearing a vacuum cleaner on her back in spooky lighting.", post.Images[0].Alt)
	assert.Equal(t, uint64(1919), post.Images[0].Width)
	assert.Equal(t, uint64(1080), post.Images[0].Height)
	assert.Equal(t, "bafkreifbcyu7bacru2dllvhvyfntddwlj67yozyycuebsuquvl5zini4da", post.Images[0].Blob.Link)
	assert.Equal(t, "image/jpeg", post.Images[0].Blob.MimeType)
	assert.Equal(t, uint64(444761), post.Images[0].Blob.Size)

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

	assert.Equal(t, 2, len(post.Images))
	assert.Equal(t, "foo", post.Images[0].Alt)
	assert.Equal(t, uint64(960), post.Images[0].Width)
	assert.Equal(t, uint64(960), post.Images[0].Height)
	assert.Equal(t, "bafkreien2xwiyz5m4nftsuac7ymbpnvpwcakr5fmrp45qca6psgvglyzxy", post.Images[0].Blob.Link)
	assert.Equal(t, "image/jpeg", post.Images[0].Blob.MimeType)
	assert.Equal(t, uint64(449738), post.Images[0].Blob.Size)

	assert.Equal(t, "", post.Images[1].Alt)
	assert.Equal(t, uint64(497), post.Images[1].Width)
	assert.Equal(t, uint64(604), post.Images[1].Height)
	assert.Equal(t, "bafkreiaq3uz7vfluz3nejy34wqmxl6fizlunhnmymnszewomemh2jgou2y", post.Images[1].Blob.Link)
	assert.Equal(t, "image/jpeg", post.Images[1].Blob.MimeType)
	assert.Equal(t, uint64(156550), post.Images[1].Blob.Size)
}

func TestParsePostImages2(t *testing.T) {
	postBytes, err := Load("testdata/image-no-aspect.json")
	assert.Nil(t, err)

	post, err := ParsePost(postBytes)
	assert.Nil(t, err)

	assert.Equal(t, "did:plc:v3wrmaefdkqejoz2c2fiaccl", post.Did)
	assert.Equal(t, 1, len(post.Images))
	assert.Equal(t, "IEMBot Image TBD", post.Images[0].Alt)
	assert.Equal(t, uint64(0), post.Images[0].Width)
	assert.Equal(t, uint64(0), post.Images[0].Height)
	assert.Equal(t, "bafkreifevvbvn7lihsmnyv5qtwvgrmrydrqqtpp73lxdhbziaxjzjlysfe", post.Images[0].Blob.Link)
	assert.Equal(t, "image/png", post.Images[0].Blob.MimeType)
	assert.Equal(t, uint64(638076), post.Images[0].Blob.Size)
}


func TestParsePostExternalWithThumb(t *testing.T) {
	postBytes, err := Load("testdata/external-with-thumb.json")
	assert.Nil(t, err)

	rv, err := ParsePost(postBytes)
	assert.Nil(t, err)
	assert.Equal(t, "image/jpeg", rv.ExternalEmbed.Thumb.MimeType)
	assert.Equal(t, uint64(91644), rv.ExternalEmbed.Thumb.Size)
	assert.Equal(t, "bafkreiah6vp4wwac3mukki3cxsxynq5bo2yjub35ej7emlyuuv3coarw6e", rv.ExternalEmbed.Thumb.Link)
	assert.Equal(t, "ALT: a red background with a few pieces of wood in the middle", rv.ExternalEmbed.Description)
	assert.Equal(t, "a red background with a few pieces of wood in the middle", rv.ExternalEmbed.Title)
	assert.Equal(t, "https://media.tenor.com/1TG5kWc6M3gAAAAC/eldenring.gif?hh=280&ww=498", rv.ExternalEmbed.Uri)
}

func TestParsePostFacets(t *testing.T) {
	postBytes, err := Load("testdata/facets.json")
	assert.Nil(t, err)

	_, err = ParsePost(postBytes)
	assert.Nil(t, err)
}

func TestParsePostVideo(t *testing.T) {
	postBytes, err := Load("testdata/video.json")
	assert.Nil(t, err)

	rv, err := ParsePost(postBytes)
	assert.Nil(t, err)

	assert.Equal(t, "at://did:plc:d3qj2eu7mqqo5u5rcgejvuhm/app.bsky.feed.post/3lb2k44wno222", rv.Quote_uri)

	assert.Equal(t, uint64(1080), rv.Video.Width)
	assert.Equal(t, uint64(1920), rv.Video.Height)

	assert.Equal(t, "bafkreieyi3vgtbjwxoxncwwkqxrglbclandonc4blslerpxm5gk42pdiqq", rv.Video.Blob.Link)
	assert.Equal(t, "video/mp4", rv.Video.Blob.MimeType)
	assert.Equal(t, uint64(1875400), rv.Video.Blob.Size)
}
