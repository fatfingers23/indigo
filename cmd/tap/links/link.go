// Package links ports constellation's (microcosm-rs) link-extraction logic to Go.
//
// It detects backlinks (AT-URIs, DIDs, and plain URIs) inside atproto records so the
// tap firehose/resyncer can forward only records that constellation would actually
// index. The behavior mirrors microcosm-rs:
//   - links/src/lib.rs       (Link, parse_any_link, parse_uri)
//   - links/src/at_uri.rs    (parse_at_uri, at_uri_collection)
//   - links/src/did.rs       (parse_did)
//   - links/src/record.rs    (walk_record, collect_links)
//   - constellation/src/lib.rs (RecordId, ActionableEvent, get_actionable)
package links

import (
	"net/url"
	"strings"
)

// LinkType identifies which kind of link a Link holds.
type LinkType int

const (
	LinkAtURI LinkType = iota
	LinkURI
	LinkDid
)

// Link is a parsed backlink target. Mirrors the Rust `Link` enum.
type Link struct {
	Type  LinkType
	Value string
}

// Name returns the metric-friendly name for the link kind, matching Rust's Link::name().
func (l Link) Name() string {
	switch l.Type {
	case LinkAtURI:
		return "at-uri"
	case LinkURI:
		return "uri"
	case LinkDid:
		return "did"
	default:
		return "unknown"
	}
}

// CollectedLink pairs a link with the JSON path it was found at. Mirrors `CollectedLink`.
type CollectedLink struct {
	Path   string
	Target Link
}

// ParseAnyLink tries to interpret s as an AT-URI, then a DID, then a plain URI, in that
// order (matching Rust's parse_any_link). Returns ok=false if s is not a link.
func ParseAnyLink(s string) (Link, bool) {
	if v, ok := parseAtURI(s); ok {
		return Link{Type: LinkAtURI, Value: v}, true
	}
	if v, ok := parseDID(s); ok {
		return Link{Type: LinkDid, Value: v}, true
	}
	if v, ok := parseURI(s); ok {
		return Link{Type: LinkURI, Value: v}, true
	}
	return Link{}, false
}

// parseAtURI ports parse_at_uri. The Rust implementation also performs percent-decoding
// and full RFC-3986 dot-segment path normalization via fluent_uri; that is intentionally
// simplified here (trailing-slash trimming only). Exact normalization is unnecessary
// because tap forwards the full original record and constellation re-parses links itself
// — this function only decides whether a string is link-shaped.
func parseAtURI(s string) (string, bool) {
	if !isASCII(s) {
		return "", false
	}
	if len(s) > 8*1024 {
		return "", false
	}

	const prefix = "at://"
	if len(s) < len(prefix) || !strings.EqualFold(s[:len(prefix)], prefix) {
		return "", false
	}
	rest := s[len(prefix):]

	// Work backwards: fragment, query, path -> authority.
	base := rest
	var fragment, query, path *string
	if i := strings.IndexByte(base, '#'); i >= 0 {
		f := base[i+1:]
		fragment = &f
		base = base[:i]
	}
	if i := strings.IndexByte(base, '?'); i >= 0 {
		q := base[i+1:]
		query = &q
		base = base[:i]
	}
	if i := strings.IndexByte(base, '/'); i >= 0 {
		p := base[i+1:]
		path = &p
		base = base[:i]
	}

	authority := base
	if authority == "" {
		return "", false
	}
	// Normalization: handles are lowercased; DIDs are left as-is.
	if !strings.HasPrefix(authority, "did:") {
		authority = strings.ToLower(authority)
	}

	var out strings.Builder
	out.WriteString("at://")
	out.WriteString(authority)
	if path != nil {
		p := strings.TrimRight(*path, "/")
		if p != "" {
			out.WriteByte('/')
			out.WriteString(p)
		}
	}
	if query != nil {
		out.WriteByte('?')
		out.WriteString(*query)
	}
	if fragment != nil {
		out.WriteByte('#')
		out.WriteString(*fragment)
	}
	return out.String(), true
}

// parseDID ports parse_did.
func parseDID(s string) (string, bool) {
	if len(s) > 2048 {
		return "", false
	}
	for _, c := range s {
		switch {
		case c >= 'A' && c <= 'Z',
			c >= 'a' && c <= 'z',
			c >= '0' && c <= '9',
			c == '.' || c == '_' || c == ':' || c == '%' || c == '-':
		default:
			return "", false
		}
	}

	unprefixed, ok := strings.CutPrefix(s, "did:")
	if !ok {
		return "", false
	}
	method, identifier, ok := strings.Cut(unprefixed, ":")
	if !ok {
		return "", false
	}
	for _, c := range method {
		if c < 'a' || c > 'z' {
			return "", false
		}
	}
	if identifier == "" || strings.HasSuffix(identifier, ":") {
		return "", false
	}
	return s, true
}

// parseURI validates s as an absolute RFC-3986 URI (a non-empty scheme is required, so
// bare strings, CIDs, dates and handles are rejected) and returns it unchanged. Rust's
// parse_uri normalizes via fluent_uri; see the note on parseAtURI for why exact
// normalization is unnecessary here.
func parseURI(s string) (string, bool) {
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" {
		return "", false
	}
	return s, true
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}
