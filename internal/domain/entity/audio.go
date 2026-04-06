package entity

// AudioItem represents a single audio guide entry.
// json tags are kept per Go convention; they can be removed in theory-only contexts.
type AudioItem struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Summary  *string `json:"summary"`
	URL      string  `json:"url"`
	FileExt  *string `json:"file_ext"`
	Modified string  `json:"modified"`
}

// AudioList is a paginated list of audio guide entries.
type AudioList struct {
	Total int         `json:"total"`
	Data  []AudioItem `json:"data"`
}
