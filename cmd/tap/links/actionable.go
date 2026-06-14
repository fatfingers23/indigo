package links

// RecordId identifies a record. Mirrors constellation's RecordId.
type RecordId struct {
	Did        string
	Collection string
	Rkey       string
}

// ActionType discriminates the variants of ActionableEvent.
type ActionType int

const (
	CreateLinks ActionType = iota
	UpdateLinks
	DeleteRecord
	ActivateAccount
	DeactivateAccount
	DeleteAccount
)

// ActionableEvent is the Go equivalent of constellation's ActionableEvent enum.
//
//   - CreateLinks / UpdateLinks: RecordID + Links populated.
//   - DeleteRecord: RecordID populated.
//   - Activate/Deactivate/DeleteAccount: Did populated.
type ActionableEvent struct {
	Type     ActionType
	RecordID RecordId
	Links    []CollectedLink
	Did      string
}

// RecordLinks collects all links for a record: the rkey itself (if it is a link, e.g. a
// DID-keyed vouch record) at path ".", followed by every link found walking the record
// body. Mirrors the link-collection done inside get_actionable for create/update.
func RecordLinks(rkey string, record map[string]any) []CollectedLink {
	var found []CollectedLink
	if link, ok := ParseAnyLink(rkey); ok {
		found = append(found, CollectedLink{Path: ".", Target: link})
	}
	WalkRecord("", record, &found)
	return found
}

// IsActionable reports whether a record operation should be forwarded to constellation,
// mirroring get_actionable's filtering: a "create" is kept only if it has at least one
// link, while "update" and "delete" are always kept. Anything else is dropped.
func IsActionable(action, rkey string, record map[string]any) bool {
	switch action {
	case "create":
		return len(RecordLinks(rkey, record)) > 0
	case "update", "delete":
		return true
	default:
		return false
	}
}

// GetActionable is a full port of constellation's get_actionable. It takes a decoded
// jetstream event (json.Unmarshal into map[string]any) and returns the actionable event
// plus the time_us cursor, or ok=false if the event is not actionable.
func GetActionable(event map[string]any) (ActionableEvent, uint64, bool) {
	timeUs, ok := asUint64(event["time_us"])
	if !ok {
		return ActionableEvent{}, 0, false
	}

	switch event["kind"] {
	case "commit":
		return getCommitActionable(event, timeUs)
	case "account":
		return getAccountActionable(event, timeUs)
	default:
		return ActionableEvent{}, 0, false
	}
}

func getCommitActionable(event map[string]any, timeUs uint64) (ActionableEvent, uint64, bool) {
	did, ok := event["did"].(string)
	if !ok {
		return ActionableEvent{}, 0, false
	}
	commit, ok := event["commit"].(map[string]any)
	if !ok {
		return ActionableEvent{}, 0, false
	}
	collection, ok := commit["collection"].(string)
	if !ok {
		return ActionableEvent{}, 0, false
	}
	rkey, ok := commit["rkey"].(string)
	if !ok {
		return ActionableEvent{}, 0, false
	}
	op, ok := commit["operation"].(string)
	if !ok {
		return ActionableEvent{}, 0, false
	}

	recordID := RecordId{Did: did, Collection: collection, Rkey: rkey}

	switch op {
	case "create", "update":
		record, _ := commit["record"].(map[string]any)
		if _, present := commit["record"]; !present {
			return ActionableEvent{}, 0, false
		}
		links := RecordLinks(rkey, record)
		if op == "create" {
			if len(links) == 0 {
				return ActionableEvent{}, 0, false
			}
			return ActionableEvent{Type: CreateLinks, RecordID: recordID, Links: links}, timeUs, true
		}
		return ActionableEvent{Type: UpdateLinks, RecordID: recordID, Links: links}, timeUs, true
	case "delete":
		return ActionableEvent{Type: DeleteRecord, RecordID: recordID}, timeUs, true
	default:
		return ActionableEvent{}, 0, false
	}
}

func getAccountActionable(event map[string]any, timeUs uint64) (ActionableEvent, uint64, bool) {
	account, ok := event["account"].(map[string]any)
	if !ok {
		return ActionableEvent{}, 0, false
	}
	did, ok := account["did"].(string)
	if !ok {
		return ActionableEvent{}, 0, false
	}
	active, ok := account["active"].(bool)
	if !ok {
		return ActionableEvent{}, 0, false
	}

	if active {
		if _, hasStatus := account["status"]; !hasStatus {
			return ActionableEvent{Type: ActivateAccount, Did: did}, timeUs, true
		}
		return ActionableEvent{}, 0, false
	}

	status, ok := account["status"].(string)
	if !ok {
		return ActionableEvent{}, 0, false
	}
	switch status {
	case "deactivated":
		return ActionableEvent{Type: DeactivateAccount, Did: did}, timeUs, true
	case "deleted":
		return ActionableEvent{Type: DeleteAccount, Did: did}, timeUs, true
	default:
		return ActionableEvent{}, 0, false
	}
}

// asUint64 extracts a uint64 from a decoded JSON number, accepting both float64 (the
// default for json.Unmarshal) and json.Number.
func asUint64(v any) (uint64, bool) {
	switch n := v.(type) {
	case float64:
		return uint64(n), true
	case interface{ Int64() (int64, error) }: // json.Number
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return uint64(i), true
	default:
		return 0, false
	}
}
