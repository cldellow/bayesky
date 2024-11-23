# bayesky

Bluesky is designed to be hackable.

It comes stock with a reverse chronological feed and a proprietary "For you"
feed. But developers can [author their own feeds](https://docs.bsky.app/docs/starter-templates/custom-feeds)
that users can subscribe to.

My goal: write a Humans Being Bros feed. The inclusion criteria for this feed
is threads where:

1. OP asks a question
2. others respond
3. OP expresses gratitude

It could be ranked by size of thread, # of likes on question or non-OP responses,
diversity of repliers, etc.

# How to discover content?

In addition to the official Relays, Bluesky publishes a lighter-weight feed that is
consumable via websocket. This is called the [Jetsream](https://github.com/bluesky-social/jetstream?tab=readme-ov-file).

You can watch the stream of new posts via the `app.bsky.feed.post` collection:

```
$ websocat wss://jetstream2.us-east.bsky.network/subscribe\?wantedCollections=app.bsky.feed.post
```

You'll see data like this sample post:

```json
{
  "did": "did:plc:w5l6zvlmyz3r2cl36bfqlq7a",
  "time_us": 1731868440607689,
  "type": "com",
  "kind": "commit",
  "commit": {
    "rev": "3lb627l4oc62h",
    "type": "c",
    "operation": "create",
    "collection": "app.bsky.feed.post",
    "rkey": "3lb627kz72s2r",
    "record": {
      "$type": "app.bsky.feed.post",
      "createdAt": "2024-11-17T18:33:58.271Z",
      "langs": [
        "en"
      ],
      "text": "Test post: testing JetStream."
    },
    "cid": "bafyreigwgz44ovvc4lu2nyklh3meclhjxnipewxaoswcdm4mj3vqled4ee"
  }
}
```

You might also want to track likes, e.g. for ranking. You can subscribe to `app.bsky.feed.like`,
to see things like:

```json
{
  "did": "did:plc:ko26dqkkmj3da6yc3fmo3ate",
  "time_us": 1731870626607952,
  "type": "com",
  "kind": "commit",
  "commit": {
    "rev": "3lb64aov7ii23",
    "type": "c",
    "operation": "create",
    "collection": "app.bsky.feed.like",
    "rkey": "3lb64aov4kq23",
    "record": {
      "$type": "app.bsky.feed.like",
      "createdAt": "2024-11-17T19:10:23.248Z",
      "subject": {
        "cid": "bafyreibn5x7unywvqytekgfg43kwruq4zyqzjnfa4kn7dp2rc7tq2mgvoy",
        "uri": "at://did:plc:65otgq6ubushgm3vk5icuxzw/app.bsky.feed.post/3lb3bqy7ibe2v"
      }
    },
    "cid": "bafyreie2avqg2zhg4dxuibxunlnjxjr4fpzdhvvpbzxquvfteax2qvjrne"
  }
}
```

# Overall approach

The firehose operates on events like "new post" or "liked post", but we want to surface
something higher-level like "threads with this kind of interaction".

The first building block will be a classifier for posts, so we can classify posts like
this:

1. Post is a top-level post that asks a question
2. Post is a non-top-level post that replies to a post in class 1, and is
   by a different author. (Wrinkle: what if the question is itself a multi-post
	 thread?)
3. Post is a non-top-level post that replies to a post in class 2, and is
   by the same author as the thread starter, and expresses gratitude.

I _think_ a naive Bayes classifier might be enough here, especially if we can help
it along by providing some clever feature extraction, e.g. emitting `AUTHOR_IS_THREAD_AUTHOR`
or `AUTHOR_IS_NOT_THREAD_AUTHOR` features.

I know LLMs are the new hotness, but they're expensive to run. A well-tuned naive Bayes
classifier should be able to handle the firehose on a single core without breaking
a sweat.

# Training the classifier

A challenge with naive Bayes is training it. The classic approach is to label a bunch
of samples as positive or negative, then train a model.

Labelling is tedious and sucks.

Maybe there's room here for an LLM to be used: you could express your desired classes
in plain language, and apply an LLM to generate best-effort labels. A human quickly
reviews them and accepts/rejects the labels, and that becomes your training set.

Perhaps Llamafile with a reasonably-sized model could be used here?

# Ops questions

The Bluesky firehose is not _that_ big at present. ~200 posts/second, ~800 likes/second.

This is just a side project, so being a little lossy is fine if it simplifies perf problems.

My overall hope is to do something like this:

- apply a Bayes classifier to the stream of posts. Hopefully we discard 99.9%+ of posts.
- track the IDs of non-discarded posts
- only track likes for non-discarded posts; buffer them in-memory and checkpoint to a SQLite
  DB on some cadence so that we can interrupt/resume Jetstream processing via `cursor`
- retain persisted data for at most 7 days

The Bayes classification can be farmed out amongst threads, but the overall processing needs
to be sequential -- e.g. we have to know we've processed post X before processing any likes for it,
or before processing replies to post X.

# Golang notes

It's been years since I wrote go code. I'm relying on ChatGPT a lot. Useful commands:

```bash
$ go test ./... # run all tests, recursively

$ gofmt -w .    # format all files, recursively
```
