package links

import "fmt"

// WalkRecord recursively traverses a decoded JSON value, collecting any link-shaped
// strings together with the JSON path they were found at. Ports walk_record.
//
// v is expected to be the output of json.Unmarshal into `any`: map[string]any for
// objects, []any for arrays, string for strings.
func WalkRecord(path string, v any, found *[]CollectedLink) {
	switch val := v.(type) {
	case map[string]any:
		for key, child := range val {
			WalkRecord(fmt.Sprintf("%s.%s", path, key), child, found)
		}
	case []any:
		for _, child := range val {
			childPath := path + "[]"
			if o, ok := child.(map[string]any); ok {
				if t, ok := o["$type"].(string); ok {
					childPath = fmt.Sprintf("%s[%s]", path, t)
				}
			}
			WalkRecord(childPath, child, found)
		}
	case string:
		if link, ok := ParseAnyLink(val); ok {
			*found = append(*found, CollectedLink{Path: path, Target: link})
		}
	}
}

// CollectLinks returns all links found in a decoded JSON value. Ports collect_links.
func CollectLinks(v any) []CollectedLink {
	var found []CollectedLink
	WalkRecord("", v, &found)
	return found
}
