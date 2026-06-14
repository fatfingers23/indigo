package links

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseEvent(t *testing.T, s string) map[string]any {
	t.Helper()
	var v map[string]any
	require.NoError(t, json.Unmarshal([]byte(s), &v))
	return v
}

// ports test_create_like
func TestCreateLike(t *testing.T) {
	rec := `{
		"did":"did:plc:icprmty6ticzracr5urz4uum",
		"time_us":1736448492661668,
		"kind":"commit",
		"commit":{"rev":"3lfddpt5qa62c","operation":"create","collection":"app.bsky.feed.like","rkey":"3lfddpt5djw2c","record":{
			"$type":"app.bsky.feed.like",
			"createdAt":"2025-01-09T18:48:10.412Z",
			"subject":{"cid":"bafyreihazf62qvmusup55ojhkzwbmzee6rxtsug3e6eg33mnjrgthxvozu","uri":"at://did:plc:lphckw3dz4mnh3ogmfpdgt6z/app.bsky.feed.post/3lfdau5f7wk23"}
		},
		"cid":"bafyreidgcs2id7nsbp6co42ind2wcig3riwcvypwan6xdywyfqklovhdjq"}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1736448492661668), cursor)
	assert.Equal(t, ActionableEvent{
		Type: CreateLinks,
		RecordID: RecordId{
			Did:        "did:plc:icprmty6ticzracr5urz4uum",
			Collection: "app.bsky.feed.like",
			Rkey:       "3lfddpt5djw2c",
		},
		Links: []CollectedLink{
			l(".subject.uri", Link{Type: LinkAtURI, Value: "at://did:plc:lphckw3dz4mnh3ogmfpdgt6z/app.bsky.feed.post/3lfdau5f7wk23"}),
		},
	}, action)
}

// ports test_update_profile
func TestUpdateProfile(t *testing.T) {
	rec := `{
		"did":"did:plc:tcmiubbjtkwhmnwmrvr2eqnx",
		"time_us":1736453696817289,"kind":"commit",
		"commit":{
			"rev":"3lfdikw7q772c",
			"operation":"update",
			"collection":"app.bsky.actor.profile",
			"rkey":"self",
			"record":{
				"$type":"app.bsky.actor.profile",
				"avatar":{"$type":"blob","ref":{"$link":"bafkreidcg5jzz3hpdtlc7um7w5masiugdqicc5fltuajqped7fx66hje54"},"mimeType":"image/jpeg","size":295764},
				"banner":{"$type":"blob","ref":{"$link":"bafkreiahaswf2yex2zfn3ynpekhw6mfj7254ra7ly27zjk73czghnz2wni"},"mimeType":"image/jpeg","size":856461},
				"createdAt":"2024-08-30T21:33:06.945Z",
				"description":"Professor, QUB | Belfast via Derry \n\nViews personal | Reposts are not an endorsement\n\nhttps://go.qub.ac.uk/charvey",
				"displayName":"Colin Harvey",
				"pinnedPost":{"cid":"bafyreifyrepqer22xsqqnqulpcxzpu7wcgeuzk6p5c23zxzctaiwmlro7y","uri":"at://did:plc:tcmiubbjtkwhmnwmrvr2eqnx/app.bsky.feed.post/3lf66ri63u22t"}
			},
			"cid":"bafyreiem4j5p7duz67negvqarq3s5h7o45fvytevhrzkkn2p6eqdkcf74m"
		}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1736453696817289), cursor)
	assert.Equal(t, ActionableEvent{
		Type: UpdateLinks,
		RecordID: RecordId{
			Did:        "did:plc:tcmiubbjtkwhmnwmrvr2eqnx",
			Collection: "app.bsky.actor.profile",
			Rkey:       "self",
		},
		Links: []CollectedLink{
			l(".pinnedPost.uri", Link{Type: LinkAtURI, Value: "at://did:plc:tcmiubbjtkwhmnwmrvr2eqnx/app.bsky.feed.post/3lf66ri63u22t"}),
		},
	}, action)
}

// ports test_delete_like
func TestDeleteLike(t *testing.T) {
	rec := `{
		"did":"did:plc:3pa2ss4l2sqzhy6wud4btqsj",
		"time_us":1736448492690783,
		"kind":"commit",
		"commit":{"rev":"3lfddpt7vnx24","operation":"delete","collection":"app.bsky.feed.like","rkey":"3lbiu72lczk2w"}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1736448492690783), cursor)
	assert.Equal(t, ActionableEvent{
		Type: DeleteRecord,
		RecordID: RecordId{
			Did:        "did:plc:3pa2ss4l2sqzhy6wud4btqsj",
			Collection: "app.bsky.feed.like",
			Rkey:       "3lbiu72lczk2w",
		},
	}, action)
}

// ports test_delete_account
func TestDeleteAccount(t *testing.T) {
	rec := `{
		"did":"did:plc:zsgqovouzm2gyksjkqrdodsw",
		"time_us":1736451739215876,
		"kind":"account",
		"account":{"active":false,"did":"did:plc:zsgqovouzm2gyksjkqrdodsw","seq":3040934738,"status":"deleted","time":"2025-01-09T19:42:18.972Z"}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1736451739215876), cursor)
	assert.Equal(t, ActionableEvent{
		Type: DeleteAccount,
		Did:  "did:plc:zsgqovouzm2gyksjkqrdodsw",
	}, action)
}

// ports test_deactivate_account
func TestDeactivateAccount(t *testing.T) {
	rec := `{
		"did":"did:plc:l4jb3hkq7lrblferbywxkiol","time_us":1736451745611273,"kind":"account","account":{"active":false,"did":"did:plc:l4jb3hkq7lrblferbywxkiol","seq":3040939563,"status":"deactivated","time":"2025-01-09T19:42:22.035Z"}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1736451745611273), cursor)
	assert.Equal(t, ActionableEvent{
		Type: DeactivateAccount,
		Did:  "did:plc:l4jb3hkq7lrblferbywxkiol",
	}, action)
}

// ports test_create_vouch_indexes_did_rkey
func TestCreateVouchIndexesDidRkey(t *testing.T) {
	rec := `{
		"did":"did:plc:voucher",
		"time_us":1746460800000000,
		"kind":"commit",
		"commit":{"rev":"3lqrvouchcreate","operation":"create","collection":"sh.tangled.graph.vouch","rkey":"did:plc:vouchedfor","record":{
			"$type":"sh.tangled.graph.vouch",
			"createdAt":"2026-05-05T12:00:00.000Z"
		}}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1746460800000000), cursor)
	assert.Equal(t, ActionableEvent{
		Type: CreateLinks,
		RecordID: RecordId{
			Did:        "did:plc:voucher",
			Collection: "sh.tangled.graph.vouch",
			Rkey:       "did:plc:vouchedfor",
		},
		Links: []CollectedLink{
			l(".", Link{Type: LinkDid, Value: "did:plc:vouchedfor"}),
		},
	}, action)
}

// ports test_update_vouch_indexes_did_rkey
func TestUpdateVouchIndexesDidRkey(t *testing.T) {
	rec := `{
		"did":"did:plc:voucher",
		"time_us":1746460800000001,
		"kind":"commit",
		"commit":{"rev":"3lqrvouchupdate","operation":"update","collection":"sh.tangled.graph.vouch","rkey":"did:plc:vouchedfor","record":{
			"$type":"sh.tangled.graph.vouch",
			"createdAt":"2026-05-05T12:00:00.000Z",
			"reason":"https://atproto.com"
		}}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1746460800000001), cursor)
	assert.Equal(t, ActionableEvent{
		Type: UpdateLinks,
		RecordID: RecordId{
			Did:        "did:plc:voucher",
			Collection: "sh.tangled.graph.vouch",
			Rkey:       "did:plc:vouchedfor",
		},
		Links: []CollectedLink{
			l(".", Link{Type: LinkDid, Value: "did:plc:vouchedfor"}),
			l(".reason", Link{Type: LinkURI, Value: "https://atproto.com"}),
		},
	}, action)
}

// ports test_activate_account
func TestActivateAccount(t *testing.T) {
	rec := `{
		"did":"did:plc:nct6zfb2j4emoj4yjomxwml2","time_us":1736451747292706,"kind":"account","account":{"active":true,"did":"did:plc:nct6zfb2j4emoj4yjomxwml2","seq":3040940775,"time":"2025-01-09T19:42:26.924Z"}
	}`
	action, cursor, ok := GetActionable(parseEvent(t, rec))
	require.True(t, ok)
	assert.Equal(t, uint64(1736451747292706), cursor)
	assert.Equal(t, ActionableEvent{
		Type: ActivateAccount,
		Did:  "did:plc:nct6zfb2j4emoj4yjomxwml2",
	}, action)
}
