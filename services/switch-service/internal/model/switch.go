package model

type (
	// Switch is the Arango model for proto.Switch
	Switch struct {
		ID     string   `json:"_id"`
		Key    string   `json:"_key"`
		Rev    string   `json:"_rev"`
		SiteID string   `json:"siteId,omitempty"`
		Name   string   `json:"name,omitempty"`
		State  string   `json:"state,omitempty"`
		States []string `json:"states,omitempty"`
	}
)
