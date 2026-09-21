package project

import (
	"cyberstrike-ai/internal/database"
)

// ApplyFactOutgoingLinks replaces a fact's outgoing edges (no change when links is nil).
func ApplyFactOutgoingLinks(db *database.DB, projectID, sourceFactKey, sourceConversationID string, links []database.ProjectFactEdgeInput) error {
	if links == nil {
		return nil
	}
	return db.ReplaceOutgoingProjectFactEdges(projectID, sourceFactKey, sourceConversationID, links)
}

// ResolveFactLinkInputs combines the links array and links_text input (the array takes precedence).
func ResolveFactLinkInputs(links []database.ProjectFactEdgeFromInput, linksText string) ([]database.ProjectFactEdgeFromInput, error) {
	if len(links) > 0 {
		return links, nil
	}
	return ParseFactLinksText(linksText)
}

// ApplyFactIncomingLinks replaces a fact's incoming edges (no change when links is nil).
func ApplyFactIncomingLinks(db *database.DB, projectID, targetFactKey string, links []database.ProjectFactEdgeFromInput) error {
	if links == nil {
		return nil
	}
	return db.ReplaceIncomingProjectFactEdges(projectID, targetFactKey, links)
}

// PersistFactIncomingLinks writes incoming edges and optionally synchronizes the current fact body's Relationships section.
func PersistFactIncomingLinks(db *database.DB, projectID, targetFactKey string, links []database.ProjectFactEdgeFromInput, syncBody bool) error {
	if links == nil {
		return nil
	}
	if err := ApplyFactIncomingLinks(db, projectID, targetFactKey, links); err != nil {
		return err
	}
	if !syncBody {
		return nil
	}
	f, err := db.GetProjectFactByKey(projectID, targetFactKey)
	if err != nil {
		return nil
	}
	in, err := db.ListIncomingProjectFactEdges(projectID, targetFactKey)
	if err != nil {
		return err
	}
	f.Body = SyncBodyLinksSection(f.Body, in)
	_, err = db.UpsertProjectFact(f)
	return err
}

// PersistFactLinksFromParsed writes parsed links (no change when parsed is nil).
func PersistFactLinksFromParsed(db *database.DB, projectID, factKey, sourceConversationID string, parsed *ParsedFactLinks, syncBody bool) error {
	if parsed == nil || parsed.Incoming == nil {
		return nil
	}
	return PersistFactIncomingLinks(db, projectID, factKey, parsed.Incoming, syncBody)
}

// PersistFactOutgoingLinks writes outgoing edges (for low-level APIs such as graph connections; use PersistFactIncomingLinks for body synchronization).
func PersistFactOutgoingLinks(db *database.DB, projectID, sourceFactKey, sourceConversationID string, links []database.ProjectFactEdgeInput, syncBody bool) error {
	if links == nil {
		return nil
	}
	return ApplyFactOutgoingLinks(db, projectID, sourceFactKey, sourceConversationID, links)
}

// LinkCountMap contains incoming/outgoing edge counts for each fact in a project.
type LinkCountMap map[string]LinkCounts

// LinkCounts contains the incoming/outgoing edge counts for one fact.
type LinkCounts struct {
	Outgoing int `json:"outgoing"`
	Incoming int `json:"incoming"`
}

// LoadProjectFactLinkCounts loads edge counts in batches.
func LoadProjectFactLinkCounts(db *database.DB, projectID string) (LinkCountMap, error) {
	edges, err := db.ListProjectFactEdgesByProject(projectID)
	if err != nil {
		return nil, err
	}
	m := LinkCountMap{}
	for _, e := range edges {
		c := m[e.SourceFactKey]
		c.Outgoing++
		m[e.SourceFactKey] = c
		c = m[e.TargetFactKey]
		c.Incoming++
		m[e.TargetFactKey] = c
	}
	return m, nil
}
