package models

// PathResponse model
type PathResponse struct {
	Start string   `json:"start"`
	End   string   `json:"end"`
	Path  []string `json:"path"`
}
