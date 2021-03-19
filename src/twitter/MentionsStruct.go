package twitter

// MentionsData Struct
type MentionsData struct {
	End      int16  `json:"end"`
	Start    int16  `json:"start"`
	Username string `json:"username"`
}

// MentionsEntityURL Struct
type MentionsEntityURL struct {
	URL         string `json:"url"`
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Start       int16  `json:"start"`
	End         int16  `json:"end"`
}

// MentionsHashtag Struct
type MentionsHashtag struct {
	Indicies []int32 `json:"indicies"`
	Text     string  `json:"text"`
}

// SizeType Struct
type SizeType struct {
	Height int32  `json:"h"`
	Width  int32  `json:"w"`
	Resize string `json:"resize"`
}

// MentionsMediaSize Struct
type MentionsMediaSize struct {
	Thumb  SizeType `json:"thumb"`
	Large  SizeType `json:"large"`
	Medium SizeType `json:"medium"`
	Small  SizeType `json:"small"`
}

// MentionsMedia Struct
type MentionsMedia struct {
	DisplayURL        string            `json:"display_url"`
	ExpandedURL       string            `json:"expanded_url"`
	ID                float64           `json:"id"`
	IDStr             string            `json:"id_str"`
	Indicies          []int32           `json:"indicies"`
	MediaURL          string            `json:"media_url"`
	MediaURLHttps     string            `json:"media_url_https"`
	Sizes             MentionsMediaSize `json:"sizes"`
	SourceStatusID    int64             `json:"source_status_id"`
	SourceStatusIDStr string            `json:"source_status_id_str"`
	Type              string            `json:"type"`
	URL               string            `json:"url"`
}

// Entities Object Structure
type Entities struct {
	Mentions []MentionsData      `json:"mentions"`
	Urls     []MentionsEntityURL `json:"urls"`
	Hashtags []MentionsHashtag   `json:"hashtags"`
	Media    []MentionsMedia     `json:"media"`
}
