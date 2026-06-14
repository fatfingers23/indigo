package links

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseJSON(t *testing.T, s string) any {
	t.Helper()
	var v any
	require.NoError(t, json.Unmarshal([]byte(s), &v))
	return v
}

func sortLinks(links []CollectedLink) {
	sort.Slice(links, func(i, j int) bool {
		if links[i].Path != links[j].Path {
			return links[i].Path < links[j].Path
		}
		return links[i].Target.Value < links[j].Target.Value
	})
}

func l(path string, target Link) CollectedLink {
	return CollectedLink{Path: path, Target: target}
}

// ports test_collect_links
func TestCollectLinks(t *testing.T) {
	rec := `{"a": "https://example.com", "b": "not a link"}`
	got := CollectLinks(parseJSON(t, rec))
	assert.Equal(t, []CollectedLink{
		l(".a", Link{Type: LinkURI, Value: "https://example.com"}),
	}, got)
}

// ports test_bsky_feed_post_record_reply
func TestBskyFeedPostRecordReply(t *testing.T) {
	rec := `{
		"$type": "app.bsky.feed.post",
		"createdAt": "2025-01-08T20:52:43.041Z",
		"langs": ["en"],
		"reply": {
			"parent": {
				"cid": "bafyreifk3bwnmulk37ezrarg4ouheqnhgucypynftqafl4limssogvzk6i",
				"uri": "at://did:plc:b3rzzkblqsxhr3dgcueymkqe/app.bsky.feed.post/3lf6yc4drhk2f"
			},
			"root": {
				"cid": "bafyreifk3bwnmulk37ezrarg4ouheqnhgucypynftqafl4limssogvzk6i",
				"uri": "at://did:plc:b3rzzkblqsxhr3dgcueymkqe/app.bsky.feed.post/3lf6yc4drhk2f"
			}
		},
		"text": "Yup!"
	}`
	got := CollectLinks(parseJSON(t, rec))
	sortLinks(got)
	assert.Equal(t, []CollectedLink{
		l(".reply.parent.uri", Link{Type: LinkAtURI, Value: "at://did:plc:b3rzzkblqsxhr3dgcueymkqe/app.bsky.feed.post/3lf6yc4drhk2f"}),
		l(".reply.root.uri", Link{Type: LinkAtURI, Value: "at://did:plc:b3rzzkblqsxhr3dgcueymkqe/app.bsky.feed.post/3lf6yc4drhk2f"}),
	}, got)
}

// ports test_bsky_feed_post_record_embed
func TestBskyFeedPostRecordEmbed(t *testing.T) {
	rec := `{
		"$type": "app.bsky.feed.post",
		"createdAt": "2025-01-08T20:52:39.539Z",
		"embed": {
			"$type": "app.bsky.embed.external",
			"external": {
				"description": "YouTube video by More Perfect Union",
				"thumb": {
					"$type": "blob",
					"ref": {"$link": "bafkreifxuvkbqksq5usi4cryex37o4absjexuouvgenlb62ojsx443b2tm"},
					"mimeType": "image/jpeg",
					"size": 477460
				},
				"title": "Corporations & Wealthy Elites Are Coopting Our Government. Who Can Stop Them?",
				"uri": "https://youtu.be/oKXm4szEP1Q?si=_0n_uPu4qNKokMnq"
			}
		},
		"facets": [
			{
				"features": [
					{
						"$type": "app.bsky.richtext.facet#link",
						"uri": "https://youtu.be/oKXm4szEP1Q?si=_0n_uPu4qNKokMnq"
					}
				],
				"index": {"byteEnd": 24, "byteStart": 0}
			}
		],
		"langs": ["en"],
		"text": "youtu.be/oKXm4szEP1Q?..."
	}`
	got := CollectLinks(parseJSON(t, rec))
	sortLinks(got)
	assert.Equal(t, []CollectedLink{
		l(".embed.external.uri", Link{Type: LinkURI, Value: "https://youtu.be/oKXm4szEP1Q?si=_0n_uPu4qNKokMnq"}),
		l(".facets[].features[app.bsky.richtext.facet#link].uri", Link{Type: LinkURI, Value: "https://youtu.be/oKXm4szEP1Q?si=_0n_uPu4qNKokMnq"}),
	}, got)
}
